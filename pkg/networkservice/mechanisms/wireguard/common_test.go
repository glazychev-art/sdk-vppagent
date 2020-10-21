package wireguard_test

import (
	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/cls"
	wireguard_mechanism "github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/wireguard"
	"strconv"
)

type wireguardTestParams struct {
	srcIP         string
	dstIP         string
	srcPort       int
	dstPort       int
}

func (p *wireguardTestParams) ToMap() map[string]string {
	return map[string]string{
		wireguard_mechanism.SrcIP: p.srcIP,
		wireguard_mechanism.DstIP: p.dstIP,
		wireguard_mechanism.SrcPort: strconv.Itoa(p.srcPort),
		wireguard_mechanism.DstPort: strconv.Itoa(p.dstPort),
	}
}

func configureTestWireguardParameters() wireguardTestParams {
	return wireguardTestParams{
		srcIP:         "1.1.1.1",
		dstIP:         "1.1.1.2",
		srcPort:       58000,
		dstPort:       58001,
	}
}

func configureTestWireguardMechanism(parameters map[string]string) *networkservice.Mechanism {
	return &networkservice.Mechanism{
		Cls:        cls.REMOTE,
		Type:       wireguard_mechanism.MECHANISM,
		Parameters: parameters,
	}
}
