package box

import (
	"github.com/SanmerDev/rules/clash"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func domain(value string) clash.Rule {
	return clash.NewDomain(value, "")
}

func domainKeyword(value string) clash.Rule {
	return clash.NewDomainKeyword(value, "")
}

func domainSuffix(value string) clash.Rule {
	return clash.NewDomainSuffix(value, "")
}

func srcIPCIDR(value string) (clash.Rule, error) {
	return clash.NewIPCIDR(value, "", true, true)
}

func destIPCIDR(value string) (clash.Rule, error) {
	return clash.NewIPCIDR(value, "", false, true)
}

func srcPort(value uint16) clash.Rule {
	return clash.NewUPort(value, "", clash.PortTypeSrc)
}

func dstPort(value uint16) clash.Rule {
	return clash.NewUPort(value, "", clash.PortTypeDest)
}

func processName(value string) clash.Rule {
	return clash.NewProcess(value, "", true)
}

func processPath(value string) clash.Rule {
	return clash.NewProcess(value, "", false)
}

func static[T any](raw option.Listable[T], new func(T) clash.Rule) clash.RuleSet {
	rules := make([]clash.Rule, len(raw))
	for i, r := range raw {
		rules[i] = new(r)
	}

	return rules
}

func dynamic[T any](raw option.Listable[T], new func(T) (clash.Rule, error)) clash.RuleSet {
	var rules []clash.Rule
	for _, r := range raw {
		v, _ := new(r)
		if v != nil {
			rules = append(rules, v)
		}
	}

	return rules
}

func ToClash(rules []option.HeadlessRule) clash.RuleSet {
	var clashRules []clash.Rule
	for _, rule := range rules {
		if rule.Type != C.RuleTypeDefault {
			continue
		}
		clashRules = append(clashRules, static(rule.DefaultOptions.Domain, domain)...)
		clashRules = append(clashRules, static(rule.DefaultOptions.DomainKeyword, domainKeyword)...)
		clashRules = append(clashRules, static(rule.DefaultOptions.DomainSuffix, domainSuffix)...)
		clashRules = append(clashRules, dynamic(rule.DefaultOptions.SourceIPCIDR, srcIPCIDR)...)
		clashRules = append(clashRules, dynamic(rule.DefaultOptions.IPCIDR, destIPCIDR)...)
		clashRules = append(clashRules, static(rule.DefaultOptions.SourcePort, srcPort)...)
		clashRules = append(clashRules, static(rule.DefaultOptions.Port, dstPort)...)
		clashRules = append(clashRules, static(rule.DefaultOptions.ProcessName, processName)...)
		clashRules = append(clashRules, static(rule.DefaultOptions.ProcessPath, processPath)...)
	}
	return clashRules
}
