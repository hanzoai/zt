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
	"embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/michaelquigley/pfxlog"
	"github.com/hanzozt/fablab"
	"github.com/hanzozt/fablab/kernel/lib/actions"
	"github.com/hanzozt/fablab/kernel/lib/actions/component"
	"github.com/hanzozt/fablab/kernel/lib/actions/semaphore"
	"github.com/hanzozt/fablab/kernel/lib/binding"
	"github.com/hanzozt/fablab/kernel/lib/runlevel/0_infrastructure/aws_ssh_key"
	semaphore0 "github.com/hanzozt/fablab/kernel/lib/runlevel/0_infrastructure/semaphore"
	terraformInit "github.com/hanzozt/fablab/kernel/lib/runlevel/0_infrastructure/terraform"
	distribution "github.com/hanzozt/fablab/kernel/lib/runlevel/3_distribution"
	"github.com/hanzozt/fablab/kernel/lib/runlevel/3_distribution/rsync"
	fablibOps "github.com/hanzozt/fablab/kernel/lib/runlevel/5_operation"
	awsSshKeyDispose "github.com/hanzozt/fablab/kernel/lib/runlevel/6_disposal/aws_ssh_key"
	"github.com/hanzozt/fablab/kernel/lib/runlevel/6_disposal/terraform"
	"github.com/hanzozt/fablab/kernel/model"
	"github.com/hanzozt/fablab/kernel/model/aws"
	"github.com/hanzozt/fablab/resources"
	"github.com/hanzozt/zt/zttest/models/test_resources"
	"github.com/hanzozt/zt/zttest/ztlab"
	ztlibActions "github.com/hanzozt/zt/zttest/ztlab/actions"
	"github.com/hanzozt/zt/zttest/ztlab/actions/edge"
	"github.com/hanzozt/zt/zttest/ztlab/models"
	ztlibOps "github.com/hanzozt/zt/zttest/ztlab/runlevel/5_operation"
)

//go:embed configs
var configResource embed.FS

func getUniqueId() string {
	if runId := os.Getenv("GITHUB_RUN_ID"); runId != "" {
		return "-" + runId + "." + os.Getenv("GITHUB_RUN_ATTEMPT")
	}
	return "-" + os.Getenv("USER")
}

var Model = &model.Model{
	Id: "circuit-perf",
	Scope: model.Scope{
		Defaults: model.Variables{
			"environment": "circuit-perf-" + getUniqueId(),
			"credentials": model.Variables{
				"aws": model.Variables{
					"managed_key": true,
				},
				"ssh": model.Variables{
					"username": "ubuntu",
				},
				"edge": model.Variables{
					"username": "admin",
					"password": "admin",
				},
			},
			"metrics": model.Variables{
				"influxdb": model.Variables{
					"url": "http://localhost:8086",
					"db":  "zt",
				},
			},
			"zt_version": "v1.6.6",
		},
	},

	Factories: []model.Factory{
		model.FactoryFunc(func(m *model.Model) error {
			pfxlog.Logger().Infof("environment [%s]", m.MustStringVariable("environment"))
			m.AddActivationActions("stop", "bootstrap", "start")
			return nil
		}),
		model.FactoryFunc(func(m *model.Model) error {
			return m.ForEachHost("*", 1, func(host *model.Host) error {
				if host.InstanceType == "" {
					host.InstanceType = "c5.xlarge"
				}
				return nil
			})
		}),
		model.FactoryFunc(func(m *model.Model) error {
			return m.ForEachComponent(".underTest", 1, func(component *model.Component) error {
				if vt, ok := component.Type.(model.VersionableComponent); ok {
					vt.SetVersion(m.GetStringVariableOr("zt_version", ""))
					return nil
				}
				return fmt.Errorf("component %s of type %T doesn't support setting version", component.Id, component.Type)
			})
		}),
		model.FactoryFunc(func(m *model.Model) error {
			simServices := ztlibOps.NewSimServices(func(s string) string {
				return "component#" + s
			})

			m.AddActivationStageF(simServices.SetupSimControllerIdentity)
			m.AddOperatingStage(simServices.CollectSimMetricStage("metrics"))

			m.AddOperatingStageF(func(run model.Run) error {
				waitC := make(chan struct{})
				<-waitC
				return nil
			})

			return nil
		}),
	},

	Resources: model.Resources{
		resources.Configs:   resources.SubFolder(configResource, "configs"),
		resources.Terraform: test_resources.TerraformResources(),
	},

	Regions: model.Regions{
		"us-east-1": {
			Region: "us-east-1",
			Site:   "us-east-1a",
			Hosts: model.Hosts{
				"ctrl": {
					Components: model.Components{
						"ctrl": {
							Scope: model.Scope{Tags: model.Tags{"ctrl", "underTest"}},
							Type:  &ztlab.ControllerType{},
						},
					},
				},
				"router-client": {
					Components: model.Components{
						"router-client": {
							Scope: model.Scope{Tags: model.Tags{"edge-router", "tunneler", "client", "test", "underTest"}},
							Type: &ztlab.RouterType{
								Debug: false,
							},
						},
					},
				},
				"router-metrics": {
					Components: model.Components{
						"router-metrics": {
							Scope: model.Scope{Tags: model.Tags{"edge-router", "no-traversal", "metrics"}},
							Type: &ztlab.RouterType{
								Debug: false,
							},
						},
					},
				},
				"zt-tunnel-client": {
					InstanceResourceType: "ondemand_iops",
					AWS: aws.EC2Host{
						Volume: aws.EC2Volume{
							Type:   "gp3",
							SizeGB: 20,
							IOPS:   1000,
						},
					},
					Components: model.Components{
						"zt-tunnel-client": {
							Scope: model.Scope{Tags: model.Tags{"zt-tunnel", "sdk-app", "client", "ssh-client", "underTest"}},
							Type: &ztlab.ZitiTunnelType{
								ControlRouterConnections: 1,
								DefaultRouterConnections: 2,
							},
						},
					},
				},
				"loop-client": {
					Scope: model.Scope{Tags: model.Tags{"loop-client"}},
					Components: model.Components{
						"loop-client": {
							Scope: model.Scope{Tags: model.Tags{"loop-client", "sdk-app", "client", "metrics-client"}},
							Type: &ztlab.Loop4SimType{
								Mode: ztlab.Loop4Dialer,
							},
						},
					},
				},
			},
		},
		"us-west-2": {
			Region: "us-west-2",
			Site:   "us-west-2b",
			Hosts: model.Hosts{
				"router-host": {
					Components: model.Components{
						"router-host": {
							Scope: model.Scope{Tags: model.Tags{"edge-router", "tunneler", "host", "ert-host", "test", "loop-host", "underTest"}},
							Type: &ztlab.RouterType{
								Debug: false,
							},
						},

						"loop-host": {
							Scope: model.Scope{Tags: model.Tags{"loop-host", "sdk-app", "host"}},
							Type: &ztlab.Loop4SimType{
								Mode: ztlab.Loop4Listener,
							},
						},
					},
				},
				"zt-tunnel-host": {
					Components: model.Components{
						"zt-tunnel-host": {
							Scope: model.Scope{Tags: model.Tags{"zt-tunnel", "sdk-app", "host", "ssh-host", "underTest"}},
							Type: &ztlab.ZitiTunnelType{
								Mode:                     ztlab.ZitiTunnelModeHost,
								ControlRouterConnections: 1,
								DefaultRouterConnections: 2,
							},
						},
					},
				},
			},
		},
	},

	Actions: model.ActionBinders{
		"bootstrap": NewBootstrapAction(),
		"start": model.BindF(func(run model.Run) error {
			workflow := actions.Workflow()
			workflow.AddAction(component.Start(".ctrl"))
			workflow.AddAction(edge.ControllerAvailable("#ctrl", 30*time.Second))
			workflow.AddAction(component.StartInParallel(models.EdgeRouterTag, 25))

			workflow.AddAction(semaphore.Sleep(5 * time.Second))
			workflow.AddAction(component.StartInParallel("loop-host", 5))
			workflow.AddAction(component.Start(".zt-tunnel"))

			workflow.AddAction(edge.Login("#ctrl"))
			workflow.AddAction(ztlibActions.Edge("list", "edge-routers", "limit none"))
			workflow.AddAction(ztlibActions.Edge("list", "terminators", "limit none"))

			return workflow.Execute(run)
		}),
		"stop":  model.Bind(component.StopInParallel("*", 15)),
		"login": model.Bind(edge.Login("#ctrl")),
		"testXgress": model.BindF(func(run model.Run) error {
			run.GetModel().Scope.PutVariable("zt_version", "")
			return run.GetModel().Operate(run)
		}),
		"testNoXgress": model.BindF(func(run model.Run) error {
			run.GetModel().Scope.PutVariable("zt_version", "v1.5.4")
			return run.GetModel().Operate(run)
		}),
	},

	Infrastructure: model.Stages{
		aws_ssh_key.Express(),
		&terraformInit.Terraform{
			Retries: 3,
			ReadyCheck: &semaphore0.ReadyStage{
				MaxWait: 90 * time.Second,
			},
		},
	},

	Distribution: model.Stages{
		distribution.DistributeSshKey("*"),
		rsync.RsyncStaged(),
	},

	Disposal: model.Stages{
		terraform.Dispose(),
		awsSshKeyDispose.Dispose(),
	},

	Operation: model.Stages{
		model.RunAction("login"),
		edge.SyncModelEdgeState(models.EdgeRouterTag),

		fablibOps.StreamSarMetrics("*", 5, 1, nil),

		fablibOps.InfluxMetricsReporter(),

		ztlibOps.ModelMetricsWithIdMapper(nil, func(id string) string {
			if id == "ctrl" {
				return "#ctrl"
			}
			id = strings.ReplaceAll(id, ".", ":")
			return "component.edgeId:" + id
		}),

		component.Stop("loop-client"),
		component.Start("loop-client"),
	},
}

func InitBootstrapExtensions() {
	model.AddBootstrapExtension(binding.AwsCredentialsLoader)
	model.AddBootstrapExtension(aws_ssh_key.KeyManager)
}

func main() {
	InitBootstrapExtensions()
	fablab.InitModel(Model)
	fablab.Run()
}
