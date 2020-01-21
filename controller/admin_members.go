package controller

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/tribalmedia/vista/model"
	"github.com/tribalmedia/vista/service"
	"github.com/tribalmedia/vista/setting"
	"go.uber.org/zap"
)

// Logger is ...
var Logger *zap.Logger

// MembersList is ...
func MembersList(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)

	templateData := map[string]interface{}{
		"title": "Member List",
		"list":  model.GetMember(),
		"auth":  auth,
		"tab":   setting.MembersTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_members/member_list.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// MembersAdd is function hanlde post request
func MembersAdd(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	// Parse form
	if err := r.ParseMultipartForm(setting.FileMaxSize); err != nil {
		Logger.Error(err.Error())
	}

	// Default avatar
	avatar := setting.SetDefaultAvatar().MemberAvatar

	member := model.Member{}
	var convertErr error
	var teamList []int

	member.LoginID = r.FormValue("loginID")
	member.Password = r.FormValue("password")
	member.Name = r.FormValue("name")
	member.RoleType, _ = strconv.Atoi(r.FormValue("role"))
	member.GenderType, convertErr = strconv.Atoi(r.FormValue("gender"))
	member.Comment = sql.NullString{String: r.FormValue("comment"), Valid: true}
	member.CompanyID, convertErr = strconv.Atoi(r.FormValue("company"))
	member.DepartmentID, convertErr = strconv.Atoi(r.FormValue("department"))
	githubAccount := r.FormValue("githubAccount")

	// timestamp is time if user not input date
	timestamp, _ := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
	birthdayInput, _ := time.Parse("01-02-2006", r.FormValue("birthday"))
	if birthdayInput == timestamp {
		member.Birthday = mysql.NullTime{Time: birthdayInput, Valid: false}
	} else {
		member.Birthday = mysql.NullTime{Time: birthdayInput, Valid: true}
	}

	for _, value := range r.PostForm["team"] {
		teamID, _ := strconv.Atoi(value)
		if teamID != 0 {
			teamList = append(teamList, teamID)
		}
	}

	if convertErr != nil {
		Logger.Error(convertErr.Error())
	}

	passwordData := map[string]string{
		"password": member.Password,
	}

	// validate member data
	validateErr := model.ValidateMember(member, r.ContentLength)

	passwordErr := model.ValidatePassword(passwordData)
	snsAccountErr := model.ValidateSNSAccount(githubAccount)

	// if validate are ok
	if len(validateErr) == 0 && len(passwordErr) == 0 && len(snsAccountErr) == 0 {

		// get data of file uploaded
		file, handler, err := r.FormFile("avatar")

		if err == nil {
			avatar = handler.Filename

			//create temporary file to save image
			tempFile, err := ioutil.TempFile(setting.ImageBaseURL, "*"+avatar)

			//get path of image
			avatar = tempFile.Name()
			if err != nil {
				Logger.Error(err.Error())
			}
			defer tempFile.Close()

			// read all of the contents of our uploaded file into a byte array
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				Logger.Error(err.Error())
			}

			// write this byte array to our temporary file
			tempFile.Write(fileBytes)

			// if app running on EC2
			if setting.UseS3Service() {
				// upload file to aws S3
				service.UploadImageToS3(avatar, setting.MemberFolderType)
				// save path on S3 bucket
				avatar = setting.S3BucketURL + setting.S3MemberFolder + filepath.Base(avatar)
			}
		} else if err != http.ErrMissingFile {
			Logger.Error(err.Error())
		}

		member.PictureURL = sql.NullString{String: avatar, Valid: true}

		memberID := model.SaveMember(member)
		sns := model.SNS{
			MemberID: memberID,
			Github:   githubAccount,
		}
		model.SaveSNSAccount(sns)

		// get array of team of user
		for _, teamID := range teamList {
			if !model.CheckTeamMemberExist(teamID, memberID) {
				model.SaveTeamMember(teamID, memberID)
			}
		}
		http.Redirect(w, r, "/admin/members", 301)
	}

	// get all team, company and department to select
	allTeams := model.GetAllTeam()
	companies := model.GetAllCompany()
	departments := model.GetAllDepartment()

	templateData := map[string]interface{}{
		"allTeams":        allTeams.List,
		"companies":       companies,
		"departments":     departments,
		"validateError":   validateErr,
		"passwordError":   passwordErr,
		"snsAccountError": snsAccountErr,
		"member":          member,
		"teamList":        teamList,
		"snsAccount":      githubAccount,
		"title":           "Add Member",
		"auth":            auth,
		"tab":             setting.MembersTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_members/member_add.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// MembersAddDisplay is function that display the add member page
func MembersAddDisplay(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	// get all team, company and department to select
	allTeams := model.GetAllTeam()
	companies := model.GetAllCompany()
	departments := model.GetAllDepartment()
	member := model.Member{}

	templateData := map[string]interface{}{
		"allTeams":      allTeams.List,
		"companies":     companies,
		"departments":   departments,
		"member":        member,
		"defaultGender": setting.MaleGenderType,
		"title":         "Add Member",
		"auth":          auth,
		"tab":           setting.MembersTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_members/member_add.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// MembersDetail is ...
func MembersDetail(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	vars := mux.Vars(r)
	id := vars["id"]

	memberID, _ := strconv.Atoi(id)

	member := model.GetMemberByID(memberID)
	teamList, _ := model.GetTeamByMember(member.ID)
	company := model.GetCompanyByID(member.CompanyID)
	department := model.GetDepartmentByID(member.DepartmentID)

	templateData := map[string]interface{}{
		"member":     member,
		"teamList":   teamList,
		"company":    company.Name,
		"department": department.Name,
		"title":      "Member Detail",
		"auth":       auth,
		"tab":        setting.MembersTab,
		"useS3":      setting.UseS3Service(),
	}
	sns, err := model.GetSNSAccountByMemberID(memberID)
	if err == nil {
		templateData["SNSAccount"] = sns
	}

	tmpl := template.Must(template.ParseFiles("template/admin_members/member_detail.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// MembersEdit is function hanlde post request
func MembersEdit(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	vars := mux.Vars(r)
	id := vars["id"]

	var teamList []int

	memberID, _ := strconv.Atoi(id)

	if err := r.ParseMultipartForm(setting.FileMaxSize); err != nil {
		Logger.Error(err.Error())
	}
	checkAvt := r.FormValue("check")
	member := model.GetMemberByID(memberID)
	member.ID = memberID
	member.Name = r.FormValue("name")
	member.LoginID = r.FormValue("loginID")
	member.RoleType, _ = strconv.Atoi(r.FormValue("role"))
	member.GenderType, _ = strconv.Atoi(r.FormValue("gender"))
	member.Status, _ = strconv.Atoi(r.FormValue("status"))
	member.Comment = sql.NullString{String: r.FormValue("comment"), Valid: true}
	member.CompanyID, _ = strconv.Atoi(r.FormValue("company"))
	member.DepartmentID, _ = strconv.Atoi(r.FormValue("department"))
	githubAccount := r.FormValue("githubAccount")

	// timestamp is time if user not input date
	timestamp, _ := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
	birthdayInput, _ := time.Parse("01-02-2006", r.FormValue("birthday"))
	if birthdayInput == timestamp {
		member.Birthday = mysql.NullTime{Time: birthdayInput, Valid: false}
	} else {
		member.Birthday = mysql.NullTime{Time: birthdayInput, Valid: true}
	}

	for _, value := range r.PostForm["team"] {
		teamID, _ := strconv.Atoi(value)
		if teamID != 0 {
			teamList = append(teamList, teamID)
		}
	}

	validateError := model.ValidateMember(member, r.ContentLength)
	snsAccountErr := model.ValidateSNSAccount(githubAccount)

	if len(validateError) == 0 && len(snsAccountErr) == 0 {
		file, handler, err := r.FormFile("avatar")
		defaultAvatar := setting.SetDefaultAvatar().MemberAvatar

		if err == nil {
			oldImage := filepath.Base(member.PictureURL.String)
			pictureURL := handler.Filename
			// Create a temporary file
			tempFile, err := ioutil.TempFile(setting.ImageBaseURL, "*"+pictureURL)
			//get path of image
			avatar := tempFile.Name()
			if err != nil {
				Logger.Error(err.Error())
			}
			defer tempFile.Close()
			// read all of the contents of our uploaded file into a byte array
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				Logger.Error(err.Error())
			}
			// write this byte array to our temporary file
			tempFile.Write(fileBytes)

			// if app running on EC2
			if setting.UseS3Service() {
				// delete old image
				if oldImage != filepath.Base(defaultAvatar) {
					service.DeleteImageFromS3(oldImage, setting.MemberFolderType)
				}

				// save new image to aws S3
				service.UploadImageToS3(avatar, setting.MemberFolderType)
				avatar = setting.S3BucketURL + setting.S3MemberFolder + filepath.Base(avatar)
			}

			member.PictureURL = sql.NullString{String: avatar, Valid: true}
		} else if err != http.ErrMissingFile {
			Logger.Error(err.Error())
		}

		// if user delete avatar
		if checkAvt == setting.DeleteAvatar && member.PictureURL.String != defaultAvatar {
			// delete image on S3
			if setting.UseS3Service() {
				fileName := filepath.Base(member.PictureURL.String)
				service.DeleteImageFromS3(fileName, setting.MemberFolderType)
			}
			// set default avatar
			member.PictureURL = sql.NullString{String: defaultAvatar, Valid: true}
		}

		model.EditMember(member)
		sns, _ := model.GetSNSAccountByMemberID(memberID)
		sns.Github = githubAccount
		model.EditSNSAccount(sns)

		// get array of team of user
		model.EditTeamOfMember(memberID, teamList)
		http.Redirect(w, r, "/admin/members", 301)
	}

	// get all team, company and department to select
	allTeams := model.GetAllTeam()
	companies := model.GetAllCompany()
	departments := model.GetAllDepartment()

	res := map[string]interface{}{
		"member":          member,
		"allTeams":        allTeams.List,
		"companies":       companies,
		"departments":     departments,
		"SNSAccount":      githubAccount,
		"validateError":   validateError,
		"snsAccountError": snsAccountErr,
		"title":           "Edit Member",
		"teamList":        teamList,
		"auth":            auth,
		"tab":             setting.MembersTab,
		"useS3":           setting.UseS3Service(),
	}

	tmpl := template.Must(template.ParseFiles("template/admin_members/member_edit.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", res); err != nil {
		Logger.Error(err.Error())
	}
}

//MembersEditDisplay is function that display the edit member page
func MembersEditDisplay(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	vars := mux.Vars(r)
	id := vars["id"]

	memberID, _ := strconv.Atoi(id)

	member := model.GetMemberByID(memberID)
	allTeams := model.GetAllTeam()
	_, teamList := model.GetTeamByMember(memberID)
	companies := model.GetAllCompany()
	departments := model.GetAllDepartment()

	templateData := map[string]interface{}{
		"member":      member,
		"allTeams":    allTeams.List,
		"teamList":    teamList,
		"companies":   companies,
		"departments": departments,
		"title":       "Edit Member",
		"auth":        auth,
		"tab":         setting.MembersTab,
		"useS3":       setting.UseS3Service(),
	}
	sns, err := model.GetSNSAccountByMemberID(memberID)
	if err == nil {
		templateData["SNSAccount"] = sns.Github
	}

	tmpl := template.Must(template.ParseFiles("template/admin_members/member_edit.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// MembersDelete is change status from In (0) to Out (1)
func MembersDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	memberID, _ := strconv.Atoi(id)

	member := model.GetMemberByID(memberID)
	member.Status = setting.MemberOutStatus

	model.EditMember(member)
	http.Redirect(w, r, "/admin/members", 301)
}

// MoveMember is ...
func MoveMember(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	memberID, _ := strconv.Atoi(r.FormValue("memberID"))
	teamID, _ := strconv.Atoi(r.FormValue("teamID"))
	teamIDPre, _ := strconv.Atoi(r.FormValue("teamIDPrevious"))
	_, teamOld := model.GetTeamByMember(memberID)
	if model.MoveMember(memberID, teamID, teamIDPre) {
		_, teamNew := model.GetTeamByMember(memberID)
		// teamDiff include new team_id of member from database
		teamDiff := service.Difference(teamNew, teamOld)

		if len(teamDiff) > 0 {
			w.Write([]byte(strconv.Itoa(teamDiff[0])))
		} else {
			// when drop one member to different teams at the same time, user drop after
			w.Write([]byte(strconv.Itoa(teamNew[0])))
		}
	} else {
		w.Write([]byte("false"))
	}
}

// ChangePassword is function that update user password
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	member := model.GetMemberByID(auth.ID)

	if err := r.ParseMultipartForm(setting.FileMaxSize); err != nil {
		Logger.Error(err.Error())
	}

	newPassword := r.FormValue("newPassword")
	confirmPassword := r.FormValue("confirmPassword")
	oldPassword := r.FormValue("oldPassword")

	passwordData := map[string]string{
		"password":        newPassword,
		"confirmPassword": confirmPassword,
	}

	validateError := model.ValidatePassword(passwordData)

	// check if password is correct
	_, err := model.Authenticated(member.LoginID, oldPassword)

	if err != nil {
		validateError["oldPassword"] = "Password is not correct."
	}

	if len(validateError) == 0 {
		model.UpdatePassword(newPassword, member)
		http.Redirect(w, r, "/my_page", 301)
	}

	templateData := map[string]interface{}{
		"title":         "Change Password",
		"validateError": validateError,
		"auth":          auth,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_members/change_password.tmpl", setting.UserTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// ChangePasswordDisplay is function dis play change password page
func ChangePasswordDisplay(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	templateData := map[string]interface{}{
		"title": "Change Password",
		"auth":  auth,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_members/change_password.tmpl", setting.UserTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// MyPageDisplay is ...
func MyPageDisplay(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)

	member := model.GetMemberByID(auth.ID)

	templateData := map[string]interface{}{
		"member": member,
		"auth":   auth,
		"title":  "My Page",
		"useS3":  setting.UseS3Service(),
	}

	sns, err := model.GetSNSAccountByMemberID(auth.ID)
	if err == nil {
		templateData["SNSAccount"] = sns
	}

	tmpl := template.Must(template.ParseFiles("template/admin_members/my_page.tmpl", setting.UserTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// MyPage is ...
func MyPage(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)

	if err := r.ParseMultipartForm(setting.FileMaxSize); err != nil {
		Logger.Error(err.Error())
	}

	checkAvt := r.FormValue("check")
	member := model.GetMemberByID(auth.ID)
	member.Name = r.FormValue("name")
	member.GenderType, _ = strconv.Atoi(r.FormValue("gender"))
	member.Comment = sql.NullString{String: r.FormValue("comment"), Valid: true}
	githubAccount := r.FormValue("githubAccount")

	// timestamp is time if user not input date
	timestamp, _ := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
	birthdayInput, _ := time.Parse("01-02-2006", r.FormValue("birthday"))
	if birthdayInput == timestamp {
		member.Birthday = mysql.NullTime{Time: birthdayInput, Valid: false}
	} else {
		member.Birthday = mysql.NullTime{Time: birthdayInput, Valid: true}
	}

	validateError := model.ValidateMember(member, r.ContentLength)

	if len(validateError) == 0 {
		file, handler, err := r.FormFile("myFile")
		defaultAvatar := setting.SetDefaultAvatar().MemberAvatar

		if err == nil {
			oldImage := filepath.Base(member.PictureURL.String)
			pictureURL := handler.Filename
			// Create a temporary file
			tempFile, err := ioutil.TempFile(setting.ImageBaseURL, "*"+pictureURL)
			//get path of image
			avatar := tempFile.Name()
			if err != nil {
				Logger.Error(err.Error())
			}
			defer tempFile.Close()
			// read all of the contents of our uploaded file into a byte array
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				Logger.Error(err.Error())
			}
			// write this byte array to our temporary file
			tempFile.Write(fileBytes)

			if setting.UseS3Service() {
				// delete old image
				if oldImage != filepath.Base(defaultAvatar) {
					service.DeleteImageFromS3(oldImage, setting.MemberFolderType)
				}

				// save new image to aws S3
				service.UploadImageToS3(avatar, setting.MemberFolderType)
				avatar = setting.S3BucketURL + setting.S3MemberFolder + filepath.Base(avatar)
			}

			member.PictureURL = sql.NullString{String: avatar, Valid: true}
		} else if err != http.ErrMissingFile {
			Logger.Error(err.Error())
		}

		// delete avatar on S3
		if checkAvt == setting.DeleteAvatar && member.PictureURL.String != defaultAvatar {
			// delete image on S3
			if setting.UseS3Service() {
				fileName := filepath.Base(member.PictureURL.String)
				service.DeleteImageFromS3(fileName, setting.MemberFolderType)
			}
			// set default avatar
			member.PictureURL = sql.NullString{String: defaultAvatar, Valid: true}
		}

		model.EditMember(member)
		sns, _ := model.GetSNSAccountByMemberID(auth.ID)
		sns.Github = githubAccount
		model.EditSNSAccount(sns)

		http.Redirect(w, r, "/team", 301)
	}

	templateData := map[string]interface{}{
		"member":        member,
		"validateError": validateError,
		"newSNSAccount": githubAccount,
		"title":         "My Page",
		"auth":          auth,
		"useS3":         setting.UseS3Service(),
	}

	tmpl := template.Must(template.ParseFiles("template/admin_members/my_page.tmpl", setting.UserTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// ResetPassword is function allow admin reset an user's password
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	vars := mux.Vars(r)
	id := vars["id"]

	memberID, _ := strconv.Atoi(id)
	member := model.GetMemberByID(memberID)

	if err := r.ParseMultipartForm(setting.FileMaxSize); err != nil {
		Logger.Error(err.Error())
	}

	newPassword := r.FormValue("newPassword")
	confirmPassword := r.FormValue("confirmPassword")

	passwordData := map[string]string{
		"password":        newPassword,
		"confirmPassword": confirmPassword,
	}

	validateError := model.ValidatePassword(passwordData)

	if len(validateError) == 0 {
		model.UpdatePassword(newPassword, member)
		http.Redirect(w, r, "/admin/members/"+id, 301)
	}

	templateData := map[string]interface{}{
		"title":         "Reset Password",
		"validateError": validateError,
		"id":            id,
		"auth":          auth,
		"tab":           setting.MembersTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_members/reset_password.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// ResetPasswordDisplay is ...
func ResetPasswordDisplay(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	vars := mux.Vars(r)
	id := vars["id"]

	templateData := map[string]interface{}{
		"title": "Reset Password",
		"id":    id,
		"auth":  auth,
		"tab":   setting.MembersTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_members/reset_password.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}
