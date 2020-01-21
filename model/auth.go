package model

import (
	"database/sql"
)

// Auth is ...
type Auth struct {
	ID         int            `db:"id"`
	LoginID    string         `db:"login_id"`
	Password   string         `db:"password"`
	Name       string         `db:"name"`
	RoleType   int            `db:"role_type"`
	PictureURL sql.NullString `db:"picture_url"`
	CompanyID  int            `db:"company_id"`
}

// Authenticated user login
func Authenticated(loginID string, password string) (*Auth, error) {
	checkMember := Auth{}

	// check member's info
	err := DB.Get(&checkMember, "SELECT id, login_id, password, name, role_type, picture_url, company_id FROM members WHERE login_id = ?", loginID)
	if err != nil && err != sql.ErrNoRows {
		Logger.Error(err.Error())
	}

	// check password
	err = CheckPasswordHash(password, checkMember.Password)
	if err != nil {
		Logger.Error(err.Error())
		return nil, err
	}
	authInfo := &checkMember
	return authInfo, nil
}
