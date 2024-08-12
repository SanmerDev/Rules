package config

import (
	"encoding/json"
	"io"
)

type Rule struct {
	Tag    string `yaml:"tag"`
	Format string `yaml:"format"`
	Url    string `yaml:"url"`
}

type Rules []Rule

func DecodeRules(reader io.Reader) (Rules, error) {
	var rules []Rule
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&rules)
	if err != nil {
		return nil, err
	}

	return rules, nil
}
