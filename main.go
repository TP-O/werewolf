package main

import (
	"fmt"
)

type AA map[string]uint

type BB struct {
	aa AA
	bb map[int]int
}

func main() {
	var a *BB
	var b *BB

	fmt.Println(a == b)
}
