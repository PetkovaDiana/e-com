package email

import (
	"bytes"
	"clean_arch/internal/dto"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

type Config struct {
	SenderEmail       string
	SenderAppPassword string
	SmtpHost          string
	SmtpPort          string
	DirectorEmail     string
	CourierEmail      string
	MimeHeaders       string
}

type email struct {
	emailConfig *Config
}

func (e *email) SendOrderEmail(emailInfo *dto.EmailOrder) error {
	var body bytes.Buffer

	auth := smtp.PlainAuth("", e.emailConfig.SenderEmail, e.emailConfig.SenderAppPassword, e.emailConfig.SmtpHost)
	address := e.emailConfig.SmtpHost + ":" + e.emailConfig.SmtpPort

	body.Write([]byte(fmt.Sprintf("Subject:Ваш заказ с сайта: «Уфа-электро» \n%s\n\n", e.emailConfig.MimeHeaders)))

	var templateUser *template.Template
	var templatePath string

	switch {
	case emailInfo.Inn != "":
		templatePath = "templates/jurdical_order.html"
	default:
		templatePath = "templates/fizik_order.html"
	}

	templateUser, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := templateUser.Execute(&body, emailInfo); err != nil {
		log.Println(err)
		return err
	}

	if err := smtp.SendMail(address, auth, e.emailConfig.SenderEmail, []string{emailInfo.Email, e.emailConfig.DirectorEmail}, body.Bytes()); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (e *email) SendCourierOrderEmail(emailInfo *dto.EmailOrder) error {
	var body bytes.Buffer

	auth := smtp.PlainAuth("", e.emailConfig.SenderEmail, e.emailConfig.SenderAppPassword, e.emailConfig.SmtpHost)
	address := e.emailConfig.SmtpHost + ":" + e.emailConfig.SmtpPort
	subject := "Subject: Информация о заказе с сайта: «Уфа-электро» \n" + e.emailConfig.MimeHeaders
	body.Write([]byte(subject))

	templateCourier, err := template.ParseFiles("templates/courier_order.html")
	if err != nil {
		return fmt.Errorf("error parsing email template: %w", err)
	}
	if err := templateCourier.Execute(&body, emailInfo); err != nil {
		return fmt.Errorf("error executing email template: %w", err)
	}

	err = smtp.SendMail(address, auth, e.emailConfig.SenderEmail, []string{emailInfo.EmailStatic.CourierEmail}, body.Bytes())
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}
