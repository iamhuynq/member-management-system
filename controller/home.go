package controller

import (
	"html/template"
	"net/http"

	"github.com/tribalmedia/vista/service"
	"github.com/tribalmedia/vista/setting"
)

// Home page is ...
func Home(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	templateData := map[string]interface{}{
		"title": "Home",
		"auth":  auth,
	}

	tmpl := template.Must(template.ParseFiles("template/home/home.tmpl", setting.UserTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}
