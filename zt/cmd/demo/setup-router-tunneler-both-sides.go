/*
	Copyright NetFoundry Inc.

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

package demo

import (
	_ "embed"
	"time"

	"github.com/hanzozt/runzmd"
	"github.com/hanzozt/runzmd/actionz"
	"github.com/hanzozt/zt/v2/zt/cmd/api"
	"github.com/hanzozt/zt/v2/zt/cmd/common"
	"github.com/spf13/cobra"
)

//go:embed setup-scripts/router-tunneler-both-sides.md
var routerTunnelerBothSidesScriptSource []byte

type routerTunnelerBothSides struct {
	api.Options
	TutorialOptions
	interactive bool
}

func newRouterTunnelerBothSidesCmd(p common.OptionsProvider) *cobra.Command {
	options := &routerTunnelerBothSides{
		Options: api.Options{
			CommonOptions: p(),
		},
	}

	cmd := &cobra.Command{
		Use:   "router-tunneler-both-sides",
		Short: "Walks you through configuration for an echo service with intercept and hosting both on router embedded tunnelers",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.Cmd = cmd
			options.Args = args
			return options.run()
		},
		SuggestFor: []string{},
	}

	// allow interspersing positional args and flags
	cmd.Flags().SetInterspersed(true)
	cmd.Flags().StringVar(&options.ControllerUrl, "controller-url", "", "The Hanzo ZT controller URL to use")
	cmd.Flags().StringVarP(&options.Username, "username", "u", "", "The Hanzo ZT controller username to use")
	cmd.Flags().StringVarP(&options.Password, "password", "p", "", "The Hanzo ZT controller password to use")
	cmd.Flags().DurationVar(&options.NewlinePause, "newline-pause", time.Millisecond*10, "How long to pause between lines when scrolling")
	cmd.Flags().BoolVar(&options.interactive, "interactive", false, "Interactive mode, waiting for user input")
	options.AddCommonFlags(cmd)

	return cmd
}

func (self *routerTunnelerBothSides) run() error {
	t := runzmd.NewRunner()
	t.NewLinePause = self.NewlinePause
	t.AssumeDefault = !self.interactive

	t.RegisterActionHandler("zt", &actionz.ZitiRunnerAction{})
	t.RegisterActionHandler("zt-login", &actionz.ZitiEnsureLoggedIn{
		LoginParams: &self.TutorialOptions,
	})
	t.RegisterActionHandler("keep-session-alive", &actionz.KeepSessionAliveAction{})
	t.RegisterActionHandler("zt-create-config", &actionz.ZitiCreateConfigAction{})
	t.RegisterActionHandler("zt-for-each", &actionz.ZitiForEach{})

	return t.Run(routerTunnelerBothSidesScriptSource)
}
