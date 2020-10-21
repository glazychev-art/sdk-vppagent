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

package wireguard

import (
	"context"
	"fmt"
	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/wireguard"
	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/vppagent"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp"
	vpp_interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
	"strings"
)

const (
	// Persistent keepalive interval (sec)
	defaultPersistentKeepalive = 20
)

func appendInterfaceConfig(ctx context.Context, conn *networkservice.Connection, privateKey, publicKey string, isIncoming bool) error {
	vppWgInterfaceName := fmt.Sprintf("wg-%v", conn.GetId())

	mechanism := wireguard.ToMechanism(conn.Mechanism)
	conf := vppagent.Config(ctx)
	vppConfig := conf.GetVppConfig()

	/* Create interface */
	var (
		localPrivateKey   string
		remotePublicKey   string
		localPort         uint32
		remotePort        uint32
		srcIP             string
		remoteIP          string
		wireguardSrcIP    string
		wireguardRemoteIP string
	)

	localPrivateKey = privateKey

	if isIncoming {

		mechanism.SetDstPublicKey(publicKey)
		conn.Mechanism.Parameters[wireguard.DstPort] = wireguard.GetPort(conn.GetId())

		remotePublicKey = mechanism.SrcPublicKey()
		localPort = mechanism.DstPort()
		remotePort = mechanism.SrcPort()
		srcIP = mechanism.DstIP().String()
		remoteIP = mechanism.SrcIP().String()

		wireguardSrcIP = conn.GetContext().GetIpContext().GetDstIpAddr()
		wireguardRemoteIP = conn.GetContext().GetIpContext().GetSrcIpAddr()
	} else {
		mechanism.SetSrcPublicKey(publicKey)

		remotePublicKey = mechanism.DstPublicKey()
		localPort = mechanism.SrcPort()
		remotePort = mechanism.DstPort()
		srcIP = mechanism.SrcIP().String()
		remoteIP = mechanism.DstIP().String()

		wireguardSrcIP = conn.GetContext().GetIpContext().GetSrcIpAddr()
		wireguardRemoteIP = conn.GetContext().GetIpContext().GetDstIpAddr()
	}

	vppConfig.Interfaces = append(vppConfig.Interfaces, &vpp.Interface{
		Name:        vppWgInterfaceName,
		IpAddresses: []string{wireguardSrcIP},
		Enabled:     true,
		Link: &vpp_interfaces.Interface_Wireguard{Wireguard: &vpp_interfaces.WireguardLink{
			PrivateKey: localPrivateKey,
			Port:       localPort,
			SrcAddr:    strings.Split(srcIP, "/")[0],
		}},
		Type: vpp_interfaces.Interface_WIREGUARD_TUNNEL,
	})

	vppConfig.WgPeers = append(vppConfig.WgPeers, &vpp.WgPeer{
		PublicKey:           remotePublicKey,
		Endpoint:            strings.Split(remoteIP, "/")[0],
		WgIfName:            vppWgInterfaceName,
		Port:                remotePort,
		PersistentKeepalive: defaultPersistentKeepalive,
		AllowedIps:          []string{wireguardRemoteIP},
	})

	return nil
}