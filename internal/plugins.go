package internal

import (
	"go.mongodb.org/mongo-driver/bson"
)

func GetPluginByName(name string) (*Plugin, error) {

	session, err := GetMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("plugins")
	decodedPluginResult := &Plugin{}

	err = collection.FindOne(session.Ctx, bson.D{{"plugin", NewCaseInsensitiveRegex(name)}}).Decode(decodedPluginResult)
	if err != nil {
		return nil, err
	}

	return decodedPluginResult, nil
}

func GetPluginAuthor(pluginName string) (string, error) {

	plugin, err := GetPluginByName(pluginName)
	if err != nil {
		return "", err
	}
	return plugin.User, nil
}

func GetPluginVersion(pluginName string) (string, error) {
	plugin, err := GetPluginByName(pluginName)
	if err != nil {
		return "", err
	}
	return plugin.Version, nil
}
