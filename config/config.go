// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/creasty/defaults"
	"github.com/imdario/mergo"

	"github.com/ChainSafe/sygma-relayer/config/relayer"
	"github.com/spf13/viper"
)

type Config struct {
	RelayerConfig relayer.RelayerConfig
	ChainConfigs  []map[string]interface{}
}

type RawConfig struct {
	RelayerConfig relayer.RawRelayerConfig `mapstructure:"relayer" json:"relayer"`
	ChainConfigs  []map[string]interface{} `mapstructure:"chains" json:"domains"`
}

// GetConfigFromENV reads config from ENV variables, validates it and parses
// it into config suitable for application
//
//
// Properties of RelayerConfig are expected to be defined as separate ENV variables
// where ENV variable name reflects properties position in structure. Each ENV variable needs to be prefixed with CBH.
//
// For example, if you want to set Config.RelayerConfig.MpcConfig.Port this would
// translate to ENV variable named CBH_RELAYER_MPCCONFIG_PORT.
//
//
// Each ChainConfig is defined as one ENV variable, where its content is JSON configuration for one chain/domain.
// Variables are named like this: CBH_DOM_X where X is domain id.
//
func GetConfigFromENV(config *Config) (*Config, error) {
	rawConfig, err := loadFromEnv()
	if err != nil {
		return config, err
	}

	return processRawConfig(rawConfig, config)
}

// GetConfigFromFile reads config from file, validates it and parses
// it into config suitable for application
func GetConfigFromFile(path string, config *Config) (*Config, error) {
	rawConfig := RawConfig{}

	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&rawConfig)
	if err != nil {
		return config, err
	}

	return processRawConfig(rawConfig, config)
}

// GetConfigFromNetwork fetches shared configuration from URL and parses it.
func GetConfigFromNetwork(url string, config *Config) (*Config, error) {
	rawConfig := RawConfig{}

	resp, err := http.Get(url)
	if err != nil {
		return &Config{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Config{}, err
	}

	err = json.Unmarshal(body, &rawConfig)
	if err != nil {
		return &Config{}, err
	}

	config.ChainConfigs = rawConfig.ChainConfigs
	return config, err
}

func processRawConfig(rawConfig RawConfig, config *Config) (*Config, error) {
	if err := defaults.Set(&rawConfig); err != nil {
		return config, err
	}

	relayerConfig, err := relayer.NewRelayerConfig(rawConfig.RelayerConfig)
	if err != nil {
		return config, err
	}

	for i, chain := range rawConfig.ChainConfigs {
		mergo.Merge(&chain, config.ChainConfigs[i])
		if chain["type"] == "" || chain["type"] == nil {
			return config, fmt.Errorf("chain 'type' must be provided for every configured chain")
		}

		config.ChainConfigs[i] = chain
	}

	config.RelayerConfig = relayerConfig
	return config, nil
}
