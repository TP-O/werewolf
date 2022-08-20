// package seeder

// import (
// 	"errors"

// 	"gorm.io/gorm"
// )

// func Seed(db *gorm.DB) {
// 	seedFactions(db)
// 	seedPhases(db)
// 	seedRoles(db)
// }

// func insert(db *gorm.DB, table interface{}, value interface{}) {
// 	// Only insert if table is empty
// 	if err := db.First(table).Error; errors.Is(err, gorm.ErrRecordNotFound) {
// 		db.Create(value)
// 	}
// }
