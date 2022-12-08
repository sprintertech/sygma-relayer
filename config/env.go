// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type wrapper struct {
	Config RawConfig `json:"syg"`
}

const EnvPrefix = "SYG"

func loadFromEnv() (RawConfig, error) {
	// load relayer config
	jsonRelayerConfig, err := loadENVToJsonStructure()
	if err != nil {
		return RawConfig{}, err
	}
	c := &wrapper{}
	err = json.Unmarshal(jsonRelayerConfig, c)
	if err != nil {
		return RawConfig{}, err
	}
	rawConfig := c.Config

	// load chain configs
	index := 1
	for {
		rawDomainConfig := os.Getenv(fmt.Sprintf("%s_DOM_%d", EnvPrefix, index))
		if rawDomainConfig == "" {
			break
		}
		var cc map[string]interface{}
		err = json.Unmarshal([]byte(rawDomainConfig), &cc)
		if err != nil {
			return RawConfig{}, err
		}
		rawConfig.ChainConfigs = append(rawConfig.ChainConfigs, cc)
		index++
	}

	err = rawConfig.RelayerConfig.Validate()
	if err != nil {
		return RawConfig{}, err
	}

	return rawConfig, nil
}

func loadENVToJsonStructure() ([]byte, error) {
	structure := map[string]interface{}{}
	for _, e := range os.Environ() {
		if strings.Contains(e, EnvPrefix) {
			pair := strings.SplitN(e, "=", 2)
			indexes := strings.Split(pair[0], "_")
			mountMap(structure, indexes, pair[1])
		}
	}
	return json.MarshalIndent(structure, "", "    ")
}

func mountMap(m map[string]interface{}, i []string, v interface{}) {
	if len(i) > 1 {
		if _, ok := m[i[0]]; !ok {
			m[i[0]] = map[string]interface{}{}
		}
		asMap, ok := m[i[0]].(map[string]interface{})
		if !ok {
			return
		}
		mountMap(asMap, i[1:], v)
		v = asMap
	}
	m[i[0]] = v
}
