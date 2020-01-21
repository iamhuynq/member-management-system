package model

import (
	"database/sql"
	"time"
)

// Seat is ...
type Seat struct {
	ID           int       `db:"id"`
	MemberID     int       `db:"member_id"`
	TableNumber  int       `db:"table_number"`
	SeatMasterID int       `db:"seat_master_id"`
	Row          int       `db:"row"`
	Col          int       `db:"col"`
	Created      time.Time `db:"created"`
	Modified     time.Time `db:"modified"`
}

// AddSeat add data to seat table
func AddSeat(seat Seat) {
	query := `INSERT INTO seats (id, member_id, table_number, seat_master_id, row, col)
		VALUES (:id, :member_id, :table_number, :seat_master_id, :row, :col)`
	_, err := DB.NamedExec(query, seat)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// GetAllSeat get all seats from seats by seatMasterID
func GetAllSeat(seatMasterID int) []Seat {
	seat := Seat{}
	seats := []Seat{}
	rows, _ := DB.Queryx("SELECT id, member_id, table_number, seat_master_id, row, col FROM seats WHERE seat_master_id = ?", seatMasterID)
	for rows.Next() {
		err := rows.StructScan(&seat)
		if err != nil {
			Logger.Error(err.Error())
		}
		seats = append(seats, seat)
	}

	return seats
}

// GetSeat get seat from seats by seatMasterID and memberID
func GetSeat(seatMasterID, memberID int) Seat {
	seat := Seat{}
	err := DB.Get(&seat, "SELECT id, member_id, table_number, seat_master_id, row, col FROM seats WHERE seat_master_id = ? AND member_id = ?", seatMasterID, memberID)
	if err != nil && err != sql.ErrNoRows {
		Logger.Error(err.Error())
	}

	return seat
}

// EditSeat edits seat
func EditSeat(seat Seat) {
	query := "UPDATE seats SET member_id=:member_id, table_number=:table_number, seat_master_id=:seat_master_id, row=:row, col=:col WHERE id=:id"
	_, err := DB.NamedExec(query, seat)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// DeleteSeat deletes seat
func DeleteSeat(id int) {
	_, err := DB.Exec("DELETE FROM seats WHERE id = ?", id)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}
