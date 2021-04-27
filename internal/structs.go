package internal

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Scopes    []string `json:"scopes"`
	Plugins   []string `json:"plugins"`
	Tag       string   `json:"tag"`
	Thumbnail []byte
}

type TemplateData struct {
	User     User
	Plugin   Plugin
	Plugins  []Plugin
	TestData string
}

type Plugin struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Readme      string
	Thumbnail   []byte
	LastUpdated int64 `json:"lastUpdated"`
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
