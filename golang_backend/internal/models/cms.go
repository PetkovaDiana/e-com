package models

import "time"

type PickUpPoint struct {
	Id                      int `gorm:"primaryKey"`
	Phone1                  string
	Phone2                  string
	Phone3                  string
	Email1                  string
	Email2                  string
	Address                 string
	PickUpPointTimeID       int
	PickUpPointStockTitleID int
	CoordinateX             string
	CoordinateY             string
	SelfDelivery            []SelfDelivery `gorm:"constraint:OnDelete:SET NULL;"`
}

type PickUpPointTime struct {
	Id           int `gorm:"primaryKey"`
	Mon          string
	Tue          string
	Wen          string
	Thu          string
	Fri          string
	Sat          string
	Sun          string
	PickUpPoints PickUpPoint `gorm:"constraint:OnDelete:SET NULL;"`
}

type PickUpPointStockTitle struct {
	Id                           int `gorm:"primaryKey"`
	Title                        string
	PickUpPoints                 []PickUpPoint                 `gorm:"constraint:OnDelete:SET NULL;"`
	PickUpPointStockDescriptions []PickUpPointStockDescription `gorm:"constraint:OnDelete:SET NULL;"`
}

type PickUpPointStockDescription struct {
	Id                      int `gorm:"primaryKey"`
	Description             string
	PickUpPointStockTitleID int
}

type RequestCall struct {
	Id        int    `gorm:"primaryKey"`
	Phone     string `gorm:"type:varchar(25)"`
	Name      string
	UserID    int
	User      User `gorm:"constraint:OnDelete:CASCADE;"`
	Message   string
	CreatedAt time.Time
}

type Vacancy struct {
	Id          int    `gorm:"primaryKey"`
	FirstPhone  string `gorm:"type:varchar(30)"`
	SecondPhone string `gorm:"type:varchar(30)"`
	Title       string
	Email       string `gorm:"type:varchar(100)"`
}

type RequestVacancy struct {
	Id        int    `gorm:"primaryKey"`
	Phone     string `gorm:"type:varchar(30)"`
	Name      string `gorm:"type:varchar(30)"`
	Lastname  string `gorm:"type:varchar(30)"`
	Surname   string `gorm:"type:varchar(30)"`
	Email     string `gorm:"type:varchar(100)"`
	CreatedAt time.Time
	Vacancy   Vacancy `gorm:"constraint:OnDelete:CASCADE;"`
	VacancyID int
	Comment   string
}

type CourierDeliveryInfo struct {
	Id                        int `gorm:"primaryKey"`
	Description               string
	CourierDeliveryTimeInfoID int
}

type CourierDeliveryTimeInfo struct {
	Id                  int                   `gorm:"primaryKey"`
	Mon                 string                `gorm:"type:varchar(30)"`
	Tue                 string                `gorm:"type:varchar(30)"`
	Wen                 string                `gorm:"type:varchar(30)"`
	Thu                 string                `gorm:"type:varchar(30)"`
	Fri                 string                `gorm:"type:varchar(30)"`
	Sat                 string                `gorm:"type:varchar(30)"`
	Sun                 string                `gorm:"type:varchar(30)"`
	CourierDeliveryInfo []CourierDeliveryInfo `gorm:"constraint:OnDelete:SET NULL;"`
}

type CDEKDeliveryInfo struct {
	Id          int `gorm:"primaryKey"`
	Description string
}

type DeliveryTypeInfo struct {
	Id            int `gorm:"primaryKey"`
	Title         string
	Description   string
	Icon          string
	DeliveryPrice float64
	CanDelivery   bool
}

type EmailStatic struct {
	Id           int `gorm:"primaryKey"`
	CarImage     string
	CartImage    string
	LikeImage    string
	LogoImage    string
	CourierEmail string
}

type Requisites struct {
	Id   int `gorm:"primaryKey"`
	Text string
}

type PrivacyPolicy struct {
	Id   int `gorm:"primaryKey"`
	Text string
}
