package routes

import (
	"text/template"
)

type TemplateData struct {
	Profile Profile
	Plugin  PluginInfo
	Plugins []PluginInfo
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("assets/templates/*.gohtml"))
}
