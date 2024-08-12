package clash

import (
	"strings"
)

type DomainKeyword struct {
	keyword string
	adapter string
}

func (d *DomainKeyword) RuleType() RuleType {
	return TypeDomainKeyword
}

func (d *DomainKeyword) Adapter() string {
	return d.adapter
}

func (d *DomainKeyword) Payload() string {
	return d.keyword
}

func (d *DomainKeyword) String() string {
	builder := strings.Builder{}
	builder.WriteString(string(RuleConfigDomainKeyword))
	builder.WriteString(",")
	builder.WriteString(d.Payload())
	if len(d.adapter) != 0 {
		builder.WriteString(",")
		builder.WriteString(d.adapter)
	}
	return builder.String()
}

func NewDomainKeyword(keyword string, adapter string) *DomainKeyword {
	return &DomainKeyword{
		keyword: strings.ToLower(keyword),
		adapter: adapter,
	}
}
