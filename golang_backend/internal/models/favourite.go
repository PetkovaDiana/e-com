package models

type Favourite struct {
	Id                int `gorm:"primaryKey"`
	UserID            int
	FavouriteProducts []FavouriteProduct `gorm:"many2many:favouritesm2ms;constraint:OnDelete:CASCADE;"`
}
