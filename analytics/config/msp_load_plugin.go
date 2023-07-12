package config

import (
	"encoding/json"
	"fmt"
	"plugin"

	"github.com/golang/glog"
	"github.com/prebid/prebid-server/analytics"
	"github.com/prebid/prebid-server/config"
)

type Builder interface {
	Build(json.RawMessage) (analytics.PBSAnalyticsModule, error)
}

func mspGetCustomAnalyticsAdapters(cfg map[string]interface{}) []analytics.PBSAnalyticsModule {
	plugins := make([]analytics.PBSAnalyticsModule, 0)

	for name, cfgData := range cfg {
		plugin, err := loadPlugin(name, cfgData)
		if err != nil {
			message := fmt.Sprintf("Failed to build analytics adapter %s, error: %+v\n", name, err)
			panic(message)
		} else {
			glog.Infof("Loaded analytics adapter: %s\n", name)
			plugins = append(plugins, plugin)
		}
	}

	return plugins
}

func loadPlugin(name string, cfgData interface{}) (analytics.PBSAnalyticsModule, error) {
	cfg, cfgJson := parsePluginConfig(name, cfgData)
	if cfg.SoPath == "" {
		message := fmt.Sprintf("The path to load analytics adapter %s is empty.\n", name)
		panic(message)
	}

	p, err := plugin.Open(cfg.SoPath)
	if err != nil {
		message := fmt.Sprintf("Failed to open shared object of analytics adapter %s, err: %+v.\n", name, err)
		panic(message)
	}

	s, err := p.Lookup("Builder")
	if err != nil {
		message := fmt.Sprintf("Failed to find Builder from analytics adapter %s, err: %+v.\n", name, err)
		panic(message)
	}

	builder, ok := s.(Builder)
	if !ok {
		message := fmt.Sprintf("Failed to convert Builder from analytics adapter %s, err: %+v.\n", name, err)
		panic(message)
	}

	return builder.Build(cfgJson)
}

func parsePluginConfig(name string, cfgData interface{}) (config.PluginConfig, json.RawMessage) {
	data, err := json.Marshal(cfgData)
	if err != nil {
		message := fmt.Sprintf("Failed to marshal config of plugin %s, err: %+v\n", name, err)
		panic(message)
	}

	var cfg config.PluginConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		message := fmt.Sprintf("Config of plugin %s is invalid , err: %+v\n", name, err)
		panic(message)
	}

	return cfg, data
}
