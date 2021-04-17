package cli

type BundleFile struct {
	Plugins map[string]string `yaml:"Plugins,omitempty"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PluginYML struct {
	Name        string `yml:"name"`
	Version     string `yml:"version"`
	Description string `yml:"description,omitempty`
}
