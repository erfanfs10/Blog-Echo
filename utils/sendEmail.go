package utils

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

var EmailChannel chan EmailJob

type EmailJob struct {
	To      string
	Subject string
	Text    string
}

func CreateEmailChannel() {
	EmailChannel = make(chan EmailJob)
	fmt.Println("Channel created")
	go func() {
		for job := range EmailChannel {
			sendEmail(job)
		}
	}()
	fmt.Println("Goroutine started")
}

func sendEmail(job EmailJob) {
	// get myEmail and emailPassword from .env
	myEmail := GetEnv("EMAIL", "example@gmail.com")
	emailPassword := GetEnv("EMAIL_PASSWORD", "somePassword")
	// create email
	e := email.NewEmail()
	from := fmt.Sprintf("My Echo Program <%v>", myEmail)
	e.From = from
	e.To = []string{job.To}
	e.Subject = job.Subject
	e.Text = []byte(job.Text)
	// send email
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth(
		"", myEmail,
		emailPassword, "smtp.gmail.com"))
	// Check error message from sending email
	fmt.Println("---------------------------------------------------------")
	if err != nil {
		log.Printf("Failed to send email to %s: %v", job.To, err)
	}
	log.Printf("Email sent successfully to %s", job.To)
}
