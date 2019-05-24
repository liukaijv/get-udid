package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	IpaFile          string `json:"ipa_file"`
	PlistFile        string `json:"plist_file"`
	MobileConfigFile string `json:"mobile_config_file"`
	HostPrefixDir    string `json:"host_prefix_dir"`
}

func New(file string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
