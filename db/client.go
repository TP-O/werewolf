package db

import (
	"uwwolf/config"

	"github.com/gocql/gocql"
)

var session *gocql.Session

func init() {
	cluster := gocql.NewCluster(config.DB().Hosts...)
	cluster.Keyspace = config.DB().Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.DB().Username,
		Password: config.DB().Password,
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
