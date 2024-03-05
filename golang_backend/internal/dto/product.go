package dto

import (
	"time"
)

type Products struct {
	Product    []*Product  `json:"data"`
	Count      *Count      `json:"count"`
	SortParams *SortParams `json:"sort_params"`
}

type Product struct {
	UUID        string    `json:"id"`
	Title       string    `json:"title"`
	VendorCode  string    `json:"vendor_code"`
	Count       string    `json:"count"`
	Quantity    string    `json:"quantity"`
	Image       string    `json:"image"`
	Price       string    `json:"price"`
	BaseUnit    string    `json:"base_unit"`
	Rating      string    `json:"rating"`
	ReviewCount string    `json:"review_count"`
	CreatedAt   time.Time `json:"-"`
	SalesCount  int       `json:"-"`
}

type SortParams struct {
	MaxPrice string `json:"max_price"`
}

//func (p *Product) ImageMediaRoot(mediaRoot string) {
//	if p.Image != "" {
//		p.Image = mediaRoot + p.Image
//	}
//}

type ProductDB struct {
	UUID             string `gorm:"primaryKey"`
	Reviews          []ReviewDB
	Rating           string
	SalesCount       string
	ReviewCount      string
	Sum5             string
	Sum4             string
	Sum3             string
	Sum2             string
	Sum1             string
	RecommendPercent string
}

type ReviewDB struct {
	Id            int
	Body          string
	Rating        int
	ProductDBUUID string `gorm:"column:product_uuid"`
	UserID        int
	Recommend     bool
	CreatedAt     time.Time
	Name          string
	Image         []ReviewImageDB
}

type ReviewImageDB struct {
	ReviewID int
	Image    []byte
}

type Count struct {
	Count string `json:"count" binding:"required"`
}

type ProductInformation struct {
	UUID            string            `json:"id"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	VendorCode      string            `json:"vendor_code"`
	Count           string            `json:"count"`
	Quantity        string            `json:"quantity"`
	Image           string            `json:"image"`
	Price           string            `json:"price"`
	BaseUnit        string            `json:"base_unit"`
	Rating          string            `json:"rating"`
	Files           []*ProductFiles   `json:"files"`
	Characteristics []*Characteristic `json:"characteristics"`
	Category        Category          `json:"category"`
}

type Characteristic struct {
	UUID  string `json:"id"`
	Title string `json:"title"`
	Value string `json:"value"`
}

type ProductFiles struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Document string `json:"document"`
}

func (p *ProductFiles) DocumentRoot(mediaRoot string) {
	if p.Document != "" {
		p.Document = mediaRoot + p.Document
	}
}

//func (p *ProductInformation) ImageMediaRoot(mediaRoot string) {
//	if p.Image != "" {
//		p.Image = mediaRoot + p.Image
//	}
//}

type CartProduct struct {
	UUID       string `json:"id"`
	Title      string `json:"title"`
	Price      string `json:"price"`
	Count      string `json:"count"`
	Quantity   string `json:"quantity"`
	Image      string `json:"image"`
	TotalPrice string `json:"total_price"`
}

//func (p *CartProduct) ImageMediaRoot(mediaRoot string) {
//	if p.Image != "" {
//		p.Image = mediaRoot + p.Image
//	}
//}

type FavouriteProduct struct {
	UUID     string `json:"id"`
	Title    string `json:"title"`
	Price    string `json:"price"`
	Image    string `json:"image"`
	Quantity string `json:"quantity"`
}

//func (p *FavouriteProduct) ImageMediaRoot(mediaRoot string) {
//	if p.Image != "" {
//		p.Image = mediaRoot + p.Image
//	}
//}

type Review struct {
	Body        string   `json:"body" binding:"required"`
	Rating      string   `json:"rating" binding:"required"`
	ProductUUID string   `json:"product_id" binding:"required"`
	Recommend   bool     `json:"recommend"`
	Image       []string `json:"image"`
}

type Stars struct {
	Star5 string `json:"star_5"`
	Star4 string `json:"star_4"`
	Star3 string `json:"star_3"`
	Star2 string `json:"star_2"`
	Star1 string `json:"star_1"`
}

type ProductWithFilesAndCharacteristics struct {
	*ProductInformation
	*ProductFiles
	*Characteristic
}

type GetReview struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Body      string   `json:"body"`
	Rating    string   `json:"rating"`
	CreatedAt string   `json:"created_at"`
	Image     []string `json:"image"`
}

type ProductStatistic struct {
	Recommend      string      `json:"recommendation"`
	ReviewCount    string      `json:"review_count"`
	SalesCount     string      `json:"sales_count"`
	StarsStatistic Stars       `json:"stars_statistic"`
	Rating         string      `json:"rating"`
	Reviews        []GetReview `json:"reviews"`
}

type ProductReviewDB struct {
	ReviewID   int
	Rating     float64
	SalesCount int
}
