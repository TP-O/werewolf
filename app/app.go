package app

import "uwwolf/app/server"

func Init() {
	go server.StartAPI()
	server.StartSocketIO()
}
