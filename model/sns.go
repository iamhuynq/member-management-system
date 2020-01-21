package model

import (
	"database/sql"
	"strconv"

	"github.com/tribalmedia/vista/setting"
)

// SNS is ...
type SNS struct {
	ID       int    `db:"id"`
	MemberID int    `db:"member_id"`
	Github   string `db:"github"`
}

// GetSNSAccountByMemberID is ...
func GetSNSAccountByMemberID(memberID int) (SNS, error) {
	sns := SNS{}
	err := DB.Get(&sns, "SELECT id, member_id, github FROM sns WHERE member_id = ?", memberID)
	if err != nil && err != sql.ErrNoRows {
		Logger.Error(err.Error())
	}

	return sns, err
}

// EditSNSAccount is ...
func EditSNSAccount(data SNS) {
	query := "UPDATE sns SET github=:github WHERE id=:id"
	_, err := DB.NamedExec(query, data)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// SaveSNSAccount is ...
func SaveSNSAccount(sns SNS) {
	query := `INSERT INTO sns (member_id, github)
			VALUES (:member_id, :github)`
	_, err := DB.NamedExec(query, sns)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// ValidateSNSAccount is function that validate SNS account data before save
func ValidateSNSAccount(account string) map[string]string {
	err := map[string]string{}

	if len(account) > setting.SNSAccountMaxLength {
		err["length"] = "Your SNS account can't be more than " + strconv.Itoa(setting.SNSAccountMaxLength) + " characters."
	}

	return err
}
