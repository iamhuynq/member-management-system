package model

import (
	"github.com/jmoiron/sqlx"
	"github.com/tribalmedia/vista/setting"
)

// DB is connection to database
var DB *sqlx.DB
var err error

// ConnectDb is function connect to database
func ConnectDb() {
	var conn setting.DBConnect
	conn = setting.SetDBConnect(conn)
	databaseSource := setting.GetDatabaseSourceName(conn)
	option := "?parseTime=true"
	DB, err = sqlx.Connect("mysql", databaseSource+option)
	if err != nil {
		Logger.Fatal(err.Error())
	}
}
