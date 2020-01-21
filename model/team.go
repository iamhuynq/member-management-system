package model

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/tribalmedia/vista/setting"
)

//Team is ...
type Team struct {
	ID          int            `db:"id"`
	Name        string         `db:"name"`
	PictureURL  sql.NullString `db:"picture_url"`
	Description sql.NullString `db:"description"`
	Created     time.Time      `db:"created"`
	Modified    time.Time      `db:"modified"`
}

// ListTeam is ...
type ListTeam struct {
	List []Team
}

// GetAllTeam is function that get all team in database with id and name
func GetAllTeam() ListTeam {
	team := Team{}
	var teams ListTeam
	rows, _ := DB.Queryx("SELECT id, name from teams")
	for rows.Next() {
		err := rows.StructScan(&team)
		if err != nil {
			Logger.Error(err.Error())
		}
		teams.List = append(teams.List, team)
	}

	return teams
}

// GetTeamByID gets team by team's id
// It returns the team
func GetTeamByID(id int) Team {
	team := Team{}
	err := DB.Get(&team, "SELECT id, name, picture_url, description, created, modified FROM teams WHERE id =?", id)
	if err != nil && err != sql.ErrNoRows {
		Logger.Error(err.Error())
	}

	return team
}

// EditTeam updates team info
// It returns true if save success
func EditTeam(data Team) {
	query := "UPDATE teams SET name=:name, description=:description, picture_url=:picture_url WHERE id=:id"

	_, err := DB.NamedExec(query, data)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// SaveTeam is function that save a new team to database
// It return true if save success
func SaveTeam(team Team) int64 {
	res, err := DB.NamedExec("INSERT INTO teams(name,picture_url,description) VALUES(:name, :picture_url, :description)", team)
	if err != nil {
		Logger.Fatal(err.Error())
	}
	idTeam, _ := res.LastInsertId()

	return idTeam
}

// ValidateTeam is function validate team data before save
// return map incase there are more validate in future
func ValidateTeam(team Team, contentLength int64) map[string]string {
	err := map[string]string{}

	// error message when violate validation
	if len(team.Name) > setting.NameMaxLength {
		err["name"] = "Team's name can't be more than " + strconv.Itoa(setting.NameMaxLength) + " characters."
	}

	if team.Name == "" {
		err["name"] = "Team's name can't be empty."
	}

	if len(team.Description.String) > setting.DescriptionMaxLength {
		err["description"] = "Your description can't be more than " + strconv.Itoa(setting.DescriptionMaxLength) + " characters."
	}

	if contentLength > setting.FileMaxSize {
		err["photo"] = "File size must be less than 3 MB"
	}

	return err
}

// GetTeams is function get list team
func GetTeams() ListTeam {
	rows, _ := DB.Queryx("SELECT id, name, picture_url, description, created, modified FROM teams ORDER BY teams.id DESC")
	var listTeam ListTeam
	for rows.Next() {
		var team Team
		err := rows.StructScan(&team)
		if err != nil {
			Logger.Error(err.Error())
		}
		listTeam.List = append(listTeam.List, team)
	}
	return listTeam
}

// GetTeamAndMems is function get list team with members
func GetTeamAndMems() []map[string]interface{} {
	rows, _ := DB.Queryx("SELECT id, name, picture_url, description, created, modified FROM teams")
	var listFull []map[string]interface{}
	for rows.Next() {
		var team Team
		err := rows.StructScan(&team)
		if err != nil {
			Logger.Error(err.Error())
		}

		m := map[string]interface{}{
			"Team":   team,
			"Member": GetMemberOfTeam(team.ID),
		}
		listFull = append(listFull, m)
	}

	return listFull
}
