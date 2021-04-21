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
