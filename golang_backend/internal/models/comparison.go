package models

type Comparison struct {
	Id                 int `gorm:"primaryKey"`
	UserID             int
	ComparisonProducts []ComparisonProduct `gorm:"many2many:comparisonm2ms;constraint:OnDelete:CASCADE;"`
}
