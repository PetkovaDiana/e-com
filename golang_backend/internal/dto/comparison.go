package dto

type Comparison struct {
	CategoryUUID      string                    `json:"category_id"`
	CategoryTitle     string                    `json:"category_title"`
	ComparisonProduct []ComparisonProduct       `json:"products"`
	Params            []Param                   `json:"characteristic"`
	Characteristics   []*CategoryCharacteristic `json:"-"`
}

type Param struct {
	Title   string         `json:"param"`
	Product []ProductValue `json:"products"`
}

type ProductValue struct {
	ProductUUID string `json:"product_id"`
	Value       string `json:"value"`
}

type ComparisonProduct struct {
	CategoryUUID    string            `json:"-"`
	CategoryTitle   string            `json:"-"`
	UUID            string            `json:"id"`
	Title           string            `json:"title"`
	Price           string            `json:"price"`
	Image           string            `json:"image"`
	Rating          string            `json:"rating"`
	TotalReviews    string            `json:"total_reviews"`
	Quantity        string            `json:"quantity"`
	Count           string            `json:"count"`
	BaseUnit        string            `json:"base_unit"`
	Characteristics []*Characteristic `json:"-"`
}

type ComparisonProductWithCharacteristic struct {
	*ComparisonProduct
	*Characteristic
}

//func (p *ComparisonProduct) ImageMediaRoot(mediaRoot string) {
//	if p.Image != "" {
//		p.Image = mediaRoot + p.Image
//	}
//}

type DeleteComparison struct {
	ProductUUID string `json:"product_id" binding:"required"`
}

type DeleteComparisonByCategory struct {
	CategoryUUID string `json:"category_id" binding:"required"`
}

type UpdateComparison struct {
	ProductUUID string `json:"product_id" binding:"required"`
}
