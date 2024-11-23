package main

import (
	"encoding/json"
	"fmt"
	"github.com/SanmerDev/rules/box"
	"github.com/SanmerDev/rules/clash"
	"github.com/SanmerDev/rules/config"
	"github.com/sagernet/sing-box/common/srs"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"net/http"
	"os"
	"sync"
)

func boxToClash(rule *config.Rule, group *sync.WaitGroup) {
	fmt.Println(rule.Tag + ": " + rule.Url)
	response, err := http.Get(rule.Url)
	if err != nil {
		fmt.Println(err)
		return
	}

	if response.StatusCode != http.StatusOK {
		fmt.Printf("status: %s\n", response.Status)
		_ = response.Body.Close()
		return
	}

	var rules []option.HeadlessRule
	switch rule.Format {
	case C.RuleSetFormatSource:
		raw := option.PlainRuleSetCompat{}
		decoder := json.NewDecoder(response.Body)
		err := decoder.Decode(&raw)
		if err != nil {
			fmt.Println(err)
			return
		}
		rules = raw.Options.Rules
	case C.RuleSetFormatBinary:
		plainRuleSet, err := srs.Read(response.Body, true)
		if err != nil {
			fmt.Println(err)
			return
		}
		rules = plainRuleSet.Options.Rules
	default:
		rules = []option.HeadlessRule{}
	}

	file, err := os.Create(rule.Tag + ".yaml")
	if err != nil {
		fmt.Println(err)
		_ = file.Close()
		return
	}

	err = clash.EncodeRuleSet(box.ToClash(rules, rule.NoResolve), file)
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = response.Body.Close()
	_ = file.Close()
	group.Done()
}

func main() {
	reader, err := os.Open("rules.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	rules, err := config.DecodeRules(reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	var group sync.WaitGroup
	for _, rule := range rules {
		group.Add(1)
		go boxToClash(&rule, &group)
	}

	group.Wait()
}
