package clash

import "strings"

type Process struct {
	adapter  string
	process  string
	nameOnly bool
}

func (p *Process) RuleType() RuleType {
	if p.nameOnly {
		return TypeProcess
	}
	return TypeProcessPath
}

func (p *Process) Adapter() string {
	return p.adapter
}

func (p *Process) Payload() string {
	return p.process
}

func (p *Process) String() string {
	builder := strings.Builder{}
	if p.nameOnly {
		builder.WriteString(string(RuleConfigProcessName))
	} else {
		builder.WriteString(string(RuleConfigProcessPath))
	}
	builder.WriteString(",")
	builder.WriteString(p.Payload())
	if len(p.adapter) != 0 {
		builder.WriteString(",")
		builder.WriteString(p.adapter)
	}
	return builder.String()
}

func NewProcess(process string, adapter string, nameOnly bool) *Process {
	return &Process{
		adapter:  adapter,
		process:  process,
		nameOnly: nameOnly,
	}
}
