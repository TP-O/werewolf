package game

type gameManager struct {
	instances []instance
}

var GameManger *gameManager

func init() {
	GameManger = &gameManager{}
}
