package dto

type Favourite struct {
	Product []*FavouriteProduct `json:"products"`
}

type DeleteFavourite struct {
	ProductUUID string `json:"product_id" binding:"required"`
}

type UpdateFavourite struct {
	ProductUUID string `json:"product_id" binding:"required"`
}
