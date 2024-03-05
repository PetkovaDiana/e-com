package dto

type UpdateCart struct {
	ProductUUID string `json:"product_id" binding:"required"`
	Count       string `json:"count" binding:"required"`
}

type DeleteCart struct {
	ProductUUID string `json:"product_id" binding:"required"`
}

type Cart struct {
	TotalPrice string         `json:"total_price"`
	Product    []*CartProduct `json:"products"`
}
