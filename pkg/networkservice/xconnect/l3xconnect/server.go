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
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/vppagent"
	"github.com/networkservicemesh/sdk/pkg/networkservice/core/next"
)

// Package l3xconnect provides a NetworkServiceServer chain element for an l3 cross connect
type l3XconnectServer struct{}

// NewServer - creates a NetworkServiceServer chain element for an l3 cross connect
func NewServer() networkservice.NetworkServiceServer {
	return &l3XconnectServer{}
}

func (l *l3XconnectServer) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*networkservice.Connection, error) {
	conf := vppagent.Config(ctx)
	appendL3XConnect(conf, request.Connection)
	return next.Server(ctx).Request(ctx, request)
}

func (l *l3XconnectServer) Close(ctx context.Context, conn *networkservice.Connection) (*empty.Empty, error) {
	conf := vppagent.Config(ctx)
	appendL3XConnect(conf, conn)
	return next.Server(ctx).Close(ctx, conn)
}