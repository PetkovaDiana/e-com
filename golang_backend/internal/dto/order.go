package dto

import (
	"time"
)

type Order struct {
	DeliveryType    DeliveryType `json:"delivery_type" binding:"required"`
	PaymentMethodID int          `json:"payment_method_id" binding:"required"`
	PromoCode       string       `json:"promo_code"`
}

type OnlineOrderChecker struct {
	ID       string  `json:"id"`
	Sum      float64 `json:"sum"`
	ClientID string  `json:"clientid"`
	OrderID  string  `json:"orderid"`
	Key      string  `json:"key"`
}

type PayKeeperOrderCancel struct {
	PaymentID     int
	Partial       bool
	SecurityToken string
}

type PayKeeperReceiptSender struct {
	InvoiceID     string
	SecurityToken string
}

type PayKeeperReceiptCreateData struct {
	PayAmount   float64
	ClientID    int
	ClientEmail string
	ClientPhone string
	Expiry      string
	Token       string
	ServiceName string
}

type DeliveryType struct {
	CourierDelivery CourierDelivery `json:"courier_delivery"`
	SelfDelivery    SelfDelivery    `json:"self_delivery"`
	CDEKDelivery    CDEKDelivery    `json:"cdek_delivery"`
}

type PaymentMethod struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Image       string `json:"image"`
}

func (p *PaymentMethod) ImageMediaRoot(mediaRoot string) {
	if p.Icon != "" {
		p.Icon = mediaRoot + p.Icon
	}
	if p.Image != "" {
		p.Image = mediaRoot + p.Image
	}
}

type PickUpPoint struct {
	Id               string                `json:"id"`
	Phone1           string                `json:"phone_1"`
	Phone2           string                `json:"phone_2"`
	Phone3           string                `json:"phone_3"`
	Email1           string                `json:"email_1"`
	Email2           string                `json:"email_2"`
	Address          string                `json:"address"`
	PickUpPointTime  PickUpPointTime       `json:"time"`
	PickUpPointStock PickUpPointStockTitle `json:"stock"`
	Coordinates      Coordinates           `json:"coordinates"`
}

type Coordinates struct {
	CoordinateX string `json:"x"`
	CoordinateY string `json:"y"`
}

type PickUpPointTime struct {
	Mon string `json:"mon"`
	Tue string `json:"tue"`
	Wen string `json:"wen"`
	Thu string `json:"thu"`
	Fri string `json:"fri"`
	Sat string `json:"sat"`
	Sun string `json:"sun"`
}

type PickUpPointStockTitle struct {
	Title       string                        `json:"title"`
	Description []PickUpPointStockDescription `json:"description"`
}

type PickUpPointStockDescription struct {
	Description string `json:"description"`
}

type CourierDelivery struct {
	Address         string `json:"address"`
	ApartmentOffice string `json:"apartment_office"`
	Index           string `json:"index"`
	Entrance        string `json:"entrance"`
	Intercom        string `json:"intercom"`
	Floor           string `json:"floor"`
	Note            string `json:"note"`
}

type SelfDelivery struct {
	PickUpPointsID int `json:"pick_up_points_id"`
	//TODO continue logic
}

type CDEKDelivery struct {
	//TODO continue logic
	PickUpPointAddress string `json:"pick_up_point_address"`
}

type OrderProduct struct {
	UUID       string `json:"id"`
	Title      string `json:"title"`
	Price      string `json:"price"`
	Image      string `json:"image"`
	TotalPrice string `json:"total_price"`
	Quantity   string `json:"quantity"`
	Article    string `json:"-"`
	Unit       string `json:"-"`
}

//func (p *OrderProduct) ImageMediaRoot(mediaRoot string) {
//	if p.Image != "" {
//		p.Image = mediaRoot + p.Image
//	}
//}

type GetOrder struct {
	Id          string          `json:"id"`
	Products    []*OrderProduct `json:"products"`
	CreatedAtDB time.Time       `json:"-" gorm:"column:created_at_db"`
	CreatedAt   string          `json:"created_at"`
	OrderStatus string          `json:"status"`
	TotalPrice  string          `json:"total_price"`
	Cancel      bool            `json:"canceled"`
	FullCount   string          `json:"-"`
	TotalOrders int             `json:"-"`
}

func (p *GetOrder) TimeFormatter(timeFormat string) string {
	p.CreatedAt = p.CreatedAtDB.Format(timeFormat)
	return p.CreatedAt
}

type UserOrdersProducts struct {
	UUID        string    `json:"id"`
	Title       string    `json:"title"`
	Image       string    `json:"image"`
	Count       string    `json:"count"`
	TotalPrice  string    `json:"total_price"`
	Rating      string    `json:"rating"`
	Quantity    string    `json:"quantity"`
	FullCount   string    `json:"-"`
	CreatedAtDB time.Time `json:"-" gorm:"column:created_at_db"`
	CreatedAt   string    `json:"created_at"`
}

//func (p *UserOrdersProducts) ImageMediaRoot(mediaRoot string) {
//	if p.Image != "" {
//		p.Image = mediaRoot + p.Image
//	}
//}

func (p *UserOrdersProducts) TimeFormatter(timeFormat string) string {
	p.CreatedAt = p.CreatedAtDB.Format(timeFormat)
	return p.CreatedAt
}

type GetReceipt struct {
	PromoCode    string `json:"promo_code"`
	DeliveryType string `json:"delivery_type" binding:"required"`
}

type Receipt struct {
	ProductCount  string `json:"products_count"`
	CartPrice     string `json:"cart_price"`
	Sale          string `json:"sale"`
	DeliveryPrice string `json:"delivery_price"`
	FinalPrice    string `json:"final_price"`
}

type OrderProductsDB struct {
	ProductUUID string
	CartID      int
}

type CartProductsDB struct {
	ProductUUID string
	CartID      int
}

type SMSOrder struct {
	Product    []*OrderProduct
	TotalPrice string
	Phone      string
	OrderID    string
}

type EmailStatic struct {
	CarImage     string
	CartImage    string
	LikeImage    string
	LogoImage    string
	CourierEmail string
}

type EmailOrder struct {
	Address         string
	CourierDelivery CourierDelivery
	PaymentMethod   string
	OrderID         string
	Email           string
	Phone           string
	Inn             string
	Product         []*OrderProduct
	EmailStatic     *EmailStatic
	TotalPrice      string
	Sale            string
	DeliveryPrice   string
	FIO             string
	ManagerName     string
}

type Order1C struct {
	Data struct {
		Id           int              `json:"id"`
		CreatedAt    time.Time        `json:"created_at"`
		Name         string           `json:"name"`
		Phone        string           `json:"phone"`
		Address      string           `json:"address"`
		Inn          string           `json:"inn"`
		CustomerType string           `json:"customer_type"`
		Products     []Order1CProduct `json:"products"`
		Comment      string           `json:"comment"`
		Branch       int              `json:"branch"`
	} `json:"data"`
}

type Order1CProduct struct {
	UUID     string  `json:"uuid"`
	Quantity float64 `json:"quantity"`
}

type InvoiceData struct {
	Id        string `json:"number"`
	CreatedAt string `json:"date"`
	Url       string `json:"-"`
	Cart      struct {
		CartProducts   []InvoiceDataCartProduct `json:"cart_products"`
		Nds            float64                  `json:"nds"`
		TotalPrice     float64                  `json:"total_price"`
		ProductCounter int                      `json:"product_counter"`
	} `json:"cart"`
	Customer InvoiceDataCustomer `json:"customer"`
}

type InvoiceDataCartProduct struct {
	Id        int     `json:"id"`
	Title     string  `json:"title"`
	Article   string  `json:"article"`
	Count     float64 `json:"count"`
	Unit      string  `json:"unit"`
	SoloPrice float64 `json:"price"`
	FullPrice float64 `json:"cost_price"`
}

type InvoiceDataCustomer struct {
	Title   string `json:"title"`
	Inn     string `json:"inn"`
	Kpp     string `json:"kpp"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

type OrderData struct {
	EmailOrder   EmailOrder
	SMSOrder     SMSOrder
	Order1C      Order1C
	ProductUUIDs []string
	PaymentID    int
	NewOrderID   int
	ReceiptData  PayKeeperReceiptCreateData
	InvoiceData  InvoiceData
}

type ReturnOrder struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}
