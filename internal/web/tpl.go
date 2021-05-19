package web

import (
	"text/template"

	"github.com/bennycio/bundle/api"
)

type profile struct {
	Id         string
	Username   string
	Email      string
	Tag        string
	StripeInfo userStripeInfo
}

type userStripeInfo struct {
	Id               string
	DetailsSubmitted bool
	ChargesEnabled   bool
}

type templateData struct {
	Profile         profile
	Plugin          *api.Plugin
	Plugins         []*api.Plugin
	Bundle          *api.Bundle
	PurchaseSession string
	Page            int
	Math            func(int, int, string) int
	Date            func(int64) string
	Contains        func([]string, string) bool
	Readme          string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("assets/templates/*.gohtml"))
}

func userToProfile(u *api.User) profile {
	return profile{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Tag:      u.Tag,
	}
}
