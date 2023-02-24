package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"uwwolf/util"
)

func main() {
	cmd := exec.Command(
		"migrate",
		"-database",
		"cassandra://"+
			util.Config().DB.Hosts[0]+
			"/"+
			util.Config().DB.Keyspace+
			"?username="+
			util.Config().DB.Username+
			"&password="+
			util.Config().DB.Password+
			"",
		"-path",
		"db/migration",
		strings.Join(os.Args[1:], " "),
	)
	cmd.Stdin = os.Stdin
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
