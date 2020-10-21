// Copyright (c) 2020 Doc.ai and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package l3xconnect

import (
	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/wireguard"
	"go.ligato.io/vpp-agent/v3/proto/ligato/configurator"
	l3 "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l3"
	"strings"
)

func appendL3XConnect(conf *configurator.Config, conn *networkservice.Connection ) {
	if len(conf.GetVppConfig().GetInterfaces()) >= 2 {
		ifaces := conf.GetVppConfig().GetInterfaces()[len(conf.GetVppConfig().Interfaces)-2:]
		nextHops := []string{ "", "" }
		if conn.GetMechanism().GetType() == wireguard.MECHANISM {
			mechanism := wireguard.ToMechanism(conn.GetMechanism())
			nextHops[0] = mechanism.SrcIP().String()
			nextHops[1] = mechanism.DstIP().String()
		}

		conf.VppConfig.L3Xconnects = append(conf.VppConfig.L3Xconnects,
			&l3.L3XConnect{
				Interface: ifaces[0].Name,
				Protocol:  getL3XProtocol(nextHops[1]),
				Paths:     []*l3.L3XConnect_Path{
					{
						OutgoingInterface: ifaces[1].Name,
						NextHopAddr:       nextHops[1],
					},
				},
			},
			&l3.L3XConnect{
				Interface: ifaces[1].Name,
				Protocol:  getL3XProtocol(nextHops[0]),
				Paths:     []*l3.L3XConnect_Path{
					{
						OutgoingInterface: ifaces[0].Name,
						NextHopAddr:       nextHops[0],
					},
				},
		})
	}
}

func getL3XProtocol(address string) l3.L3XConnect_Protocol {
	if strings.Count(address, ":") >= 2 {
		return l3.L3XConnect_IPV6
	}
	return l3.L3XConnect_IPV4
}