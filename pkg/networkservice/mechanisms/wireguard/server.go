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
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/wireguard"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/sdk/pkg/networkservice/core/next"
)

type wireguardServer struct {
	privateKey string
	publicKey  string
}

func NewServer() networkservice.NetworkServiceServer {
	key, _ := wgtypes.GeneratePrivateKey()
	return &wireguardServer{
		privateKey: key.String(),
		publicKey:  key.PublicKey().String(),
	}
}

func (v *wireguardServer) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*networkservice.Connection, error) {
	mechanism := wireguard.ToMechanism(request.GetConnection().GetMechanism())
	if mechanism != nil && request.GetConnection().GetMechanism().Parameters != nil {
		if err := appendInterfaceConfig(ctx, request.GetConnection(), v.privateKey, v.publicKey, true); err != nil {
			return nil, err
		}
	}
	return next.Server(ctx).Request(ctx, request)
}

func (v *wireguardServer) Close(ctx context.Context, conn *networkservice.Connection) (*empty.Empty, error) {
	if err := appendInterfaceConfig(ctx, conn, v.privateKey, v.publicKey, true); err != nil {
		return nil, err
	}
	return next.Server(ctx).Close(ctx, conn)
}
