package models

import (
	"time"
)

type Order struct {
	Id              int `gorm:"primaryKey"`
	UserID          int
	OrderStatusID   int `gorm:"default:1"`
	CartID          int
	CreatedAt       time.Time
	DeliveryType    DeliveryType `gorm:"constraint:OnDelete:CASCADE;"`
	PaymentMethodID int
	TotalPrice      float64
	Cancel          bool          `gorm:"default:false"`
	Promo           bool          `gorm:"default:false"`
	PayKeeperInfo   PayKeeperInfo `gorm:"constraint:OnDelete:CASCADE;"`
}

type PayKeeperInfo struct {
	Id        int `gorm:"primaryKey"`
	OrderID   int
	PaymentID int
}

type OnlinePay struct {
	Id          int   `gorm:"primaryKey"`
	OrderID     Order `gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time
	Cancel      bool `gorm:"default:false"`
	PayKeeperID int
}

type PaymentMethod struct {
	Id          int `gorm:"primaryKey"`
	Title       string
	Description string
	Icon        string
	Image       string
	Order       []Order `gorm:"constraint:OnDelete:SET NULL;"`
}

type DeliveryType struct {
	Id              int             `gorm:"primaryKey"`
	CourierDelivery CourierDelivery `gorm:"constraint:OnDelete:CASCADE;"`
	SelfDelivery    SelfDelivery    `gorm:"constraint:OnDelete:CASCADE;"`
	CDEKDelivery    CDEKDelivery    `gorm:"constraint:OnDelete:CASCADE;"`
	OrderID         int
}

type OrderStatus struct {
	Id    int    `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Order Order  `json:"orders" gorm:"constraint:OnDelete:SET NULL;"`
}

type CourierDelivery struct {
	Id              int `gorm:"primaryKey"`
	Address         string
	ApartmentOffice string
	Index           string
	Entrance        string
	Intercom        string
	Floor           string
	Note            string
	DeliveryTypeID  int
}

type SelfDelivery struct {
	Id             int `gorm:"primaryKey"`
	PickUpPointID  int
	DeliveryTypeID int
	//TODO continue logic
}

type CDEKDelivery struct {
	Id                 int `gorm:"primaryKey"`
	DeliveryTypeID     int
	PickUpPointAddress string
	//TODO continue logic
}
