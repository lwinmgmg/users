package services

import (
	"fmt"
	"net/smtp"
)

var (
	MailSender *MailService
)

type MailService struct {
	senderMail string
	password   string
	Host       string
	Port       string
}

func (sender *MailService) getAuth() smtp.Auth {
	return smtp.PlainAuth("", sender.senderMail, sender.password, sender.Host)
}

func (sender *MailService) Send(message string, recipient []string) error {
	// err := smtp.SendMail(fmt.Sprintf("%v:%v", sender.Host, sender.Port), sender.getAuth(), sender.senderMail, recipient, []byte(message))
	// if err != nil {
	// 	fmt.Printf("Error on sending email %v\n", err)
	// }
	fmt.Println(message, recipient)
	return nil
}

func NewMailService(email, password, host, port string) *MailService {
	return &MailService{
		senderMail: email,
		password:   password,
		Host:       host,
		Port:       port,
	}
}
