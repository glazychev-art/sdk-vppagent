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

package wireguard_test

import (
	"github.com/networkservicemesh/api/pkg/api/networkservice"
	wireguard_mechanism "github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/wireguard"
	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/mechanisms/checkvppagentmechanism"
	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/mechanisms/wireguard"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.ligato.io/vpp-agent/v3/proto/ligato/configurator"
	"io/ioutil"
	"testing"
)

func TestWireguardClient(t *testing.T) {
	// Turn off log output
	logrus.SetOutput(ioutil.Discard)

	tp := configureTestWireguardParameters()
	testMech := configureTestWireguardMechanism(tp.ToMap())

	testRequest := &networkservice.NetworkServiceRequest{
		Connection: &networkservice.Connection{
			Id: "1",
			Mechanism: testMech,
		},
	}

	suite.Run(t, checkvppagentmechanism.NewClientSuite(
		wireguard.NewClient(),
		wireguard_mechanism.MECHANISM,
		func(t *testing.T, mechanism *networkservice.Mechanism) {
			m := wireguard_mechanism.ToMechanism(mechanism)
			require.NotNil(t, m)
		},
		func(t *testing.T, conf *configurator.Config) {
			// Basic interface check
			vppConfig := conf.GetVppConfig()
			vppInterfaces := vppConfig.GetInterfaces()
			require.Greater(t, len(vppInterfaces), 0)
			vppInterface := vppInterfaces[len(vppInterfaces)-1]
			assert.NotNil(t, vppInterface)

			// Wireguard interface parameters check
			wireguardInterface := vppInterface.GetWireguard()
			assert.NotNil(t, wireguardInterface)
			assert.NotEqual(t, "", wireguardInterface.GetPrivateKey())
			assert.Equal(t, tp.srcIP, wireguardInterface.GetSrcAddr())
			assert.Equal(t, uint32(tp.srcPort), wireguardInterface.GetPort())

			// Basic peer check
			wireguardPeers := vppConfig.GetWgPeers()
			require.Greater(t, len(wireguardPeers), 0)
			wireguardPeer := wireguardPeers[len(wireguardPeers)-1]
			assert.NotNil(t, wireguardPeer)

			// Wireguard peer parameters check
			//assert.NotEqual(t, "", wireguardPeer.GetPublicKey())
			assert.Equal(t, uint32(tp.dstPort), wireguardPeer.GetPort())
			assert.Equal(t, tp.dstIP, wireguardPeer.GetEndpoint())
			assert.Equal(t, vppInterface.GetName(), wireguardPeer.GetWgIfName())
		},
		testRequest,
		testRequest.GetConnection(),
	))
}
