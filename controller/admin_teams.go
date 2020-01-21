package controller

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tribalmedia/vista/model"
	"github.com/tribalmedia/vista/service"
	"github.com/tribalmedia/vista/setting"
)

// TeamsList is ...
func TeamsList(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	templateData := map[string]interface{}{
		"title": "Team List",
		"list":  model.GetTeams(),
		"auth":  auth,
		"tab":   setting.TeamsTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_teams/team_list.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// TeamsAdd is function that add a new team
func TeamsAdd(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	// Parse form
	if err := r.ParseMultipartForm(setting.FileMaxSize); err != nil {
		Logger.Error(err.Error())
	}

	// Default team icon
	avatar := setting.SetDefaultAvatar().TeamAvatar

	team := model.Team{}
	team.Name = r.FormValue("team")
	team.Description = sql.NullString{String: r.FormValue("description"), Valid: true}
	leaderID, _ := strconv.Atoi(r.FormValue("leader"))

	// validate team data
	validateError := model.ValidateTeam(team, r.ContentLength)

	if len(validateError) == 0 {

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
				service.UploadImageToS3(avatar, setting.TeamFolderType)
				// save path on S3 bucket
				avatar = setting.S3BucketURL + setting.S3TeamFolder + filepath.Base(avatar)
			}
		} else if err != http.ErrMissingFile {
			Logger.Error(err.Error())
		}

		team.PictureURL = sql.NullString{String: avatar, Valid: true}
		// save data to database
		teamID := model.SaveTeam(team)
		if leaderID != 0 {
			model.AddLeaderNewTeam(leaderID, int(teamID))
		}
		http.Redirect(w, r, "/admin/teams", 301)
	}

	templateData := map[string]interface{}{
		"validateError": validateError,
		"team":          team,
		"title":         "Add Team",
		"allMember":     model.GetMember(),
		"leaderID":      leaderID,
		"method":        "post",
		"auth":          auth,
		"tab":           setting.TeamsTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_teams/team_add.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// TeamsAddDisplay is function that display add Team page
func TeamsAddDisplay(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	templateData := map[string]interface{}{
		"title":     "Add Team",
		"allMember": model.GetMember(),
		"method":    "get",
		"auth":      auth,
		"tab":       setting.TeamsTab,
	}

	tmpl := template.Must(template.ParseFiles("template/admin_teams/team_add.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// TeamsDetail is ...
func TeamsDetail(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	var leader map[string]interface{}
	vars := mux.Vars(r)
	id := vars["id"]
	teamID, _ := strconv.Atoi(id)
	team := model.GetTeamByID(teamID)
	if model.CheckTeamHasLeader(teamID) {
		leader = model.LeaderOfTeam(teamID)
	}
	templateData := map[string]interface{}{
		"team":   team,
		"title":  "Team Detail",
		"leader": leader,
		"auth":   auth,
		"tab":    setting.TeamsTab,
		"useS3":  setting.UseS3Service(),
	}

	tmpl := template.Must(template.ParseFiles("template/admin_teams/team_detail.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// TeamsEdit is function hanlde post request
func TeamsEdit(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	vars := mux.Vars(r)
	id := vars["id"]
	if err := r.ParseMultipartForm(setting.FileMaxSize); err != nil {
		Logger.Error(err.Error())
	}

	teamID, _ := strconv.Atoi(id)
	checkAvt := r.FormValue("check")
	team := model.GetTeamByID(teamID)
	team.ID = teamID
	team.Name = r.FormValue("name")
	team.Description = sql.NullString{String: r.FormValue("description"), Valid: true}
	leaderID, _ := strconv.Atoi(r.FormValue("leader"))
	validateError := model.ValidateTeam(team, r.ContentLength)

	if len(validateError) == 0 {
		file, handler, err := r.FormFile("myFile")
		defaultAvatar := setting.SetDefaultAvatar().TeamAvatar

		if err == nil {
			oldImage := filepath.Base(team.PictureURL.String)
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
					service.DeleteImageFromS3(oldImage, setting.TeamFolderType)
				}

				// save new image to aws S3
				service.UploadImageToS3(avatar, setting.TeamFolderType)
				avatar = setting.S3BucketURL + setting.S3TeamFolder + filepath.Base(avatar)
			}

			team.PictureURL = sql.NullString{String: avatar, Valid: true}
		} else if err != http.ErrMissingFile {
			Logger.Error(err.Error())
		}

		// if user delete team's avatar
		if checkAvt == setting.DeleteAvatar && team.PictureURL.String != defaultAvatar {
			// delete image on S3
			if setting.UseS3Service() {
				fileName := filepath.Base(team.PictureURL.String)
				service.DeleteImageFromS3(fileName, setting.TeamFolderType)
			}
			// set default avatar
			team.PictureURL = sql.NullString{String: defaultAvatar, Valid: true}
		}

		model.EditTeam(team)
		model.EditTeamLeader(teamID, leaderID)
		http.Redirect(w, r, "/admin/teams", 301)
	}
	leader := map[string]interface{}{
		"leaderID": leaderID,
	}

	templateData := map[string]interface{}{
		"team":          team,
		"validateError": validateError,
		"title":         "Edit Team",
		"memberList":    model.GetMemberOfTeam(teamID),
		"leader":        leader,
		"auth":          auth,
		"tab":           setting.TeamsTab,
		"useS3":         setting.UseS3Service(),
	}

	tmpl := template.Must(template.ParseFiles("template/admin_teams/team_edit.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

// TeamsEditDisplay is function that display the edit team page
func TeamsEditDisplay(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	var leader map[string]interface{}
	vars := mux.Vars(r)
	id := vars["id"]

	teamID, _ := strconv.Atoi(id)
	if model.CheckTeamHasLeader(teamID) {
		leader = model.LeaderOfTeam(teamID)
	}
	team := model.GetTeamByID(teamID)
	templateData := map[string]interface{}{
		"team":       team,
		"title":      "Edit Team",
		"memberList": model.GetMemberOfTeam(teamID),
		"leader":     leader,
		"auth":       auth,
		"tab":        setting.TeamsTab,
		"useS3":      setting.UseS3Service(),
	}

	tmpl := template.Must(template.ParseFiles("template/admin_teams/team_edit.tmpl", setting.AdminTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}
