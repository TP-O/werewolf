package db

import (
	"uwwolf/util"

	"github.com/gocql/gocql"
)

var session *gocql.Session

func init() {
	cluster := gocql.NewCluster(util.Config().DB.Hosts...)
	cluster.Keyspace = util.Config().DB.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: util.Config().DB.Username,
		Password: util.Config().DB.Password,
	}

	if ss, err := cluster.CreateSession(); err != nil {
		panic(err)
	} else {
		session = ss
	}
}

func Client() *gocql.Session {
	return session
}
