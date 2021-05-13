package web

import (
	"text/template"

	"github.com/bennycio/bundle/api"
)

type Profile struct {
	Id       string
	Username string
	Email    string
	Tag      string
	Bundles  []string
}

type TemplateData struct {
	Profile Profile
	Plugin  *api.Plugin
	Plugins []*api.Plugin
	Page    int
	Math    func(int, int, string) int
	Date    func(int64) string
	Readme  string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("assets/templates/*.gohtml"))
}

func userToProfile(u *api.User) Profile {
	return Profile{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Tag:      u.Tag,
	}
}
