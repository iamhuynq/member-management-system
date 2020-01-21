package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tribalmedia/vista/model"
	"github.com/tribalmedia/vista/service"
	"github.com/tribalmedia/vista/setting"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  setting.ReadBufferSize,
	WriteBufferSize: setting.WriteBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan MoveTeamData)      // broadcast channel
var broadcastUser = make(chan []model.Auth)  // broadcast user channel
var arrUser = []model.Auth{}                 // arrUser is array save user online
var broadcastSeat = make(chan MoveSeatData)  // broadcast seat channel

// HandleConnectionsSocket handles websocket requests from the peer
func HandleConnectionsSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Logger.Error(err.Error())
	}
	// Register our new client
	clients[conn] = true

	member := service.GetSessionMember(r)
	arrUser = append(arrUser, member)
	broadcastUser <- arrUser

	for {
		var mtd MoveTeamData
		var msd MoveSeatData
		var incMessage map[string]interface{}
		_, message, err := conn.ReadMessage()
		_ = json.Unmarshal(message, &incMessage)
		if incMessage["MoveType"] != nil {
			switch int(incMessage["MoveType"].(float64)) {
			case setting.MoveOnTeam:
				_ = json.Unmarshal(message, &mtd)
				broadcast <- mtd
			case setting.MoveOnSeat:
				_ = json.Unmarshal(message, &msd)
				broadcastSeat <- msd
			}
		}

		// Handle close event
		if err != nil {
			var index = SliceIndex(len(arrUser), func(i int) bool { return arrUser[i].LoginID == member.LoginID })
			arrUser = Remove(arrUser, index)
			delete(clients, conn)
			broadcastUser <- arrUser
			break
		}
	}
}

//ShowUserOnline sends all user online to clients connected
func ShowUserOnline() {
	for {
		arrUser := <-broadcastUser
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(arrUser)
			if err != nil {
				Logger.Error(err.Error())
				client.Close()
				delete(clients, client)
			}
		}
	}
}

//SliceIndex return index of value
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

//Remove user from online user list...
func Remove(slice []model.Auth, i int) []model.Auth {
	return append(slice[:i], slice[i+1:]...)
}
