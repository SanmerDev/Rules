package clash

import (
	"strings"
)

type Domain struct {
	domain  string
	adapter string
}

func (d *Domain) RuleType() RuleType {
	return TypeDomain
}

func (d *Domain) Adapter() string {
	return d.adapter
}

func (d *Domain) Payload() string {
	return d.domain
}

func (d *Domain) String() string {
	builder := strings.Builder{}
	builder.WriteString(string(RuleConfigDomain))
	builder.WriteString(",")
	builder.WriteString(d.Payload())
	if len(d.adapter) != 0 {
		builder.WriteString(",")
		builder.WriteString(d.adapter)
	}
	return builder.String()
}

func NewDomain(domain string, adapter string) *Domain {
	return &Domain{
		domain:  strings.ToLower(domain),
		adapter: adapter,
	}
}
