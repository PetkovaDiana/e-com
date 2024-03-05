package phone

import (
	"clean_arch/internal/dto"
)

type Phone interface {
	SendCode(code, phone string) error
	PhoneOrder(smsOrder *dto.SMSOrder) error
}

func NewPhone(apiKey string) Phone {
	return &phone{
		apiKey: apiKey,
	}
}
