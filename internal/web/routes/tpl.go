package routes

import (
	"text/template"

	bundle "github.com/bennycio/bundle/internal"
)

var tpl *template.Template

const ReqFileType = bundle.RequiredFileType

func init() {
	tpl = template.Must(template.ParseGlob("assets/templates/*.gohtml"))
}
