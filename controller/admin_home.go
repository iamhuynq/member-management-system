package controller

import (
	"html/template"
	"net/http"

	"github.com/tribalmedia/vista/service"
	"github.com/tribalmedia/vista/setting"
)

// AdminHome page is ...
func AdminHome(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	templateData := map[string]interface{}{
		"title": "Admin",
		"auth":  auth,
		"tab":   setting.AdminHomeTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_home/admin_home.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}
