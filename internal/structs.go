package internal

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Scopes   []string `json:"scopes"`
}

type Profile struct {
	Username string
	Email    string
	Tags     []string
	Scopes   []string
}

type TemplateData struct {
	Profile  Profile
	Plugin   Plugin
	Markdown string
}

type Plugin struct {
	Plugin  string `json:"plugin"`
	User    string `json:"user"`
	Version string `json:"version"`
}

type Mongo struct {
	Client *mongo.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}
