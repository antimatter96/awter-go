package server

import "html/template"

func (server *Server) parseTemplates(templatePath string) {
	server.shortnerTemplate = template.Must(template.ParseFiles(
		[]string{
			templatePath + "/shortner.html",
			templatePath + "/_heading.html",
			templatePath + "/_showHidePassword.html",
		}...,
	))
	server.createdTemplate = template.Must(template.ParseFiles(
		[]string{
			templatePath + "/created.html",
			templatePath + "/_heading.html",
			templatePath + "/_showHidePassword.html",
		}...,
	))
	server.elongateTemplate = template.Must(template.ParseFiles(
		[]string{
			templatePath + "/elongate.html",
			templatePath + "/_heading.html",
			templatePath + "/_showHidePassword.html",
		}...,
	))
}
