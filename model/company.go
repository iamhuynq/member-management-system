package model

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/tribalmedia/vista/setting"
)

//Company is ...
type Company struct {
	ID       int       `db:"id"`
	Name     string    `db:"name"`
	Color    string    `db:"color"`
	Created  time.Time `db:"created"`
	Modified time.Time `db:"modified"`
}

// GetAllCompany is function that get all the company in database
func GetAllCompany() []Company {
	company := Company{}
	companies := []Company{}
	rows, _ := DB.Queryx("SELECT id, name, color, created, modified FROM companies ORDER BY id DESC")
	for rows.Next() {
		err := rows.StructScan(&company)
		if err != nil {
			Logger.Error(err.Error())
		}
		companies = append(companies, company)
	}

	return companies
}

// GetCompanyByID is function get information of company by company_id
func GetCompanyByID(id int) Company {
	company := Company{}
	err := DB.Get(&company, "SELECT id, name, color, created, modified FROM companies WHERE id = ?", id)
	if err != nil && err != sql.ErrNoRows {
		Logger.Error(err.Error())
	}

	return company
}

// SaveCompany is function that add new company and return company_id for add new departments
func SaveCompany(company Company) (bool, int) {
	// create new company in database
	query := `INSERT INTO companies (name, color)
			VALUES (:name, :color)`
	_, err := DB.NamedExec(query, company)
	if err != nil {
		Logger.Fatal(err.Error())
	}

	// return company_id already added
	var companyID int
	err = DB.Get(&companyID, "SELECT id FROM companies ORDER BY id DESC LIMIT 1")
	if err != nil && err != sql.ErrNoRows {
		Logger.Error(err.Error())
		return false, 0
	}

	return true, companyID
}

// EditCompany is function update information of company
func EditCompany(company Company) {
	query := "UPDATE companies SET name=:name, color=:color WHERE id=:id"
	_, err := DB.NamedExec(query, company)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// CheckCompanyColorExist is function that check if company color exist
func CheckCompanyColorExist(color string, companyID int) bool {
	var resultCompanyID int

	err := DB.Get(&resultCompanyID, "SELECT id FROM companies WHERE color = ? AND id != ?", color, companyID)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		Logger.Error(err.Error())
	}

	return true
}

// ValidateCompany is function that validate company data before save
func ValidateCompany(company Company) map[string]string {
	err := map[string]string{}

	if company.Name == "" {
		err["name"] = "Company name can't be empty."
	}

	if len(company.Name) > setting.NameMaxLength {
		err["name"] = "Company name can't be more than " + strconv.Itoa(setting.NameMaxLength) + " characters."
	}

	if company.Color == "" {
		err["color"] = "Company color can't be empty."
	}

	if len(company.Color) > setting.ColorCodeMaxLength {
		err["color"] = "Company color can't be more than " + strconv.Itoa(setting.ColorCodeMaxLength) + " characters."
	}

	if CheckCompanyColorExist(company.Color, company.ID) {
		err["color"] = "Company color is exist."
	}

	return err
}
