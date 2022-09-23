package model

import (
	"time"

	"gorm.io/gorm"

	"uwwolf/app/types"
)

type Faction struct {
	Id        types.FactionId `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string          `gorm:"type:varchar(50);unique" json:"name"`
	CreatedAt time.Time       `gorm:"" json:"createdAt"`
	UpdatedAt time.Time       `gorm:"" json:"updatedAt"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"deletedAt"`
}
