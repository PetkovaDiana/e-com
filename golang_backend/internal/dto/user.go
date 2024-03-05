package dto

type RegisterUser struct {
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	CompanyName    string `json:"company_name"`
	CompanyAddress string `json:"company_address"`
	Kpp            string `json:"kpp"`
	Inn            string `json:"inn"`
	ManagerName    string `json:"manager_name"`
	Phone          string `json:"phone" binding:"required"`
	Code           string `json:"code" binding:"required"`
}

type UserAuth struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type Code struct {
	Phone string `json:"phone" binding:"required"`
}

type CodeGenerate struct {
	Phone string `json:"phone" binding:"required"`
}

type CodeResponse struct {
	Code string `json:"code"`
}

type UserData struct {
	Id                      string `json:"id"`
	Name                    string `json:"name"`
	Surname                 string `json:"surname"`
	CompanyName             string `json:"company_name"`
	CompanyAddress          string `json:"company_address"`
	Kpp                     string `json:"kpp"`
	Inn                     string `json:"inn"`
	ManagerName             string `json:"manager_name"`
	Phone                   string `json:"phone"`
	Email                   string `json:"email"`
	CanToSendNews           bool   `json:"can_to_send_news"`
	CanToSendPersonalOffers bool   `json:"can_to_send_personal_offers"`
}

type UpdateEmail struct {
	Email string `json:"email" binding:"required"`
}

type UpdateManagerName struct {
	ManagerName string `json:"manager_name" binding:"required"`
}

type CanToSendEmail struct {
	CanToSendNews           bool `json:"can_to_send_news"`
	CanToSendPersonalOffers bool `json:"can_to_send_personal_offers"`
}

type UserOrderData struct {
	ID             int
	Inn            string
	Name           string
	CompanyName    string
	CompanyAddress string
	Kpp            string
	ManagerName    string
	Phone          string
	Email          string
	PaymentID      int
	PaymentMethod  string
}

type SiteReview struct {
	Rating  int    `json:"rating" binding:"required"`
	Comment string `json:"comment" binding:"required"`
}
