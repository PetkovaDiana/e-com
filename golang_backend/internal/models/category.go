package models

import "github.com/google/uuid"

type Category struct {
	UUID            uuid.UUID `gorm:"primaryKey;type:uuid"`
	Title           string
	Image           string
	CanToView       bool             `gorm:"default:true"`
	Level           int              `gorm:"default:0"`
	ParentUuid      uuid.UUID        `gorm:"type:uuid"`
	Products        []Product        `gorm:"many2many:category_product;constraint:OnDelete:CASCADE;"`
	Characteristics []Characteristic `gorm:"many2many:category_characteristic;constraint:OnDelete:CASCADE;"`
	Product         []Product        `gorm:"constraint:OnDelete:CASCADE;"`
}
