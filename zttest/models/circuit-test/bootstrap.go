/*
	(c) Copyright NetFoundry Inc.

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
	"time"

	"github.com/hanzozt/fablab/kernel/lib/actions/semaphore"
	"github.com/hanzozt/zt/zttest/ztlab"

	"github.com/hanzozt/fablab/kernel/lib/actions"
	"github.com/hanzozt/fablab/kernel/lib/actions/component"
	"github.com/hanzozt/fablab/kernel/lib/actions/host"
	"github.com/hanzozt/fablab/kernel/model"
	ztlib_actions "github.com/hanzozt/zt/zttest/ztlab/actions"
	"github.com/hanzozt/zt/zttest/ztlab/actions/edge"
	"github.com/hanzozt/zt/zttest/ztlab/models"
)

type bootstrapAction struct{}

func NewBootstrapAction() model.ActionBinder {
	action := &bootstrapAction{}
	return action.bind
}

func (a *bootstrapAction) bind(m *model.Model) model.Action {
	workflow := actions.Workflow()

	workflow.AddAction(component.StopInParallel("*", 300))
	workflow.AddAction(host.GroupExec("*", 25, "rm -f logs/* ctrl.db"))
	workflow.AddAction(host.GroupExec("component.ctrl", 5, "rm -rf ./fablab/ctrldata"))

	isHA := len(m.SelectComponents(".ctrl")) > 1

	if !isHA {
		workflow.AddAction(component.Exec("#ctrl1", ztlab.ControllerActionInitStandalone))
	}

	workflow.AddAction(component.Start(".ctrl"))

	if isHA {
		workflow.AddAction(semaphore.Sleep(2 * time.Second))
		workflow.AddAction(edge.InitRaftController("#ctrl1"))
	}

	workflow.AddAction(edge.ControllerAvailable("#ctrl1", 30*time.Second))
	workflow.AddAction(edge.Login("#ctrl1"))

	workflow.AddAction(component.StopInParallel(models.EdgeRouterTag, 25))
	workflow.AddAction(edge.InitEdgeRouters(models.EdgeRouterTag, 2))
	workflow.AddAction(edge.InitIdentities(models.SdkAppTag, 2))

	// Loop Service
	workflow.AddAction(ztlib_actions.Edge("create", "config", "loop-host", "host.v1", `
		{
			"address" : "localhost",
			"port" : 3456,
			"protocol" : "tcp"
		}`))

	workflow.AddAction(ztlib_actions.Edge("create", "config", "throughput-intercept", "intercept.v1", `
		{
			"addresses": ["throughput.zt"],
			"portRanges" : [ { "low": 3456, "high": 3456 } ],
			"protocols": ["tcp"]
		}`))

	workflow.AddAction(ztlib_actions.Edge("create", "config", "latency-intercept", "intercept.v1", `
		{
			"addresses": ["latency.zt"],
			"portRanges" : [ { "low": 3456, "high": 3456 } ],
			"protocols": ["tcp"]
		}`))

	workflow.AddAction(ztlib_actions.Edge("create", "config", "throughput-xg-intercept", "intercept.v1", `
		{
			"addresses": ["throughput-xg.zt"],
			"portRanges" : [ { "low": 3456, "high": 3456 } ],
			"protocols": ["tcp"]
		}`))

	workflow.AddAction(ztlib_actions.Edge("create", "config", "latency-xg-intercept", "intercept.v1", `
		{
			"addresses": ["latency-xg.zt"],
			"portRanges" : [ { "low": 3456, "high": 3456 } ],
			"protocols": ["tcp"]
		}`))

	workflow.AddAction(ztlib_actions.Edge("create", "config", "slow-xg-intercept", "intercept.v1", `
		{
			"addresses": ["slow-xg.zt"],
			"portRanges" : [ { "low": 3456, "high": 3456 } ],
			"protocols": ["tcp"]
		}`))

	workflow.AddAction(ztlib_actions.Edge("create", "config", "throughput-ert-intercept", "intercept.v1", `
		{
			"addresses": ["throughput-ert.zt"],
			"portRanges" : [ { "low": 3456, "high": 3456 } ],
			"protocols": ["tcp"]
		}`))

	workflow.AddAction(ztlib_actions.Edge("create", "config", "latency-ert-intercept", "intercept.v1", `
		{
			"addresses": ["latency-ert.zt"],
			"portRanges" : [ { "low": 3456, "high": 3456 } ],
			"protocols": ["tcp"]
		}`))

	workflow.AddAction(ztlib_actions.Edge("create", "config", "slow-ert-intercept", "intercept.v1", `
		{
			"addresses": ["slow-ert.zt"],
			"portRanges" : [ { "low": 3456, "high": 3456 } ],
			"protocols": ["tcp"]
		}`))

	workflow.AddAction(ztlib_actions.Edge("create", "service", "throughput", "-a", "loop,loop-host", "-c", "throughput-intercept"))
	workflow.AddAction(ztlib_actions.Edge("create", "service", "latency", "-a", "loop,loop-host", "-c", "latency-intercept"))

	workflow.AddAction(ztlib_actions.Edge("create", "service", "throughput-xg", "-a", "loop,loop-host-xg", "-c", "throughput-xg-intercept"))
	workflow.AddAction(ztlib_actions.Edge("create", "service", "latency-xg", "-a", "loop,loop-host-xg", "-c", "latency-xg-intercept"))
	workflow.AddAction(ztlib_actions.Edge("create", "service", "slow-xg", "-a", "loop,loop-host-xg", "-c", "slow-xg-intercept"))

	workflow.AddAction(ztlib_actions.Edge("create", "service", "throughput-ert", "-a", "loop,loop-host-ert", "-c", "loop-host,throughput-ert-intercept"))
	workflow.AddAction(ztlib_actions.Edge("create", "service", "latency-ert", "-a", "loop,loop-host-ert", "-c", "loop-host,latency-ert-intercept"))
	workflow.AddAction(ztlib_actions.Edge("create", "service", "slow-ert", "-a", "loop,loop-host-ert", "-c", "loop-host,slow-ert-intercept"))

	workflow.AddAction(ztlib_actions.Edge("create", "service-policy", "loop-hosts", "Bind", "--service-roles", "#loop-host", "--identity-roles", "#loop-host"))
	workflow.AddAction(ztlib_actions.Edge("create", "service-policy", "loop-hosts-xg", "Bind", "--service-roles", "#loop-host-xg", "--identity-roles", "#loop-host-xg"))
	workflow.AddAction(ztlib_actions.Edge("create", "service-policy", "loop-hosts-ert", "Bind", "--service-roles", "#loop-host-ert", "--identity-roles", "#loop-host-ert"))
	workflow.AddAction(ztlib_actions.Edge("create", "service-policy", "loop-clients", "Dial", "--service-roles", "#loop", "--identity-roles", "#loop-client"))
	workflow.AddAction(ztlib_actions.Edge("create", "service-edge-router-policy", "loop", "--service-roles", "#loop", "--edge-router-roles", "#test"))

	// Sim Services
	workflow.AddAction(ztlib_actions.Edge("create", "service", "metrics", "-a", "sim-services"))
	workflow.AddAction(ztlib_actions.Edge("create", "service", "sim-control", "-a", "sim-services"))

	workflow.AddAction(ztlib_actions.Edge("create", "service-policy", "sim-service-hosts", "Bind", "--service-roles", "#sim-services", "--identity-roles", "#sim-services-host"))
	workflow.AddAction(ztlib_actions.Edge("create", "service-policy", "sim-service-clients", "Dial", "--service-roles", "#sim-services", "--identity-roles", "#sim-services-client"))

	workflow.AddAction(ztlib_actions.Edge("create", "edge-router-policy", "sim-services-hosts", "--edge-router-roles", "#sim-services", "--identity-roles", "#all"))
	workflow.AddAction(ztlib_actions.Edge("create", "service-edge-router-policy", "sim-services", "--service-roles", "#sim-services", "--edge-router-roles", "#sim-services"))

	// Shared policies
	workflow.AddAction(ztlib_actions.Edge("create", "edge-router-policy", "hosts", "--edge-router-roles", "#host", "--identity-roles", "#host"))
	workflow.AddAction(ztlib_actions.Edge("create", "edge-router-policy", "clients", "--edge-router-roles", "#client", "--identity-roles", "#client"))

	workflow.AddAction(semaphore.Sleep(2 * time.Second))
	workflow.AddAction(edge.RaftJoin("ctrl1", ".ctrl"))
	workflow.AddAction(semaphore.Sleep(5 * time.Second))

	workflow.AddAction(component.StartInParallel(".tcpdump", 10))
	workflow.AddAction(component.StartInParallel(models.EdgeRouterTag, 10))
	workflow.AddAction(semaphore.Sleep(2 * time.Second))
	workflow.AddAction(component.StartInParallel(".sim-services-host", 50))
	workflow.AddAction(component.StartInParallel(".sim-services-client", 50))

	workflow.AddAction(model.ActionFunc(func(run model.Run) error {
		if run.GetModel().BoolVariable("tcpdump") {
			time.Sleep(30 * time.Second)
			return run.GetModel().ForEachComponent(".tcpdump", 10, func(c *model.Component) error {
				return c.Host.ExecLogOnlyOnError(fmt.Sprintf("cp /home/ubuntu/logs/%s.pcap0 /home/ubuntu/logs/%s.pcap.first", c.Id, c.Id))
			})
		}
		return nil
	}))

	return workflow
}
