package model

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/tribalmedia/vista/setting"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Logger is ...
var Logger *zap.Logger

// Member is ...
type Member struct {
	ID           int            `db:"id"`
	LoginID      string         `db:"login_id"`
	Password     string         `db:"password"`
	Name         string         `db:"name"`
	RoleType     int            `db:"role_type"`
	Birthday     mysql.NullTime `db:"birthday"`
	PictureURL   sql.NullString `db:"picture_url"`
	GenderType   int            `db:"gender_type"`
	Comment      sql.NullString `db:"comment"`
	Status       int            `db:"status"`
	DepartmentID int            `db:"department_id"`
	CompanyID    int            `db:"company_id"`
	Created      time.Time      `db:"created"`
	Modified     time.Time      `db:"modified"`
}

// HashPassword is function that encrypt a password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash is function that check password match with encrypted password
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// CheckMemberExist is function that check if user exist with login_id
func CheckMemberExist(memberID int, loginID string) bool {
	var id int

	err := DB.Get(&id, "SELECT members.id FROM members WHERE members.login_id = ? AND members.id != ?", loginID, memberID)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		Logger.Error(err.Error())
	}

	return true
}

// SaveMember is function that save member data to database
// It return true if save success
func SaveMember(member Member) int {
	member.Password, _ = HashPassword(member.Password)
	query := `INSERT INTO members (login_id,password,name,role_type,birthday,picture_url,gender_type,comment,status,department_id,company_id)
			VALUES (:login_id, :password, :name, :role_type, :birthday, :picture_url, :gender_type, :comment, :status, :department_id, :company_id)`
	res, err := DB.NamedExec(query, member)
	if err != nil {
		Logger.Fatal(err.Error())
	}

	memberID, _ := res.LastInsertId()
	return int(memberID)
}

// UpdatePassword is ...
func UpdatePassword(newPassword string, member Member) {
	member.Password, _ = HashPassword(newPassword)

	query := "UPDATE members SET password=:password  WHERE id=:id"

	_, err := DB.NamedExec(query, member)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// ValidatePassword is for password only. Seperate with other field
func ValidatePassword(passwordData map[string]string) map[string]string {
	err := map[string]string{}

	if len(passwordData["password"]) < setting.PasswordMinLength {
		err["password"] = "Password must be at least " + strconv.Itoa(setting.PasswordMinLength) + " characters."
	}

	if passwordData["password"] == "" {
		err["password"] = "Password can not be empty."
	}

	// check if confirm_password exist
	if _, key := passwordData["confirmPassword"]; key {
		if passwordData["confirmPassword"] != passwordData["password"] {
			err["confirmPassword"] = "Confirm password must match with new password."
		}
	}

	return err
}

// ValidateMember is function that validate member data before save
func ValidateMember(member Member, contentLength int64) map[string]string {
	err := map[string]string{}
	// error message when violate validation
	if len(member.Name) > setting.NameMaxLength {
		err["name"] = "Your name can't be more than " + strconv.Itoa(setting.NameMaxLength) + " characters."
	}

	if member.Name == "" {
		err["name"] = "Your name can't be empty."
	}

	if len(member.Comment.String) > setting.CommentMaxLength {
		err["comment"] = "Your comment can't be more than " + strconv.Itoa(setting.CommentMaxLength) + " characters."
	}

	if member.CompanyID == 0 {
		err["company"] = "Please select company."
	}

	if member.LoginID == "" {
		err["loginID"] = "Login ID can not be empty."
	}

	if len(member.LoginID) > setting.LoginIDMaxLength {
		err["loginID"] = "Login ID can't be more than " + strconv.Itoa(setting.LoginIDMaxLength) + " characters."
	}

	if contentLength > setting.FileMaxSize {
		err["photo"] = "File size must be less than 3 MB"
	}

	if CheckMemberExist(member.ID, member.LoginID) {
		err["loginID"] = "Account already exists."
	}

	return err
}

// GetMemberByID gets member by member's id
// It returns the member
func GetMemberByID(id int) Member {
	member := Member{}
	err := DB.Get(&member, "SELECT id,login_id,password,name,role_type,birthday,picture_url,gender_type,comment,status,department_id,company_id FROM members WHERE id=?", id)
	if err != nil && err != sql.ErrNoRows {
		Logger.Error(err.Error())
	}

	return member
}

// EditMember updates info member
// It returns true if save success
func EditMember(member Member) {
	query := "UPDATE members SET name=:name, role_type=:role_type, login_id=:login_id, birthday=:birthday, picture_url=:picture_url, gender_type=:gender_type, comment=:comment, status=:status, department_id=:department_id, company_id=:company_id WHERE id=:id"

	_, err := DB.NamedExec(query, member)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// GetMember is function get a list of all member
func GetMember() []map[string]interface{} {
	rows, _ := DB.Queryx("SELECT members.id,members.name,members.role_type,members.login_id,members.company_id,members.status,members.created,members.modified FROM members ORDER BY members.id DESC")
	var listMember []map[string]interface{}
	for rows.Next() {
		var member Member
		err := rows.StructScan(&member)
		if err != nil {
			Logger.Error(err.Error())
		}

		teamList, _ := GetTeamByMember(member.ID)

		m := map[string]interface{}{
			"Member":  member,
			"Teams":   teamList,
			"Company": GetCompanyByID(member.CompanyID).Name,
		}
		listMember = append(listMember, m)
	}
	return listMember
}

//GetMemberOfTeam is ...
func GetMemberOfTeam(id int) []map[string]interface{} {
	var listMember []map[string]interface{}
	rows, _ := DB.Queryx("SELECT team_member.member_id FROM team_member WHERE team_member.team_id=?", id)
	for rows.Next() {
		var memberID int
		err := rows.Scan(&memberID)
		if err != nil {
			Logger.Error(err.Error())
		}
		teamList, _ := GetTeamByMember(memberID)
		member := GetMemberByID(memberID)
		m := map[string]interface{}{
			"Member":   member,
			"Team":     teamList,
			"Company":  GetCompanyByID(member.CompanyID),
			"IsLeader": CheckLeader(id, memberID),
		}
		sns, err := GetSNSAccountByMemberID(member.ID)
		if err == nil {
			m["SNSAccount"] = sns.Github
		}
		if member.Status == setting.MemberInStatus {
			listMember = append(listMember, m)
		}
	}
	return listMember
}

// GetMembersOfDepartment gets members of department
// Id is id of department
func GetMembersOfDepartment(id int) []Member {
	member := Member{}
	members := []Member{}
	rows, _ := DB.Queryx("SELECT id, name FROM members WHERE department_id = ?", id)
	for rows.Next() {
		err := rows.StructScan(&member)
		if err != nil {
			Logger.Error(err.Error())
		}
		members = append(members, member)
	}

	return members
}
