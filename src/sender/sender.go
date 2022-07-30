package sender

import (
	"btc_service/src/btc"
	"btc_service/src/model"
	"btc_service/src/persistance"
	"crypto/tls"
	"strconv"

	"gopkg.in/gomail.v2"
)

type (
	sender struct {
	}
	database interface {
		Save()
		Exists(model.Email) bool
	}
)

type MailSender struct {
	From      string
	To        []model.Email
	Subject   string
	Text      string
	Attach    string
	password  string
	smtp_name string
}

func New(from string, smtpName string, password string, text string, subject string) *MailSender {
	var mailSender MailSender
	mailSender.From = from
	mailSender.smtp_name = smtpName
	mailSender.password = password
	mailSender.Text = text
	mailSender.Subject = subject
	return &mailSender
}

func (mailSender *MailSender) send() error {
	m := gomail.NewMessage()

	m.SetHeader("From", mailSender.From)

	addresses := make([]string, len(mailSender.To))
	for i, recipient := range mailSender.To {
		addresses[i] = m.FormatAddress(string(recipient), "")
	}

	m.SetHeader("To", addresses...)

	m.SetHeader("Subject", mailSender.Subject)

	m.SetBody("text/plain", mailSender.Text)

	d := gomail.NewDialer(mailSender.smtp_name, 465, mailSender.From, mailSender.password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(m)
	return err
}

func (mailSender *MailSender) SendRate(fdb *persistance.FileDatabase, apiName string) (error, error) {
	btc_rate, rateErr := btc.GetRate(apiName)
	mailSender.To = fdb.Buffer
	mailSender.Text += " " + strconv.Itoa(btc_rate)
	sendErr := mailSender.send()
	return sendErr, rateErr
}
