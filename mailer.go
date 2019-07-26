package main

import (
	"gopkg.in/gomail.v2"
)

type Mail struct {
	From			string
	To				string
	Subject			string
	Body			string
}

func SendMail(cfg ConfigMail, mail Mail) {
	m := gomail.NewMessage()
	m.SetHeader("From", mail.From)
	m.SetHeader("To", mail.To)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Body)
	
	d := gomail.NewDialer(cfg.Server, cfg.Port, cfg.Username, cfg.Password)
	
	// Send the email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}