package shortner

import (
	"html/template"
)

var shortnerTemplate *template.Template
var createdTemplate *template.Template
var elongateTemplate *template.Template

func parseTemplates() {
	shortnerTemplate = template.Must(template.ParseFiles("./template/shortner.html"))
	createdTemplate = template.Must(template.ParseFiles("./template/created.html"))
	elongateTemplate = template.Must(template.ParseFiles("./template/elongate.html"))
}
