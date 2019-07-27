package main

import (
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

type Mail struct {
	From    string
	To      []string
	Subject string
	Body    string
}

var cfg ConfigMail
var mailList map[int]Mail
var mailQueu map[int]int64

func InitMailer(config ConfigMail) {
	cfg = config
	mailList = make(map[int]Mail)
	mailQueu = make(map[int]int64)
	fmt.Println("INIT")
}

func MailSender() {
	for range time.Tick(1 * time.Second) {
		for k, v := range mailQueu {
			interval := int(time.Now().Unix() - v)
			if interval > cfg.Interval {
				fmt.Println("SEND MAIL", mailList[k])
				delete(mailQueu, k)
				SendMail(mailList[k])
				time.Sleep(3 * time.Second)
			}
		}
	}
}

func SendMail(mail Mail) {
	m := gomail.NewMessage()
	m.SetHeader("From", mail.From)
	m.SetHeader("To", mail.To[0], mail.To[1], mail.To[2])
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Body)

	d := gomail.NewDialer(cfg.Server, cfg.Port, cfg.Username, cfg.Password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func PostMail(id int, mail Mail) {
	_, ok := mailQueu[id]
	if !ok {
		mailQueu[id] = time.Now().Unix()
		mailList[id] = mail
	}
}
