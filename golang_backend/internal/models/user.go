package models

import "time"

type User struct {
	Id int `gorm:"primaryKey;index:user_id_idx"`
	//Phyz
	Name    string
	Surname string
	//Ur
	CompanyName    string
	CompanyAddress string
	Inn            string `gorm:"index:inn_idx"`
	Kpp            string
	ManagerName    string
	//Same
	Phone string `gorm:"index:phone_idx"`

	Cart       []Cart       `gorm:"constraint:OnDelete:CASCADE;"`
	Favourites []Favourite  `gorm:"constraint:OnDelete:CASCADE;"`
	Comparison []Comparison `gorm:"constraint:OnDelete:CASCADE;"`
	Reviews    []Review     `gorm:"constraint:OnDelete:CASCADE;"`
	Orders     []Order      `gorm:"constraint:OnDelete:CASCADE;"`
	Email      Email        `gorm:"constraint:OnDelete:CASCADE;"`
}

type Session struct {
	Id        int       `gorm:"primaryKey"`
	Session   string    `gorm:"index:session_idx"`
	User      User      `gorm:"constraint:OnDelete:CASCADE;"`
	UserID    int       `gorm:"index:user_id_idx"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Email struct {
	Id                      int `gorm:"primaryKey"`
	Email                   string
	CanToSendNews           bool `gorm:"default:true"`
	CanToSendPersonalOffers bool `gorm:"default:true"`
	UserID                  int  `gorm:"unique"`
}
