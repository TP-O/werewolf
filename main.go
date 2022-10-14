package main

import (
	"fmt"
	"strings"
)

func main() {
	// app.Init()

	a := "aaa"

	fmt.Println(strings.Split(a, "Bearer")[1])
}
