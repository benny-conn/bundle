package internal

import (
	"context"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateData struct {
	Profile  Profile
	Plugin   api.Plugin
	Plugins  []api.Plugin
	TestData string
}

type Profile struct {
	Username string
	Email    string
	Tag      string
}

type Mongo struct {
	Client *mongo.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}

type PluginYML struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description,omitempty"`
}
