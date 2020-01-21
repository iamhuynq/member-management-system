package main

import (
	"encoding/gob"
	"net/http"

	"github.com/tribalmedia/vista/setting"

	"github.com/gorilla/mux"
	"github.com/tribalmedia/vista/controller"
	"github.com/tribalmedia/vista/model"
	"github.com/tribalmedia/vista/service"

	_ "github.com/go-sql-driver/mysql"
)

// getSessionMiddleware ...
func getSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store := service.ConnectToRedis()
		defer store.Close()
		//get session info
		session, _ := store.Get(r, "vista_member")
		if session.Values["member"] == nil {
			// redirect login page
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// getSessionMiddleware ...
func adminRoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		memberLogin := service.GetSessionMember(r)
		if memberLogin.RoleType == setting.AdminRoleType {
			next.ServeHTTP(w, r)
		} else {
			// redirect to home page
			http.Redirect(w, r, "/", http.StatusFound)
		}
	})
}

func main() {
	// connect database before do anything
	logger := setting.InitLogger()
	model.Logger = logger
	controller.Logger = logger
	service.Logger = logger

	model.ConnectDb()
	gob.Register(model.Auth{})

	fs := http.FileServer(http.Dir("webroot"))
	http.Handle("/webroot/", http.StripPrefix("/webroot/", fs))

	r := mux.NewRouter().StrictSlash(true)
	// This is for nomal user
	t := r.PathPrefix("/").Subrouter()
	// This is for admin
	s := t.PathPrefix("/admin").Subrouter()

	s.HandleFunc("/members", controller.MembersList)
	s.HandleFunc("/members/{id}/delete", controller.MembersDelete).Methods("POST")
	s.HandleFunc("/members/add", controller.MembersAdd).Methods("POST")
	s.HandleFunc("/members/add", controller.MembersAddDisplay).Methods("GET")
	s.HandleFunc("/members/{id}", controller.MembersDetail)
	s.HandleFunc("/members/{id}/edit", controller.MembersEdit).Methods("POST")
	s.HandleFunc("/members/{id}/edit", controller.MembersEditDisplay).Methods("GET")
	s.HandleFunc("/members/{id}/reset_password", controller.ResetPassword).Methods("POST")
	s.HandleFunc("/members/{id}/reset_password", controller.ResetPasswordDisplay).Methods("GET")

	t.HandleFunc("/team", controller.TeamsTop).Methods("GET")
	s.HandleFunc("/moveMember", controller.MoveMember).Methods("POST")

	s.HandleFunc("/teams", controller.TeamsList)
	s.HandleFunc("/teams/add", controller.TeamsAdd).Methods("POST")
	s.HandleFunc("/teams/add", controller.TeamsAddDisplay).Methods("GET")
	s.HandleFunc("/teams/{id}", controller.TeamsDetail)
	s.HandleFunc("/teams/{id}/edit", controller.TeamsEdit).Methods("POST")
	s.HandleFunc("/teams/{id}/edit", controller.TeamsEditDisplay).Methods("GET")

	s.HandleFunc("/companies", controller.CompanyList)
	s.HandleFunc("/companies/add", controller.CompanyAddDisplay).Methods("GET")
	s.HandleFunc("/companies/add", controller.CompanyAdd).Methods("POST")
	s.HandleFunc("/companies/{id}", controller.CompanyDetail).Methods("GET")
	s.HandleFunc("/companies/{id}/edit", controller.CompanyEditDisplay).Methods("GET")
	s.HandleFunc("/companies/{id}/edit", controller.CompanyEdit).Methods("POST")

	s.HandleFunc("/seats", controller.SeatsList).Methods("GET")
	s.HandleFunc("/seats/add", controller.SeatsAdd).Methods("GET")
	s.HandleFunc("/seats/add", controller.AddSeats).Methods("POST")
	s.HandleFunc("/seats/{id}", controller.SeatsDetail).Methods("GET")
	s.HandleFunc("/seats/{id}/delete", controller.SeatsDelete).Methods("POST")
	s.HandleFunc("/seats/edit", controller.EditSeats).Methods("POST")
	s.HandleFunc("/seats/{id}/edit", controller.SeatsEdit).Methods("GET")
	t.HandleFunc("/seats/getSeats", controller.SeatsAjaxCall).Methods("POST")

	t.HandleFunc("/seat", controller.SeatsTop).Methods("GET")
	s.HandleFunc("/addSeat", controller.AddSeat).Methods("POST")
	s.HandleFunc("/getSeat", controller.GetSeat).Methods("POST")
	s.HandleFunc("/editSeat", controller.EditSeat).Methods("POST")
	s.HandleFunc("/deleteSeat", controller.DeleteSeat).Methods("POST")
	t.HandleFunc("/getAllSeat", controller.GetAllSeat).Methods("POST")

	t.HandleFunc("/", controller.Home).Methods("GET")
	s.HandleFunc("", controller.AdminHome).Methods("GET")
	t.HandleFunc("/ws", controller.HandleConnectionsSocket)
	t.HandleFunc("/my_page", controller.MyPageDisplay).Methods("GET")
	t.HandleFunc("/my_page", controller.MyPage).Methods("POST")
	t.HandleFunc("/members/change_password", controller.ChangePassword).Methods("POST")
	t.HandleFunc("/members/change_password", controller.ChangePasswordDisplay).Methods("GET")

	// Start listening for incoming messages
	go controller.HandleDataTeam()
	go controller.HandleDataSeat()
	go controller.ShowUserOnline()

	r.HandleFunc("/login", controller.LoginView).Methods("GET")
	r.HandleFunc("/login", controller.Login).Methods("POST")
	r.HandleFunc("/logout", controller.Logout).Methods("GET")

	http.Handle("/", r)
	t.Use(getSessionMiddleware)
	s.Use(adminRoleMiddleware)

	http.ListenAndServe(":8080", nil)
}
