package clash

import (
	"fmt"
	"strconv"
	"strings"
)

type PortType int

const (
	PortTypeSrc PortType = iota
	PortTypeDest
)

type Port struct {
	port     uint16
	adapter  string
	portType PortType
}

func (p *Port) RuleType() RuleType {
	switch p.portType {
	case PortTypeSrc:
		return TypeSrcPort
	case PortTypeDest:
		return TypeDstPort
	default:
		panic(fmt.Errorf("unknown port type: %v", p.portType))
	}
}

func (p *Port) Adapter() string {
	return p.adapter
}

func (p *Port) Payload() string {
	return strconv.FormatUint(uint64(p.port), 10)
}

func (p *Port) String() string {
	builder := strings.Builder{}
	switch p.portType {
	case PortTypeSrc:
		builder.WriteString(string(RuleConfigSrcPort))
	case PortTypeDest:
		builder.WriteString(string(RuleConfigDstPort))
	}
	builder.WriteString(",")
	builder.WriteString(p.Payload())
	if len(p.adapter) != 0 {
		builder.WriteString(",")
		builder.WriteString(p.adapter)
	}
	return builder.String()
}

func NewPort(port string, adapter string, portType PortType) (*Port, error) {
	p, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil, errPayload
	}

	return &Port{
		port:     uint16(p),
		adapter:  adapter,
		portType: portType,
	}, nil
}

func NewUPort(port uint16, adapter string, portType PortType) *Port {
	return &Port{
		port:     port,
		adapter:  adapter,
		portType: portType,
	}
}
