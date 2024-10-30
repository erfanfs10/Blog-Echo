package utils

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

func SendEmail(userEmail, subject, text string) error {
	// get myEmail and emailPassword from .env
	myEmail := os.Getenv("EMAIL")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	// create email
	e := email.NewEmail()
	from := fmt.Sprintf("Mahdi FarahmandSafa <%v>", myEmail)
	e.From = from
	e.To = []string{userEmail}
	e.Subject = subject
	e.Text = []byte(text)
	// send email
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth(
		"", myEmail,
		emailPassword, "smtp.gmail.com"))
	// Check error message form sending email
	if err != nil {
		return err
	}
	return nil
}
