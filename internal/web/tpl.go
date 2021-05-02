package web

import (
	"text/template"

	"github.com/bennycio/bundle/api"
)

type TemplateData struct {
	User    *api.User
	Plugin  *api.Plugin
	Plugins []*api.Plugin
	Readme  string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("assets/templates/*.gohtml"))
}
