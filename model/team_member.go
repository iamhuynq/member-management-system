package model

import (
	"database/sql"
)

// TeamMember is ...
type TeamMember struct {
	ID       int `db:"id"`
	TeamID   int `db:"team_id"`
	MemberID int `db:"member_id"`
	Position int `db:"position"`
}

// GetTeamByMember is function get all team have that member by member ID and return as an string array
func GetTeamByMember(id int) ([]string, []int) {
	var teamListName []string
	var teamListID []int
	rows, _ := DB.Queryx("SELECT id, team_id, member_id, position FROM team_member WHERE member_id=? ORDER BY id DESC", id)
	for rows.Next() {
		var teamMember TeamMember
		err := rows.StructScan(&teamMember)

		if err != nil {
			Logger.Error(err.Error())
		}

		teamName := GetTeamByID(teamMember.TeamID)

		teamListName = append(teamListName, teamName.Name)
		teamListID = append(teamListID, teamName.ID)
	}

	return teamListName, teamListID
}

// SaveTeamMember is ...
func SaveTeamMember(teamID, memberID int) {
	query := `INSERT INTO team_member (team_id,member_id)
			VALUES (?, ?)`
	_, err := DB.Exec(query, teamID, memberID)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

// DeleteTeamMember is ...
func DeleteTeamMember(teamID, memberID int) {
	_, err := DB.Exec("DELETE FROM team_member WHERE team_id = ? AND member_id = ?", teamID, memberID)
	if err != nil {
		Logger.Error(err.Error())
	}
}

// EditTeamOfMember is
// Delete old team then update new Team
func EditTeamOfMember(memberID int, teamsID []int) {
	var oldTeams []int
	rows, _ := DB.Queryx("SELECT id, team_id, member_id, position FROM team_member WHERE member_id=? ORDER BY id DESC", memberID)
	for rows.Next() {
		var teamMember TeamMember
		err := rows.StructScan(&teamMember)
		if err != nil {
			Logger.Error(err.Error())
		}

		// check if old team exist in new teams
		// if exist => not delete
		// if not exist => delete
		isDelete := true
		for _, teamID := range teamsID {
			if teamMember.TeamID == teamID {
				isDelete = false
				// add team to old teams list
				oldTeams = append(oldTeams, teamID)
			}
		}
		if isDelete {
			DeleteTeamMember(teamMember.TeamID, teamMember.MemberID)
		}
	}

	for _, teamID := range teamsID {
		// check if teamID is a new team
		// if is new => save
		isNew := true
		for _, oldTeam := range oldTeams {
			if teamID == oldTeam {
				isNew = false
			}
		}

		if isNew && !CheckTeamMemberExist(teamID, memberID) {
			SaveTeamMember(teamID, memberID)
		}
	}
}

// EditTeamLeader is ...
func EditTeamLeader(teamID, leaderID int) {
	rows, _ := DB.Queryx("SELECT id, team_id, member_id, position FROM team_member WHERE team_id=? ", teamID)
	for rows.Next() {
		var teamMember TeamMember
		err := rows.StructScan(&teamMember)
		if err != nil {
			Logger.Error(err.Error())
		}
		if teamMember.MemberID == leaderID && teamMember.Position == 1 {
			_, err := DB.Exec("UPDATE team_member SET team_member.position = 2 WHERE team_member.team_id = ? AND team_member.member_id = ?", teamID, teamMember.MemberID)
			if err != nil {
				Logger.Fatal(err.Error())
			}
		}
		if teamMember.MemberID != leaderID && teamMember.Position == 2 {
			_, err := DB.Exec("UPDATE team_member SET team_member.position = 1 WHERE team_member.team_id = ? AND team_member.member_id = ?", teamID, teamMember.MemberID)
			if err != nil {
				Logger.Fatal(err.Error())
			}
		}
	}
}

// CheckTeamMemberExist is ...
func CheckTeamMemberExist(teamID, memberID int) bool {
	var count int
	err := DB.Get(&count, "SELECT count(id) FROM team_member WHERE team_id = ? AND member_id = ?", teamID, memberID)
	if err != nil {
		Logger.Error(err.Error())
	}

	if count != 0 {
		return true
	}
	return false
}

// MoveMember is ...
func MoveMember(memberID, teamID, teamIDPre int) bool {
	if CheckDuplicateMember(memberID, teamID) {
		_, sqlErr := DB.Exec("UPDATE team_member SET team_member.team_id=?, team_member.position = '1' WHERE team_member.team_id=? AND team_member.member_id=?", teamID, teamIDPre, memberID)
		if sqlErr == nil {

			return true
		}
		Logger.Fatal(err.Error())
	}

	return false
}

// CheckDuplicateMember checks one team has many same members
func CheckDuplicateMember(memberID, teamID int) bool {
	var check int
	err := DB.Get(&check, "SELECT team_member.id FROM team_member WHERE team_member.member_id =? AND team_member.team_id = ?", memberID, teamID)
	if err == sql.ErrNoRows {
		return true
	} else if err != nil {
		Logger.Error(err.Error())
	}

	return false
}

// LeaderOfTeam is ...
func LeaderOfTeam(teamID int) map[string]interface{} {
	teamMember := TeamMember{}
	err := DB.Get(&teamMember, "SELECT id, team_id, member_id, position FROM team_member WHERE team_member.position != 1 AND team_member.team_id=?", teamID)
	member := GetMemberByID(teamMember.MemberID)
	position := teamMember.Position
	leader := map[string]interface{}{
		"Leader":   member.Name,
		"leaderID": member.ID,
		"Position": position,
	}
	if err != nil && err != sql.ErrNoRows {
		Logger.Error(err.Error())
	}
	return leader
}

// CheckTeamHasLeader is ...
func CheckTeamHasLeader(teamID int) bool {
	var id int
	err := DB.Get(&id, "SELECT id FROM team_member WHERE team_member.position != 1 AND team_member.team_id=?", teamID)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		Logger.Error(err.Error())
	}
	return true
}

// CheckLeader is ...
func CheckLeader(teamID, memberID int) bool {
	var id int
	err := DB.Get(&id, "SELECT team_member.id FROM team_member WHERE team_member.team_id = ? AND team_member.member_id = ? AND team_member.position = 2", teamID, memberID)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		Logger.Error(err.Error())
	}
	return true
}

// AddLeaderNewTeam is ...
func AddLeaderNewTeam(leaderID, teamID int) {
	query := `INSERT INTO team_member (team_id, member_id, position)
			VALUES (?, ?, 2)`
	_, err := DB.Exec(query, teamID, leaderID)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}
