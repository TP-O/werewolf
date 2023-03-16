package game

import (
	"testing"
	"uwwolf/util"
)

func TestMain(m *testing.M) {
	util.LoadConfig("../")
	m.Run()
}
