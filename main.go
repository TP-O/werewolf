package main

import (
	"uwwolf/config"
)

func main() {
	// cluster := gocql.NewCluster("192.168.1.1", "192.168.1.2", "192.168.1.3")
	// cluster.Keyspace = "example"
	// session, _ := cluster.CreateSession()

	// session.Query("aaa").GetConsistency().MarshalText()

	println(config.DB().Username)
	println(config.Game().MinCapacity)
}
