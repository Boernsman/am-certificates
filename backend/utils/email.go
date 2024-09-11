package utils

import (
	"gopkg.in/gomail.v2"
	"log"
)

// SendEmail sends the certificate to the user
func SendEmail(email string, attachmentPath string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "youremail@domain.com")
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Your Certificate")
	mailer.SetBody("text/plain", "Here is your certificate!")
	mailer.Attach(attachmentPath)

	d := gomail.NewDialer("smtp.gmail.com", 587, "youremail@domain.com", "yourpassword")

	err := d.DialAndSend(mailer)
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}
	return nil
}
