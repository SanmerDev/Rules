package clash

import (
	"strings"
)

type DomainSuffix struct {
	suffix  string
	adapter string
}

func (d *DomainSuffix) RuleType() RuleType {
	return TypeDomainSuffix
}

func (d *DomainSuffix) Adapter() string {
	return d.adapter
}

func (d *DomainSuffix) Payload() string {
	return d.suffix
}

func (d *DomainSuffix) String() string {
	builder := strings.Builder{}
	builder.WriteString(string(RuleConfigDomainSuffix))
	builder.WriteString(",")
	builder.WriteString(d.Payload())
	if len(d.adapter) != 0 {
		builder.WriteString(",")
		builder.WriteString(d.adapter)
	}
	return builder.String()
}

func NewDomainSuffix(suffix string, adapter string) *DomainSuffix {
	return &DomainSuffix{
		suffix:  strings.ToLower(suffix),
		adapter: adapter,
	}
}
