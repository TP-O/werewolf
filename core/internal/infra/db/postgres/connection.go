package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"uwwolf/internal/config"

	"github.com/avast/retry-go"
)

func Connect(config config.Postgres) Store {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	))
	if err != nil {
		log.Panic(err)
	}

	if err := retry.Do(
		func() error {
			if err := db.Ping(); err != nil {
				log.Println(err.Error())
				return err
			}
			return nil
		},
		retry.Attempts(10),
		retry.DelayType(retry.RandomDelay),
		retry.MaxJitter(10*time.Second),
	); err != nil {
		log.Panic(err)
	}

	db.SetMaxOpenConns(config.PollSize)
	db.SetMaxIdleConns(config.PollSize)
	db.SetConnMaxIdleTime(0)

	return &store{
		Queries: New(db),
		db:      db,
	}
}
