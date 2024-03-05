package email

import "clean_arch/internal/dto"

type Email interface {
	SendOrderEmail(emailInfo *dto.EmailOrder) error
	SendCourierOrderEmail(emailInfo *dto.EmailOrder) error
}

func NewEmail(emailConfig *Config) Email {
	return &email{
		emailConfig: emailConfig,
	}
}
