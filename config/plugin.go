package config

import (
	"fmt"

	"github.com/tidwall/gjson"
)

type PluginConfiguration struct {
	Allow     []string
	Deny      []string
	DenyAll   bool
	CacheSize int
}

func ParsPluginConfiguration(cByte []byte) (PluginConfiguration, error) {
	config := PluginConfiguration{
		DenyAll:   false,
		CacheSize: 200,
	}

	if !gjson.ValidBytes(cByte) {
		return PluginConfiguration{}, fmt.Errorf("the plugin configuration is not a valid json")
	}
	jsonData := gjson.ParseBytes(cByte)

	for _, v := range jsonData.Get("allow").Array() {
		config.Allow = append(config.Allow, v.Str)
	}
	for _, v := range jsonData.Get("deny").Array() {
		config.Deny = append(config.Deny, v.Str)
	}

	config.DenyAll = jsonData.Get("denyAll").Bool()

	if cacheSize := int(jsonData.Get("cacheSize").Int()); cacheSize != 0 {
		config.CacheSize = cacheSize
	}

	return config, nil
}
