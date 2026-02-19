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
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hanzozt/sdk-golang/zt"
)

func NewZitiEchoClient(identityJson string) (*ztEchoClient, error) {
	config, err := zt.NewConfigFromFile(identityJson)
	if err != nil {
		return nil, err
	}

	ztContext, err := zt.NewContext(config)

	if err != nil {
		return nil, err
	}

	dial := func(_ context.Context, _ string, addr string) (net.Conn, error) {
		service := strings.Split(addr, ":")[0] // assume host is service
		return ztContext.Dial(service)
	}

	ztTransport := http.DefaultTransport.(*http.Transport).Clone()
	ztTransport.DialContext = dial

	return &ztEchoClient{
		httpClient: &http.Client{Transport: ztTransport},
	}, nil
}

type ztEchoClient struct {
	httpClient *http.Client
}

func (self *ztEchoClient) echo(input string) error {
	u := fmt.Sprintf("http://echo?input=%v", url.QueryEscape(input))
	resp, err := self.httpClient.Get(u)
	if err == nil {
		c := color.New(color.FgGreen, color.Bold)
		_, _ = c.Print("\nzt-http-echo-client: ")
		_, err = io.Copy(os.Stdout, resp.Body)
	}
	return err
}
