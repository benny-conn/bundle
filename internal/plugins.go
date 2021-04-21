package internal

func GetPluginAuthor(pluginName string) (string, error) {

	plugin, err := GetPlugin(pluginName)
	if err != nil {
		return "", err
	}
	return plugin.User, nil
}

func GetPluginVersion(pluginName string) (string, error) {
	plugin, err := GetPlugin(pluginName)
	if err != nil {
		return "", err
	}
	return plugin.Version, nil
}
