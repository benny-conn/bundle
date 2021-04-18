package cli

type BundleFile struct {
	Plugins map[string]string `yaml:"Plugins,omitempty"`
}

type PluginYML struct {
	Name        string `yml:"name"`
	Version     string `yml:"version"`
	Description string `yml:"description,omitempty`
}
