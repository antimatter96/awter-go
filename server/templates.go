package server

import "html/template"

func (server *Server) parseTemplates(templatePath string) {
	server.shortnerTemplate = template.Must(template.ParseFiles(templatePath + "/shortner.html"))
	server.createdTemplate = template.Must(template.ParseFiles(templatePath + "/created.html"))
	server.elongateTemplate = template.Must(template.ParseFiles(templatePath + "/elongate.html"))
}
