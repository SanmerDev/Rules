package clash

const (
	RuleConfigDomain        RuleConfig = "DOMAIN"
	RuleConfigDomainSuffix  RuleConfig = "DOMAIN-SUFFIX"
	RuleConfigDomainKeyword RuleConfig = "DOMAIN-KEYWORD"
	RuleConfigIPCIDR        RuleConfig = "IP-CIDR"
	RuleConfigIPCIDR6       RuleConfig = "IP-CIDR6"
	RuleConfigSrcIPCIDR     RuleConfig = "SRC-IP-CIDR"
	RuleConfigSrcPort       RuleConfig = "SRC-PORT"
	RuleConfigDstPort       RuleConfig = "DST-PORT"
	RuleConfigProcessName   RuleConfig = "PROCESS-NAME"
	RuleConfigProcessPath   RuleConfig = "PROCESS-PATH"
)

type RuleConfig string

const (
	TypeDomain RuleType = iota
	TypeDomainSuffix
	TypeDomainKeyword
	TypeIPCIDR
	TypeSrcIPCIDR
	TypeSrcPort
	TypeDstPort
	TypeProcess
	TypeProcessPath
)

type RuleType int

type Rule interface {
	RuleType() RuleType
	Adapter() string
	Payload() string
	String() string
}
