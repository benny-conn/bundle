package internal

import (
	"context"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateData struct {
	User     api.User
	Plugin   api.Plugin
	Plugins  []api.Plugin
	TestData string
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
