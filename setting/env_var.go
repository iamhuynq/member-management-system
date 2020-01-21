package setting

import (
	"os"
)

// DBConnect Store db connect info
type DBConnect struct {
	UserName string
	Password string
	Host     string
	Schema   string
}

// RedisConnect store Redis connect info
type RedisConnect struct {
	Host string
}

// DefaultAvatar store path of default avatar for member
type DefaultAvatar struct {
	MemberAvatar string
	TeamAvatar   string
}

// GetDatabaseSourceName return data source name to connect to DB
func GetDatabaseSourceName(conn DBConnect) string {
	return conn.UserName + ":" + conn.Password + "@" + conn.Host + "/" + conn.Schema
}

// SetDBConnect test
func SetDBConnect(conn DBConnect) DBConnect {
	switch os.Getenv("APP_MODE") {
	case "dev":
		// SetDBConnect when app on dev mode
		conn.UserName = os.Getenv("LOCAL_USERNAME")
		conn.Password = os.Getenv("LOCAL_PASSWORD")
		conn.Host = os.Getenv("LOCAL_HOST")
		conn.Schema = os.Getenv("LOCAL_SCHEMA")
	case "prod":
		// SetDBConnect when app on product mode
		conn.UserName = os.Getenv("PROD_USERNAME")
		conn.Password = os.Getenv("PROD_PASSWORD")
		conn.Host = os.Getenv("PROD_HOST")
		conn.Schema = os.Getenv("PROD_SCHEMA")
	default:
		// SetDBConnect when app on dev mode
		conn.UserName = os.Getenv("LOCAL_USERNAME")
		conn.Password = os.Getenv("LOCAL_PASSWORD")
		conn.Host = os.Getenv("LOCAL_HOST")
		conn.Schema = os.Getenv("LOCAL_SCHEMA")
	}

	return conn
}

// SetRedisConnect is ...
func SetRedisConnect(conn RedisConnect) RedisConnect {
	switch os.Getenv("APP_MODE") {
	case "dev":
		// Set Redis connect for local
		conn.Host = os.Getenv("LOCAL_REDIS_HOST")
	case "prod":
		// Set Redis connect for product
		conn.Host = os.Getenv("PROD_REDIS_HOST")
	default:
		// Set Redis connect for local
		conn.Host = os.Getenv("LOCAL_REDIS_HOST")
	}

	return conn
}

// UseS3Service is function check if user upload image to AWS S3
func UseS3Service() bool {
	switch os.Getenv("APP_MODE") {
	case "dev":
		// Upload image to webroot folder
		return false
	case "prod":
		// Upload image to S3
		return true
	default:
		return true
	}
}

// SetDefaultAvatar is ...
func SetDefaultAvatar() DefaultAvatar {
	var defaultAvatar DefaultAvatar
	switch os.Getenv("APP_MODE") {
	case "dev":
		defaultAvatar.MemberAvatar = DefaultLocalMemberAvatar
		defaultAvatar.TeamAvatar = DefaultLocalTeamAvatar
	case "prod":
		defaultAvatar.MemberAvatar = DefaultS3MemberAvatar
		defaultAvatar.TeamAvatar = DefaultS3TeamAvatar
	default:
		defaultAvatar.MemberAvatar = DefaultS3MemberAvatar
		defaultAvatar.TeamAvatar = DefaultS3TeamAvatar
	}

	return defaultAvatar
}
