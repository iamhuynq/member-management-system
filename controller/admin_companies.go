package controller

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tribalmedia/vista/model"
	"github.com/tribalmedia/vista/service"
	"github.com/tribalmedia/vista/setting"
)

// CompanyAddDisplay is ...
func CompanyAddDisplay(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	templateData := map[string]interface{}{
		"title": "Add Company",
		"auth":  auth,
		"tab":   setting.CompaniesTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_companies/company_add.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// CompanyAdd is ...
func CompanyAdd(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	if err := r.ParseForm(); err != nil {
		Logger.Error(err.Error())
	}

	company := model.Company{}
	department := model.Department{}
	departmentList := []model.Department{}

	// Get company data
	company.Name = r.FormValue("companyName")
	company.Color = r.FormValue("companyColor")
	// Get departments data
	for index, value := range r.PostForm["departmentName"] {
		department.Name = value
		department.Color = r.PostForm["departmentColor"][index]
		// Add department to departmentList
		if department.Name != "" || department.Color != "" {
			departmentList = append(departmentList, department)
		}
	}

	// validate company
	companyError := model.ValidateCompany(company)

	// validate list department
	departmentListError := model.ValidateListDepartment(departmentList)

	if len(companyError) == 0 && len(departmentListError) == 0 {
		saveSuccess, companyID := model.SaveCompany(company)
		// Save each department
		if saveSuccess {
			for _, newDepartment := range departmentList {
				newDepartment.CompanyID = companyID
				model.SaveDepartment(newDepartment)
			}
			http.Redirect(w, r, "/admin/companies", 301)
		}
	}

	templateData := map[string]interface{}{
		"title":               "Add Company",
		"company":             company,
		"departmentList":      departmentList,
		"companyError":        companyError,
		"departmentListError": departmentListError,
		"auth":                auth,
		"tab":                 setting.CompaniesTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_companies/company_add.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// CompanyList is ...
func CompanyList(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	templateData := map[string]interface{}{
		"title":       "List Company",
		"companyList": model.GetAllCompany(),
		"auth":        auth,
		"tab":         setting.CompaniesTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_companies/company_list.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// CompanyDetail is ...
func CompanyDetail(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	vars := mux.Vars(r)
	id := vars["id"]
	companyID, _ := strconv.Atoi(id)

	templateData := map[string]interface{}{
		"company":        model.GetCompanyByID(companyID),
		"departmentList": model.GetDepartmentByCompanyID(companyID),
		"title":          "Company Detail",
		"auth":           auth,
		"tab":            setting.CompaniesTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_companies/company_detail.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// CompanyEditDisplay is ...
func CompanyEditDisplay(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	vars := mux.Vars(r)
	id := vars["id"]
	companyID, _ := strconv.Atoi(id)

	templateData := map[string]interface{}{
		"company":        model.GetCompanyByID(companyID),
		"departmentList": model.GetDepartmentByCompanyID(companyID),
		"title":          "Edit Company",
		"auth":           auth,
		"tab":            setting.CompaniesTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_companies/company_edit.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// CompanyEdit is ...
func CompanyEdit(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	if err := r.ParseForm(); err != nil {
		Logger.Error(err.Error())
	}

	vars := mux.Vars(r)
	id := vars["id"]
	companyID, _ := strconv.Atoi(id)

	// get information for edit
	company := model.GetCompanyByID(companyID)
	departmentList := model.GetDepartmentByCompanyID(companyID)

	department := model.Department{}
	newDepartmentList := []model.Department{}

	company.Name = r.FormValue("companyName")
	company.Color = r.FormValue("companyColor")
	for index, value := range r.PostForm["departmentID"] {
		department.ID, _ = strconv.Atoi(value)
		department.CompanyID = companyID
		department.Name = r.PostForm["departmentName"][index]
		department.Color = r.PostForm["departmentColor"][index]
		// add new department to departmentList
		if department.Name != "" || department.Color != "" {
			newDepartmentList = append(newDepartmentList, department)
		}
	}

	// validate company
	companyError := model.ValidateCompany(company)

	// validate list department
	departmentListError := model.ValidateListDepartment(newDepartmentList)

	if len(companyError) == 0 && len(departmentListError) == 0 {
		model.EditCompany(company)
		model.EditDepartmentsInCompany(companyID, newDepartmentList)
		http.Redirect(w, r, "/admin/companies/"+id, 301)
	}

	templateData := map[string]interface{}{
		"title":               "Edit Company",
		"company":             company,
		"departmentList":      departmentList,
		"newDepartmentList":   newDepartmentList,
		"companyError":        companyError,
		"departmentListError": departmentListError,
		"auth":                auth,
		"tab":                 setting.CompaniesTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_companies/company_edit.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}
