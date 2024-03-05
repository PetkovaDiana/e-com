package dto

type RequestCall struct {
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	Surname     string `json:"-"`
	ManagerName string `json:"manager_name"`
	Message     string `json:"message"`
	Email       string `json:"email"`
	CompanyName string `json:"-"`
	Inn         string `json:"-"`
}

type CourierDeliveryInfo struct {
	Description             string                  `json:"description"`
	CourierDeliveryTimeInfo CourierDeliveryTimeInfo `json:"time"`
}

type CourierDeliveryTimeInfo struct {
	Mon string `json:"mon"`
	Tue string `json:"tue"`
	Wen string `json:"wen"`
	Thu string `json:"thu"`
	Fri string `json:"fri"`
	Sat string `json:"sat"`
	Sun string `json:"sun"`
}

type CDEKDeliveryInfo struct {
	Description string
}

type DeliveryTypeInfo struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	CanDelivery bool   `json:"can_delivery"`
}

func (p *DeliveryTypeInfo) ImageMediaRoot(mediaRoot string) {
	if p.Icon != "" {
		p.Icon = mediaRoot + p.Icon
	}
}

type RequestVacancy struct {
	Phone        string `json:"phone"`
	Name         string `json:"name"`
	Lastname     string `json:"last_name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	VacancyID    string `json:"vacancy_id"`
	Comment      string `json:"comment"`
	VacancyTitle string `json:"-"`
}

type Vacancy struct {
	Id          string `json:"id"`
	FirstPhone  string `json:"first_phone"`
	SecondPhone string `json:"second_phone"`
	Title       string `json:"title"`
	Email       string `json:"email"`
}

type Requisites struct {
	Text string `json:"body"`
}

type PrivacyPolicy struct {
	Text string `json:"body"`
}
