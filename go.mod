module github.com/networkservicemesh/sdk-vppagent

go 1.15

require (
	github.com/edwarnicke/exechelper v1.0.2
	github.com/edwarnicke/serialize v1.0.0
	github.com/golang/protobuf v1.4.3
	github.com/networkservicemesh/api v0.0.0-20201108204718-89d65b3605cf
	github.com/networkservicemesh/sdk v0.0.0-20201106151841-39537ac8948d
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	go.ligato.io/vpp-agent/v3 v3.2.0
	go.uber.org/goleak v1.1.10
	golang.org/x/net v0.0.0-20201010224723-4f7140c49acb // indirect
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20200609130330-bd2cb7843e1b
	google.golang.org/genproto v0.0.0-20201014134559-03b6142f0dc9 // indirect
	google.golang.org/grpc v1.33.2
)

replace (
	github.com/networkservicemesh/api v0.0.0-20201108204718-89d65b3605cf => ../api
	github.com/networkservicemesh/sdk v0.0.0-20201106151841-39537ac8948d => ../sdk
)
