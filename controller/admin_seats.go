package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tribalmedia/vista/model"
	"github.com/tribalmedia/vista/service"
	"github.com/tribalmedia/vista/setting"
)

// SeatsList is ...
func SeatsList(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	templateData := map[string]interface{}{
		"title":          "Seats List",
		"seatMasterList": model.GetAllSeatMaster(setting.SeatMasterAllStatus),
		"auth":           auth,
		"tab":            setting.SeatsTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_seats/seats_list.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// SeatsAdd is ...
func SeatsAdd(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	templateData := map[string]interface{}{
		"title":       "Seats Add",
		"companyList": model.GetAllCompany(),
		"auth":        auth,
		"tab":         setting.SeatsTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_seats/seats_add.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// AddSeats is ...
func AddSeats(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	seats := model.SeatMaster{}
	seats.Title = r.FormValue("title")
	seats.CompanyID, _ = strconv.Atoi(r.FormValue("company"))
	seats.SeatMaster = r.FormValue("seats")
	model.AddSeats(seats)
}

// SeatsDetail is ...
func SeatsDetail(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	vars := mux.Vars(r)
	seatsID, _ := strconv.Atoi(vars["id"])
	seatMaster := model.GetSeatsByID(seatsID)
	company := model.GetCompanyByID(seatMaster.CompanyID).Name
	templateData := map[string]interface{}{
		"title":      "Seats Detail",
		"seatMaster": seatMaster,
		"company":    company,
		"auth":       auth,
		"tab":        setting.SeatsTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_seats/seats_detail.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// SeatsEdit is ...
func SeatsEdit(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	vars := mux.Vars(r)
	seatsID, _ := strconv.Atoi(vars["id"])
	seatMaster := model.GetSeatsByID(seatsID)
	templateData := map[string]interface{}{
		"title":       "Seats Edit",
		"seatMaster":  seatMaster,
		"companyList": model.GetAllCompany(),
		"auth":        auth,
		"tab":         setting.SeatsTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_seats/seats_edit.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// EditSeats is ...
func EditSeats(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	seats := model.SeatMaster{}
	seats.ID, _ = strconv.Atoi(r.FormValue("id"))
	seats.Title = r.FormValue("title")
	seats.CompanyID, _ = strconv.Atoi(r.FormValue("company"))
	seats.Status, _ = strconv.Atoi(r.FormValue("status"))
	seats.SeatMaster = r.FormValue("seats")
	model.EditSeats(seats)
	w.Write([]byte("1"))
}

// SeatsDelete is ...
func SeatsDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	seatMasterID, _ := strconv.Atoi(id)

	model.DeleteSeatsByID(seatMasterID)
	http.Redirect(w, r, "/admin/seats", 301)
}

// SeatsAjaxCall is ...
func SeatsAjaxCall(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, _ := strconv.Atoi((r.FormValue("id")))
	seats := model.GetSeatsByID(id)
	b, _ := json.Marshal(seats)
	w.Write([]byte(b))
}
