package clash

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
)

type RawRuleSet struct {
	Payload []string `yaml:"payload"`
}

type RuleSet []Rule

func EncodeRuleSet(rules RuleSet, writer io.Writer) error {
	rawRules := make([]string, len(rules))
	for i, rule := range rules {
		rawRules[i] = rule.String()
	}

	raw := &RawRuleSet{Payload: rawRules}
	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)
	err := encoder.Encode(raw)
	if err != nil {
		_ = encoder.Close()
		return err
	}

	return encoder.Close()
}

func DecodeRuleSet(reader io.Reader) (RuleSet, error) {
	raw := &RawRuleSet{}
	decoder := yaml.NewDecoder(reader)
	err := decoder.Decode(raw)
	if err != nil {
		return nil, err
	}

	return parseRules(raw.Payload)
}

func parseRules(rawRules []string) (RuleSet, error) {
	var rules []Rule
	for i, line := range rawRules {
		rule := trimArr(strings.Split(line, ","))
		var (
			payload string
			target  string
			params  []string
		)

		switch length := len(rule); {
		case length == 2:
			payload = rule[1]
		case length == 3:
			payload = rule[1]
			target = rule[2]
		case length >= 4:
			payload = rule[1]
			target = rule[2]
			params = rule[3:]
		default:
			return nil, fmt.Errorf("rules[%d] [%s] error: format invalid", i, line)
		}

		params = trimArr(params)
		parsed, err := parseRule(rule[0], payload, target, params)
		if err != nil {
			return nil, fmt.Errorf("rules[%d] [%s] error: %s", i, line, err.Error())
		}

		rules = append(rules, parsed)
	}

	return rules, nil
}

func parseRule(tp, payload, target string, params []string) (Rule, error) {
	var (
		err  error
		rule Rule
	)

	switch RuleConfig(tp) {
	case RuleConfigDomain:
		rule = NewDomain(payload, target)
	case RuleConfigDomainSuffix:
		rule = NewDomainSuffix(payload, target)
	case RuleConfigDomainKeyword:
		rule = NewDomainKeyword(payload, target)
	case RuleConfigIPCIDR, RuleConfigIPCIDR6:
		noResolve := HasNoResolve(params)
		rule, err = NewIPCIDR(payload, target, false, noResolve)
	case RuleConfigSrcIPCIDR:
		rule, err = NewIPCIDR(payload, target, true, true)
	case RuleConfigSrcPort:
		rule, err = NewPort(payload, target, PortTypeSrc)
	case RuleConfigDstPort:
		rule, err = NewPort(payload, target, PortTypeDest)
	case RuleConfigProcessName:
		rule = NewProcess(payload, target, true)
	case RuleConfigProcessPath:
		rule = NewProcess(payload, target, false)
	default:
		err = fmt.Errorf("unsupported rule type %s", tp)
	}

	return rule, err
}
