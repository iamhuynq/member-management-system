package controller

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/tribalmedia/vista/model"
	"github.com/tribalmedia/vista/service"
	"github.com/tribalmedia/vista/setting"
)

// MoveTeamData save position of member when drag drop
type MoveTeamData struct {
	TeamID    string `json:"team_id"`
	TeamIDPre string `json:"team_id_previous"`
	MemberID  string `json:"member_id"`
	UserRole  string `json:"user_role"`
	UserID    string `json:"user_id"`
	MoveType  int
}

// TeamsTop is ...
func TeamsTop(w http.ResponseWriter, r *http.Request) {
	auth := service.GetSessionMember(r)
	templateData := map[string]interface{}{
		"teamData": model.GetTeamAndMems(),
		"title":    "Team",
		"auth":     auth,
		"useS3":    setting.UseS3Service(),
	}

	tmpl := template.Must(template.ParseFiles("template/team/team.tmpl", setting.UserTemplate))
	if err := tmpl.ExecuteTemplate(w, "base", templateData); err != nil {
		Logger.Error(err.Error())
	}
}

//HandleDataTeam sends msg to clients connected use websocket
func HandleDataTeam() {
	for {
		// Grab the next message from the broadcast channel
		mtd := <-broadcast
		if mtd.UserRole == strconv.Itoa(setting.AdminRoleType) {
			// Send it out to every client that is currently connected
			for client := range clients {
				err := client.WriteJSON(mtd)
				if err != nil {
					Logger.Error(err.Error())
					client.Close()
					delete(clients, client)
				}
			}
		} else {
			Logger.Info("Permission denied")
		}

	}
}
