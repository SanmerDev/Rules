package clash

import (
	"net"
	"strings"
)

type IPCIDR struct {
	ipNet       *net.IPNet
	adapter     string
	isSourceIP  bool
	noResolveIP bool
}

func (i *IPCIDR) RuleType() RuleType {
	if i.isSourceIP {
		return TypeSrcIPCIDR
	}
	return TypeIPCIDR
}

func (i *IPCIDR) Adapter() string {
	return i.adapter
}

func (i *IPCIDR) Payload() string {
	return i.ipNet.String()
}

func (i *IPCIDR) String() string {
	builder := strings.Builder{}
	if i.isSourceIP {
		builder.WriteString(string(RuleConfigSrcIPCIDR))
	} else {
		switch len(i.ipNet.IP) {
		case net.IPv4len:
			builder.WriteString(string(RuleConfigIPCIDR))
		case net.IPv6len:
			builder.WriteString(string(RuleConfigIPCIDR6))
		}
	}
	builder.WriteString(",")
	builder.WriteString(i.Payload())
	if len(i.adapter) != 0 {
		builder.WriteString(",")
		builder.WriteString(i.adapter)
	}
	if i.noResolveIP {
		builder.WriteString(",")
		WriteNoResolve(&builder)
	}
	return builder.String()
}

func NewIPCIDR(s string, adapter string, source bool, noResolve bool) (*IPCIDR, error) {
	_, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, errPayload
	}

	ipcidr := &IPCIDR{
		ipNet:       ipNet,
		adapter:     adapter,
		isSourceIP:  source,
		noResolveIP: noResolve,
	}

	return ipcidr, nil
}
