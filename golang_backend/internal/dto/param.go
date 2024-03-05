package dto

type Params struct {
	Limit       int
	Page        int
	PriceMin    float64
	PriceMax    float64
	RatingMin   float64
	RatingMax   float64
	Sort        []string
	CatId       string
	ProductUUID []string
	Difference  string
	NotNull     string
}
