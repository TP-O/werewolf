package database

import (
	"strconv"
	"uwwolf/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var client *gorm.DB

func init() {
	dsn := "host=" + config.DB.Host +
		" user=" + config.DB.Username +
		" password=" + config.DB.Password +
		" dbname=" + config.DB.Name +
		" port=" + strconv.Itoa(config.DB.Port) +
		" sslmode=disable TimeZone=Asia/Ho_Chi_Minh"

	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		panic(err)
	} else {
		client = db
	}
}

func Client() *gorm.DB {
	return client
}
