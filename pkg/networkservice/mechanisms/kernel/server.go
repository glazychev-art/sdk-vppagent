// Copyright (c) 2020 Cisco Systems, Inc.
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

// +build !windows

// Package kernel provides a networkservice chain element that properly handles
// the kernel Mechanism
package kernel

import (
	"os"

	"github.com/networkservicemesh/api/pkg/api/networkservice"

	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/mechanisms/kernel/kerneltap"
	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/mechanisms/kernel/kernelvethpair"
)

// NewServer return a NetworkServiceServer chain element that correctly handles the kernel Mechanism
func NewServer() networkservice.NetworkServiceServer {
	if _, err := os.Stat(vnetFilename); err == nil {
		return kerneltap.NewServer()
	}
	return kernelvethpair.NewServer()
}
