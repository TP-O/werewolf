package itf

type IGame interface {
	NextTurn()
	Pipe(pub *chan string)
	NumberOfVillagers() uint
	NumberOfWerewolves() uint
}
