package model

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/tribalmedia/vista/setting"
)

// SeatMaster is ..
type SeatMaster struct {
	ID         int       `db:"id"`
	Title      string    `db:"title"`
	CompanyID  int       `db:"company_id"`
	SeatMaster string    `db:"seat_master"`
	Status     int       `db:"status"`
	Created    time.Time `db:"created"`
	Modified   time.Time `db:"modified"`
}

// GetAllSeatMaster is ...
func GetAllSeatMaster(status int) []map[string]interface{} {
	query := `SELECT seat_master.id, seat_master.title, seat_master.company_id, seat_master.seat_master,
			seat_master.status, seat_master.created, seat_master.modified FROM seat_master`
	if status == setting.SeatMasterActiveStatus {
		query += ` WHERE seat_master.status = ` + strconv.Itoa(status)
	}
	query += ` ORDER BY seat_master.id DESC`
	rows, _ := DB.Queryx(query)
	var listSeats []map[string]interface{}
	for rows.Next() {
		var seatMaster SeatMaster
		err := rows.StructScan(&seatMaster)
		if err != nil {
			Logger.Error(err.Error())
		}
		company := GetCompanyByID(seatMaster.CompanyID).Name
		departments := GetDepartmentsAndMemsByCompanyID(seatMaster.CompanyID)
		m := map[string]interface{}{
			"SeatMaster":      seatMaster,
			"Company":         company,
			"DepartmentsList": departments,
		}
		listSeats = append(listSeats, m)
	}
	return listSeats
}

// GetAllSeatsOfCompany gets all seats_master of user's company
// Id is id of user's company logined
func GetAllSeatsOfCompany(id int) []map[string]interface{} {
	var listSeats []map[string]interface{}
	rows, _ := DB.Queryx("SELECT id, title, company_id, status FROM seat_master WHERE company_id = ? AND status = ? ORDER BY id DESC", id, setting.SeatMasterActiveStatus)
	for rows.Next() {
		var seatMaster SeatMaster
		err := rows.StructScan(&seatMaster)
		if err != nil {
			Logger.Error(err.Error())
		}
		m := map[string]interface{}{
			"SeatMaster":      seatMaster,
			"DepartmentsList": GetDepartmentsAndMemsByCompanyID(id),
		}
		listSeats = append(listSeats, m)
	}

	return listSeats
}

// AddSeats is ...
func AddSeats(seatMaster SeatMaster) {
	query := `INSERT INTO seat_master (id, title, company_id, seat_master)
			VALUES (:id, :title, :company_id, :seat_master)`
	_, err := DB.NamedExec(query, seatMaster)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// GetSeatsByID is ...
func GetSeatsByID(id int) SeatMaster {
	seatMaster := SeatMaster{}
	err := DB.Get(&seatMaster, "SELECT id, title, seat_master, company_id, status, created, modified FROM seat_master WHERE id=?", id)
	if err != nil && err != sql.ErrNoRows {
		Logger.Error(err.Error())
	}
	return seatMaster
}

// EditSeats is ...
func EditSeats(data SeatMaster) {
	query := "UPDATE seat_master SET title=:title, company_id=:company_id, seat_master=:seat_master, status=:status WHERE id=:id"
	_, err := DB.NamedExec(query, data)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// DeleteSeatsByID function
func DeleteSeatsByID(id int) {
	_, err := DB.Exec("UPDATE seat_master SET status = 0 WHERE id = ?", id)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}
