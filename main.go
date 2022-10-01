package main

import (
	"uwwolf/app/server"
)

func main() {
	go server.StartSocketIO()
	server.StartAPI()
}
