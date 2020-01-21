package service

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/gorilla/sessions"
	"github.com/tribalmedia/vista/model"
	"github.com/tribalmedia/vista/setting"
	"gopkg.in/boj/redistore.v1"
)

// Logger is ...
var Logger *zap.Logger

// ConnectToRedis ...
func ConnectToRedis() *redistore.RediStore {
	var conn setting.RedisConnect
	conn = setting.SetRedisConnect(conn)
	store, err := redistore.NewRediStore(10, "tcp", conn.Host, "", []byte("secret-key"))
	if err != nil {
		Logger.Fatal(err.Error())
	}
	return store
}

// SaveSession ...
func SaveSession(session *sessions.Session, r *http.Request, w http.ResponseWriter) {
	err := session.Save(r, w)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}

//GetSessionMember gets info of user who is logged in
func GetSessionMember(r *http.Request) model.Auth {
	store := ConnectToRedis()
	defer store.Close()
	//get session info
	session, _ := store.Get(r, "vista_member")
	val := session.Values["member"]
	member, _ := val.(model.Auth)
	return member
}
