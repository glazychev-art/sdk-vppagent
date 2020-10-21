package wireguard_test

import (
	"context"
	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/cls"
	"github.com/networkservicemesh/api/pkg/api/registry"
	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/connectioncontextkernel"
	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/mechanisms/wireguard"
	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/vppagent"
	"github.com/networkservicemesh/sdk/pkg/networkservice/chains/client"
	"github.com/networkservicemesh/sdk/pkg/networkservice/chains/endpoint"
	"github.com/networkservicemesh/sdk/pkg/networkservice/common/authorize"
	"github.com/networkservicemesh/sdk/pkg/networkservice/common/clienturl"
	"github.com/networkservicemesh/sdk/pkg/networkservice/common/connect"
	"github.com/networkservicemesh/sdk/pkg/networkservice/common/mechanisms"
	"github.com/networkservicemesh/sdk/pkg/networkservice/common/null"
	"github.com/networkservicemesh/sdk/pkg/networkservice/core/adapters"
	"github.com/networkservicemesh/sdk/pkg/tools/addressof"
	"github.com/networkservicemesh/sdk/pkg/tools/sandbox"
	"github.com/networkservicemesh/sdk/pkg/tools/token"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/url"
	"testing"
	"time"
)

func wireguardSupplier(ctx context.Context, name string, generateToken token.GeneratorFunc, connectTo *url.URL, dialOptions ...grpc.DialOption) endpoint.Endpoint {
	var result endpoint.Endpoint
	result = endpoint.NewServer(ctx,
		name,
		authorize.NewServer(),
		generateToken,
		vppagent.NewServer(),
		mechanisms.NewServer(map[string]networkservice.NetworkServiceServer{
			//kernel.MECHANISM:    kernel.NewServer(),
			wireguard.MECHANISM: wireguard.NewServer(),
		}),
		clienturl.NewServer(connectTo),
		connect.NewServer(
			ctx,
			client.NewClientFactory(
				name,
				// What to call onHeal
				addressof.NetworkServiceClient(adapters.NewServerToClient(result)),
				generateToken,
				connectioncontextkernel.NewClient(),
				// Preference ordered list of mechanisms we support for outgoing connections
				//kernel.NewClient(),
				wireguard.NewClient(),
			),
			grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		),
		//directmemif.NewServer(),
		connectioncontextkernel.NewServer(),
		)

	return result
}

func TestNSMGR_RemoteUsecase2(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())
	logrus.SetOutput(ioutil.Discard)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5000)
	defer cancel()
	domain := sandbox.NewBuilder(t).
		SetNodesCount(2).
		SetForwarderSupplier(wireguardSupplier).
		SetRegistryProxySupplier(nil).
		SetContext(ctx).
		Build()
	defer domain.Cleanup()

	nseReg := &registry.NetworkServiceEndpoint{
		Name:                "final-endpoint",
		NetworkServiceNames: []string{"my-service-remote"},
	}

	entryy, err := sandbox.NewEndpoint(ctx, nseReg, sandbox.GenerateTestToken, domain.Nodes[0].NSMgr, null.NewServer())
	require.NoError(t, err)
	if entryy == nil {

	}

	request := &networkservice.NetworkServiceRequest{
		MechanismPreferences: []*networkservice.Mechanism{
			{Cls: cls.REMOTE, Type: wireguard.MECHANISM},
		},
		Connection: &networkservice.Connection{
			Id:             "1",
			NetworkService: "my-service-remote",
			Context:        &networkservice.ConnectionContext{},
		},
	}

	nsc, err := sandbox.NewClient(ctx, sandbox.GenerateTestToken, domain.Nodes[1].NSMgr.URL)
	require.NoError(t, err)

	conn, err := nsc.Request(ctx, request)
	require.NoError(t, err)
	require.NotNil(t, conn)

	require.Equal(t, 8, len(conn.Path.PathSegments))
}