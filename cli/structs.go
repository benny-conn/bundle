package cli

type BundleFile struct {
	Plugins map[string]string `yaml:"Plugins,omitempty"`
}
