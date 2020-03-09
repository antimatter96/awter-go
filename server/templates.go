package server

import "text/template"

func (server *server) parseTemplates(templatePath string) {
	server.shortnerTemplate = template.Must(template.ParseFiles(templatePath + "/shortner.html"))
	server.createdTemplate = template.Must(template.ParseFiles(templatePath + "/created.html"))
	server.elongateTemplate = template.Must(template.ParseFiles(templatePath + "/elongate.html"))
}
