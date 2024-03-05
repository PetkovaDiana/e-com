package phone

import (
	"clean_arch/internal/dto"
	"fmt"
	"net/http"
)

const (
	brandName = "BESM"
	phoneUrl  = "http://api.sms-prosto.ru/?method=push_msg&key=%s&text=%s&phone=%s&sender_name=%s"
)

type phone struct {
	apiKey string
}

func (p *phone) SendCode(code, phone string) error {
	message := fmt.Sprintf("Ваш код - %s", code)

	url := fmt.Sprintf(phoneUrl, p.apiKey, message, phone, brandName)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (p *phone) PhoneOrder(smsOrder *dto.SMSOrder) error {
	var messageData string

	for _, product := range smsOrder.Product {
		messageData += fmt.Sprintf("%s %s шт. %s руб | ", product.Title, product.Quantity, product.Price)
	}
	message := fmt.Sprintf(`Ваш заказ с сайта "Уфа-электро" Номер заказа - %s Товары | Кол-во | Цена: %s Итоговая цена: %s`, smsOrder.OrderID, messageData, smsOrder.TotalPrice)

	url := fmt.Sprintf(phoneUrl, p.apiKey, message, smsOrder.Phone, brandName)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
