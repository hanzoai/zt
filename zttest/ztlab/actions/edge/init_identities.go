package edge

import (
	"path/filepath"
	"strings"

	"github.com/hanzozt/fablab/kernel/libssh"
	"github.com/hanzozt/fablab/kernel/model"
	ztlib_actions "github.com/hanzozt/zt/zttest/ztlab/actions"
	"github.com/hanzozt/zt/zttest/ztlab/cli"
)

func InitIdentities(componentSpec string, concurrency int) model.Action {
	return &initIdentitiesAction{
		componentSpec: componentSpec,
		concurrency:   concurrency,
	}
}

func (action *initIdentitiesAction) Execute(run model.Run) error {
	return run.GetModel().ForEachComponent(action.componentSpec, action.concurrency, func(c *model.Component) error {
		if err := ztlib_actions.EdgeExec(run.GetModel(), "delete", "identity", c.Id); err != nil {
			return err
		}

		return action.createAndEnrollIdentity(run, c)
	})
}

func (action *initIdentitiesAction) createAndEnrollIdentity(run model.Run, c *model.Component) error {
	ssh := c.GetHost().NewSshConfigFactory()

	jwtFileName := filepath.Join(run.GetTmpDir(), c.Id+".jwt")

	err := ztlib_actions.EdgeExec(c.GetModel(), "create", "identity", c.Id,
		"--jwt-output-file", jwtFileName,
		"-a", strings.Join(c.Tags, ","))

	if err != nil {
		return err
	}

	configFileName := filepath.Join(run.GetTmpDir(), c.Id+".json")

	_, err = cli.Exec(c.GetModel(), "edge", "enroll", "--jwt", jwtFileName, "--out", configFileName)

	if err != nil {
		return err
	}

	remoteConfigFile := "/home/ubuntu/fablab/cfg/" + c.Id + ".json"
	return libssh.SendFile(ssh, configFileName, remoteConfigFile)
}

type initIdentitiesAction struct {
	componentSpec string
	concurrency   int
}
