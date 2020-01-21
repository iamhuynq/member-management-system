package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/tribalmedia/vista/model"
	"github.com/tribalmedia/vista/service"
	"github.com/tribalmedia/vista/setting"
)

// MoveSeatData save position of member on seat page
type MoveSeatData struct {
	model.Seat
	DepartmentID int
	MoveType     int
}

// SeatsTop displays user screen for seat
func SeatsTop(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	var seats []map[string]interface{}
	if auth.RoleType == setting.AdminRoleType {
		seats = model.GetAllSeatMaster(setting.SeatMasterActiveStatus)
	}
	if auth.RoleType == setting.UserRoleType {
		seats = model.GetAllSeatsOfCompany(auth.CompanyID)
	}

	templateData := map[string]interface{}{
		"title":     "Seat",
		"seatsList": seats,
		"auth":      auth,
	}

	tmpl := template.Must(template.ParseFiles("template/seat/seat.tmpl", setting.UserTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// AddSeat adds member to seat called by ajax
func AddSeat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	seat := model.Seat{}
	seat.MemberID, _ = strconv.Atoi(r.FormValue("memberID"))
	seat.TableNumber, _ = strconv.Atoi(r.FormValue("tableNumber"))
	seat.SeatMasterID, _ = strconv.Atoi(r.FormValue("seatMasterID"))
	seat.Row, _ = strconv.Atoi(r.FormValue("row"))
	seat.Col, _ = strconv.Atoi(r.FormValue("column"))
	model.AddSeat(seat)
}

// GetAllSeat gets all seats from seats called by ajax
func GetAllSeat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	seatMasterID, _ := strconv.Atoi((r.FormValue("id")))
	seats := model.GetAllSeat(seatMasterID)
	b, _ := json.Marshal(seats)
	w.Write([]byte(b))
}

// GetSeat get seat from seats by seatMasterID and MemberID called by ajax
func GetSeat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	seatMasterID, _ := strconv.Atoi((r.FormValue("seatMasterID")))
	memberID, _ := strconv.Atoi((r.FormValue("memberID")))
	seats := model.GetSeat(seatMasterID, memberID)
	b, _ := json.Marshal(seats)
	w.Write([]byte(b))
}

// EditSeat edits seat called by ajax
func EditSeat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	seat := model.Seat{}
	seat.ID, _ = strconv.Atoi(r.FormValue("seatID"))
	seat.MemberID, _ = strconv.Atoi(r.FormValue("memberID"))
	seat.TableNumber, _ = strconv.Atoi(r.FormValue("tableNumber"))
	seat.SeatMasterID, _ = strconv.Atoi(r.FormValue("seatMasterID"))
	seat.Row, _ = strconv.Atoi(r.FormValue("row"))
	seat.Col, _ = strconv.Atoi(r.FormValue("column"))
	model.EditSeat(seat)
}

// DeleteSeat edits seat called by ajax
func DeleteSeat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	seatID, _ := strconv.Atoi(r.FormValue("seatID"))
	model.DeleteSeat(seatID)
}

//HandleDataSeat sends msg to clients connected use websocket
func HandleDataSeat() {
	for {
		// Grab the next message from the broadcast channel
		msd := <-broadcastSeat

		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msd)
			if err != nil {
				Logger.Error(err.Error())
				client.Close()
				delete(clients, client)
			}
		}

	}
}
