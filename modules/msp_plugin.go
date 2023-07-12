package modules

import (
	"encoding/json"
	"fmt"
	"plugin"

	"github.com/golang/glog"
	"github.com/prebid/prebid-server/config"
	"github.com/prebid/prebid-server/modules/moduledeps"
)

type ModuleBuilder interface {
	Build(json.RawMessage, moduledeps.ModuleDeps) (interface{}, error)
}

func mspLoadModulesThroughPlugin(modules map[string]interface{}, cfg config.Modules, deps moduledeps.ModuleDeps) map[string]interface{} {
	for vendor, moduleBuilders := range cfg {
		for moduleName, moduleCfg := range moduleBuilders {
			id := fmt.Sprintf("%s.%s", vendor, moduleName)

			if _, ok := modules[id]; ok {
				// skip loading modules that have already been loaded through hardcoded builder
				continue
			}

			cfg, cfgJson := parsePluginConfig(id, moduleCfg)
			if !cfg.Enabled {
				glog.Infof("Skip loading module %s as it is disabled.", id)
				continue
			}

			p, err := loadPlugin(id, cfg, cfgJson, deps)
			if err != nil {
				message := fmt.Sprintf("Failed to load module %s, err: %+v.\n", id, err)
				panic(message)
			}

			modules[id] = p
			glog.Infof("Loaded Module plugin %s.\n", id)
		}
	}

	return modules
}

func loadPlugin(name string, cfg config.PluginConfig, cfgJson json.RawMessage, deps moduledeps.ModuleDeps) (interface{}, error) {
	if cfg.SoPath == "" {
		message := fmt.Sprintf("The path to load module %s is empty.\n", name)
		panic(message)
	}

	p, err := plugin.Open(cfg.SoPath)
	if err != nil {
		message := fmt.Sprintf("Failed to open shared object of module %s, err: %+v.\n", name, err)
		panic(message)
	}

	s, err := p.Lookup("Builder")
	if err != nil {
		message := fmt.Sprintf("Failed to find Builder from module %s, err: %+v.\n", name, err)
		panic(message)
	}

	builder, ok := s.(ModuleBuilder)
	if !ok {
		message := fmt.Sprintf("Failed to convert Builder from module %s, err: %+v.\n", name, err)
		panic(message)
	}

	return builder.Build(cfgJson, deps)
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
