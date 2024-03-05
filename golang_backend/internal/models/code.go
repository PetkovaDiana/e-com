package models

type Code struct {
	Id    int `gorm:"primaryKey"`
	Phone string
	Code  string
}
