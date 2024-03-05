package models

type Cart struct {
	Id           int           `gorm:"primaryKey;index:id_idx"`
	UserID       int           `gorm:"index:user_id_idx"`
	InOrder      bool          `gorm:"default:false"`
	CartProducts []CartProduct `gorm:"many2many:cartm2ms;constraint:OnDelete:CASCADE;"`
	Order        Order         `gorm:"constraint:OnDelete:CASCADE;"`
}
