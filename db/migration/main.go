package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"uwwolf/config"
)

func main() {
	cmd := exec.Command(
		"migrate",
		"-database",
		"cassandra://"+
			config.DB().Hosts[0]+
			":9042/"+
			config.DB().Keyspace+
			"?username="+
			config.DB().Username+
			"&password="+
			config.DB().Password+
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
