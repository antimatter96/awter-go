package handlers

import (
	"html/template"
)

var shortnerTemplate *template.Template
var createdTemplate *template.Template

func parseTemplates() {
	shortnerTemplate = template.Must(template.ParseFiles("./template/shortner.html"))
	createdTemplate = template.Must(template.ParseFiles("./template/created.html"))
}
