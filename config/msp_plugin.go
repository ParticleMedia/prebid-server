package config

// General config for all plugins. Different plugins may have different configs structure, but
// two fields are mandatory for all of them:
//  1. `so_path`, a string value telling Prebid-server where to load the shared object.
//  2. `enabled`, a boolean value telling whether or not to load it.
type PluginConfig struct {
	SoPath  string `mapstructure:"so_path" json:"so_path"`
	Enabled bool   `mapstructure:"enabled" json:"enabled"`
}
