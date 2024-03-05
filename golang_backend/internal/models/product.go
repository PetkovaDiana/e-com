package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	UUID            uuid.UUID `gorm:"primaryKey;type:uuid;index:prod_uuid_idx"`
	Title           string
	Description     string
	VendorCode      string
	Categories      []Category              `gorm:"many2many:category_product;constraint:OnDelete:CASCADE;"`
	Characteristics []ProductCharacteristic `gorm:"constraint:OnDelete:CASCADE;"`
	Category        Category                `gorm:"constraint:OnDelete:CASCADE;"`
	BaseUnit        string
	//Image            []byte  `gorm:"type:bytea"`
	Image            string  `gorm:"type:text"`
	Price            float64 `gorm:"index:price_idx"`
	CreatedAt        time.Time
	CanToView        bool `gorm:"default:true;index:can_to_view_idx"`
	Quantity         float64
	ProductStatistic ProductStatistic `gorm:"constraint:OnDelete:CASCADE;"`
	ProductFiles     []*ProductFiles  `gorm:"many2many:product_file;constraint:OnDelete:CASCADE;"`
	Reviews          []Review         `gorm:"constraint:OnDelete:CASCADE;"`
	CategoryID       uuid.UUID
}

type ProductFiles struct {
	Id       int    `gorm:"primaryKey"`
	Title    string `gorm:"type:VARCHAR(255)"`
	Document string
	Products []*Product `gorm:"many2many:product_file;constraint:OnDelete:CASCADE;"`
}

type ProductStatistic struct {
	Id                 int       `gorm:"primaryKey"`
	Rating             float64   `gorm:"index:rating_idx"`
	SalesCount         int       `gorm:"default:0;index:sales_count_idx"`
	RequestDetailCount int       `gorm:"index:req_detail_count_idx"`
	ProductUUID        uuid.UUID `gorm:"constraint:unique;type:uuid;unique;index:prod_uuid_idx"`
}

type Characteristic struct {
	UUID                  uuid.UUID `gorm:"primaryKey;type:uuid"`
	Title                 string
	Categories            []Category              `gorm:"many2many:category_characteristic;constraint:OnDelete:CASCADE;"`
	ProductCharacteristic []ProductCharacteristic `gorm:"constraint:OnDelete:CASCADE;"`
}

type ProductCharacteristic struct {
	ProductUUID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	CharacteristicUUID uuid.UUID `gorm:"primaryKey;type:uuid"`
	Value              string
	Product            Product        `gorm:"constraint:OnDelete:CASCADE;"`
	Characteristic     Characteristic `gorm:"constraint:OnDelete:CASCADE;"`
}

type Review struct {
	Id           int    `gorm:"primaryKey"`
	Body         string `gorm:"type:text"`
	Rating       int
	ProductUUID  uuid.UUID `gorm:"type:uuid;index:prod_uuid_idx"`
	UserID       int
	Recommend    bool
	CreatedAt    time.Time
	ReviewPhotos []ReviewPhotos
}

type ReviewPhotos struct {
	Id       int    `gorm:"primaryKey"`
	Image    []byte `gorm:"type:bytea"`
	ReviewID int
}

type CartProduct struct {
	Id          int       `gorm:"primaryKey"`
	Product     Product   `gorm:"constraint:OnDelete:CASCADE;"`
	ProductUUID uuid.UUID `gorm:"type:uuid;index:prod_uuid_idx"`
	Count       int
	TotalPrice  float64
	Carts       []Cart `gorm:"many2many:cartm2ms;constraint:OnDelete:CASCADE;"`
}

type FavouriteProduct struct {
	Id          int         `gorm:"primaryKey"`
	Product     Product     `gorm:"constraint:OnDelete:CASCADE;"`
	ProductUUID uuid.UUID   `gorm:"type:uuid;index:prod_uuid_idx"`
	Favourites  []Favourite `gorm:"many2many:favouritesm2ms;constraint:OnDelete:CASCADE;"`
}

type ComparisonProduct struct {
	Id           int       `gorm:"primaryKey"`
	Product      Product   `gorm:"constraint:OnDelete:CASCADE;"`
	ProductUUID  uuid.UUID `gorm:"type:uuid;index:prod_uuid_idx"`
	Category     Category  `gorm:"constraint:OnDelete:CASCADE;"`
	CategoryUUID uuid.UUID
	Comparison   []Comparison `gorm:"many2many:comparisonm2ms;constraint:OnDelete:CASCADE;"`
}

type PromoCode struct {
	Id              int `gorm:"primaryKey"`
	PromoCode       string
	DiscountPercent int
	DiscountSum     int
	NumberOfUses    int
	ExpiresAt       time.Time
	CreatedAt       time.Time
}
