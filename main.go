package main

import (
	"context"
	"fmt"
	"log"
	"uwwolf/db"
)

func main() {
	// handler.APIRouter().Run(":" + strconv.Itoa(config.App().ApiPort))
	// grpc.Run()

	// println(db.Client())

	scanner := db.Client().Query(`SELECT id, name FROM phases WHERE id = ?`,
		"1").WithContext(context.Background()).Iter().Scanner()

	for scanner.Next() {
		var (
			id   int
			name string
		)
		err := scanner.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Tweet:", id, name)
	}
	// scanner.Err() closes the iterator, so scanner nor iter should be used afterwards.
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
