/*
	Copyright 2019 NetFoundry Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hanzozt/fablab/kernel/lib/parallel"
	"github.com/hanzozt/fablab/kernel/model"
	"github.com/hanzozt/zt/zttest/ztlab"
	ztlib_actions "github.com/hanzozt/zt/zttest/ztlab/actions"
	"github.com/hanzozt/zt/zttest/ztlab/cli"
)

func provisionIdentities(identities []string, run model.Run) error {
	var tasks []parallel.Task

	identitiesDir := model.MakeBuildPath("identities")
	_ = os.MkdirAll(identitiesDir, 0770)
	for _, v := range identities {
		id := v
		task := func() error {
			jwtFileName := filepath.Join(identitiesDir, id+".jwt")
			args := []string{"create", "enrollment", "ott", "--jwt-output-file", jwtFileName, "--", id}

			if err := ztlib_actions.EdgeExec(run.GetModel(), args...); err != nil {
				return err
			}

			if _, err := cli.Exec(m, "edge", "enroll", "--", jwtFileName); err != nil {
				return fmt.Errorf("failed to enroll %s (%w)", id, err)
			}
			return nil
		}
		tasks = append(tasks, task)
	}

	if err := parallel.Execute(tasks, 10); err != nil {
		return err
	}

	return nil
}

func provisionRouters(run model.Run) error {
	routerPkiDir := model.MakeBuildPath("router-jwts")
	_ = os.MkdirAll(routerPkiDir, 0770)

	return run.GetModel().ForEachComponent(".router.pre-created", 15, func(c *model.Component) error {
		jwtFileName := filepath.Join(routerPkiDir, c.Id+".jwt")
		args := []string{"re-enroll", "edge-router", "-j", "--jwt-output-file", jwtFileName, "--", c.Id}
		return ztlib_actions.EdgeExec(c.GetModel(), args...)
	})
}

func enrollRouters(run model.Run) error {
	ztVersion := ""
	ctrls := run.GetModel().SelectComponents(".ctrl")
	for _, ctrl := range ctrls {
		if ctrl.Type != nil {
			ztVersion = ctrl.Type.GetVersion()
			break
		}
	}

	return run.GetModel().ForEachHost("component.router", 15, func(h *model.Host) error {
		var cmds []string
		cmds = append(cmds, "mkdir -p /home/ubuntu/router-pki/")
		hostUser := h.GetSshUser()
		for _, c := range h.Components {
			var routerType *ztlab.RouterType
			var ok bool
			if routerType, ok = c.Type.(*ztlab.RouterType); !ok || !c.HasTag("pre-created") {
				continue
			}

			remoteJwt := fmt.Sprintf("/home/%s/router-jwts/%s.jwt", hostUser, c.Id)
			tmpl := "set -o pipefail; %s router enroll /home/ubuntu/fablab/cfg/%s -j %s 2>&1 | tee /home/ubuntu/logs/%s.router.enroll.log "
			cmd := fmt.Sprintf(tmpl, ztlab.GetZitiBinaryPath(c, ztVersion), routerType.GetConfigName(c), remoteJwt, c.Id)
			cmds = append(cmds, cmd)
		}
		return h.ExecLogOnlyOnError(cmds...)
	})
}
