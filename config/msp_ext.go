package config

// General config for all plugins. Different plugins may have different configs structure, but
// `so_path` is a mandatory field for all of their configs. This field tells Prebid-server where
// to load the shared object.
type PluginConfig struct {
	SoPath string `mapstructure:"so_path" json:"so_path"`
}
