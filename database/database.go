package database

import "gorm.io/gorm"

func DB() *gorm.DB {
	return dbInstance
}
