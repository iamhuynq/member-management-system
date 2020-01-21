package controller

import (
	"html/template"
	"net/http"

	"github.com/tribalmedia/vista/model"
	"github.com/tribalmedia/vista/service"
)

// LoginView is ...
func LoginView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/auth/login.tmpl"))
	if err := tmpl.Execute(w, nil); err != nil {
		Logger.Error(err.Error())
	}
}

// Login is ...
func Login(w http.ResponseWriter, r *http.Request) {
	// get form value
	memberID, password := r.FormValue("loginID"), r.FormValue("password")
	templateData := map[string]interface{}{
		"loginID": memberID,
	}

	auth, err := model.Authenticated(memberID, password)

	if err == nil {
		store := service.ConnectToRedis()
		defer store.Close()
		session, _ := store.Get(r, "vista_member")
		session.Values["member"] = auth
		service.SaveSession(session, r, w)
		//redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		templateData["loginError"] = true
	}

	tmpl := template.Must(template.ParseFiles("template/auth/login.tmpl"))
	if err := tmpl.Execute(w, templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// Logout is ...
func Logout(w http.ResponseWriter, r *http.Request) {
	store := service.ConnectToRedis()
	defer store.Close()
	session, _ := store.Get(r, "vista_member")
	// Delete session
	session.Options.MaxAge = -1
	service.SaveSession(session, r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
