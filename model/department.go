package model

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/tribalmedia/vista/setting"
)

// Department is ...
type Department struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CompanyID int       `db:"company_id"`
	Color     string    `db:"color"`
	Status    int       `db:"status"`
	Created   time.Time `db:"created"`
	Modified  time.Time `db:"modified"`
}

// GetAllDepartment is function that get all the department in database
func GetAllDepartment() []Department {
	department := Department{}
	departments := []Department{}
	rows, _ := DB.Queryx("SELECT id, name, color, company_id, status FROM departments")
	for rows.Next() {
		err := rows.StructScan(&department)
		if err != nil {
			Logger.Error(err.Error())
		}
		departments = append(departments, department)
	}

	return departments
}

// GetDepartmentByID is function get department by department_id
func GetDepartmentByID(departmentID int) Department {
	department := Department{}
	err := DB.Get(&department, "SELECT id, name, color, company_id, status FROM departments WHERE id = ?", departmentID)
	if err != nil && err != sql.ErrNoRows {
		Logger.Error(err.Error())
	}

	return department
}

// GetDepartmentByCompanyID is function return all departments in any company
func GetDepartmentByCompanyID(companyID int) []Department {
	department := Department{}
	departments := []Department{}
	rows, _ := DB.Queryx("SELECT id, name, color, company_id, status FROM departments WHERE company_id = ?", companyID)
	for rows.Next() {
		err := rows.StructScan(&department)
		if err != nil {
			Logger.Error(err.Error())
		}
		departments = append(departments, department)
	}

	return departments
}

// EditDepartment is function that update information of department
func EditDepartment(department Department) {
	query := "UPDATE departments SET name=:name, color=:color WHERE id=:id"
	_, err := DB.NamedExec(query, department)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// EditDepartmentsInCompany is function that edit all departments in company
func EditDepartmentsInCompany(companyID int, newDepartments []Department) {
	var oldDepartments []Department
	rows, _ := DB.Queryx("SELECT id, name, color, company_id, status FROM departments WHERE company_id = ?", companyID)
	for rows.Next() {
		// check every department in database
		var department Department
		err := rows.StructScan(&department)
		if err != nil {
			Logger.Error(err.Error())
		}

		// if department is not deleted, update it, else delete it
		isDelete := true
		for _, newValue := range newDepartments {
			if department.ID == newValue.ID {
				isDelete = false
				EditDepartment(newValue)
				oldDepartments = append(oldDepartments, newValue)
			}
		}
		if isDelete {
			DeleteDepartment(department)
		}
	}

	// if department input is not exist in database, create it
	for _, newValue := range newDepartments {
		isNew := true
		for _, oldValue := range oldDepartments {
			if newValue.ID == oldValue.ID {
				isNew = false
			}
		}

		if isNew {
			SaveDepartment(newValue)
		}
	}
}

// ValidateDepartment is function that validate department data before save
func ValidateDepartment(department Department) map[string]string {
	err := map[string]string{}

	if department.Name == "" {
		err["name_empty"] = "Department name can't be empty."
	}

	if len(department.Name) > setting.NameMaxLength {
		err["name_length"] = "Department name can't be more than " + strconv.Itoa(setting.NameMaxLength) + " characters."
	}

	if department.Color == "" {
		err["color_empty"] = "Department color can't be empty."
	}

	if len(department.Color) > setting.ColorCodeMaxLength {
		err["color_length"] = "Department color can't be more than " + strconv.Itoa(setting.ColorCodeMaxLength) + " characters."
	}

	return err
}

// ValidateListDepartment is function check if list department has empty value or duplicate value
func ValidateListDepartment(departments []Department) map[string]string {
	err := map[string]string{}

	var departmentNames []string
	var departmentColors []string

	// company has at least 1 department
	if len(departments) < setting.LeastDepartment {
		err["limited"] = "Company must has at least " + strconv.Itoa(setting.LeastDepartment) + " department"
		return err
	}

	// limit to 20 departments
	if len(departments) > setting.MaxDepartment {
		err["limited"] = "Limit to " + strconv.Itoa(setting.MaxDepartment) + " departments"
		return err
	}

	// check each department in list departments
	for _, department := range departments {
		errDepartment := ValidateDepartment(department)
		if len(errDepartment) != 0 {
			for errKey, errDetail := range errDepartment {
				err[errKey] = errDetail
			}
		} else {
			departmentNames = append(departmentNames, department.Name)
			departmentColors = append(departmentColors, department.Color)
		}
	}

	checkList := make(map[string]bool)
	for _, name := range departmentNames {
		if _, ok := checkList[name]; ok {
			// duplicate department name
			err["duplicate_name"] = "Duplicate department name"
		}
		checkList[name] = true
	}

	for _, color := range departmentColors {
		if _, ok := checkList[color]; ok {
			// duplicate department color
			err["duplicate_color"] = "Duplicate department color"
		}
		checkList[color] = true
	}

	return err
}

// SaveDepartment is function that insert new department
func SaveDepartment(department Department) {
	query := `INSERT INTO departments (name, company_id, color, status)
			VALUES (:name, :company_id, :color, :status)`
	_, err := DB.NamedExec(query, department)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// DeleteDepartment is ...
func DeleteDepartment(department Department) {
	query := "DELETE FROM departments WHERE id=:id"
	_, err := DB.NamedExec(query, department)
	if err != nil {
		Logger.Error(err.Error())
	}
}

// GetDepartmentsAndMemsByCompanyID get all departments and members of department by companyID
func GetDepartmentsAndMemsByCompanyID(companyID int) []map[string]interface{} {
	var DepartmentsList []map[string]interface{}
	var department Department
	rows, _ := DB.Queryx("SELECT id, name, color, company_id, status FROM departments WHERE company_id = ?", companyID)
	for rows.Next() {
		err := rows.StructScan(&department)
		if err != nil {
			Logger.Error(err.Error())
		}
		members := GetMembersOfDepartment(department.ID)
		m := map[string]interface{}{
			"Members":     members,
			"Departments": department,
		}
		DepartmentsList = append(DepartmentsList, m)
	}

	return DepartmentsList
}
