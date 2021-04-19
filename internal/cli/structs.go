package cli

type BundleFile struct {
	Plugins map[string]string `yaml:"Plugins,omitempty"`
}

type PluginYML struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description,omitempty"`
}
