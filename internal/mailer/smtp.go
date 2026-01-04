package mailer

import (
	"fmt"
	"gmail_microservice/internal/model"
	"net/smtp"
)

// Send menerima parameter 'account' untuk menentukan siapa pengirimnya
func Send(job model.EmailPayload, account model.GmailAccount) error {
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	// Header Email
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=\"utf-8\"\r\n"+
		"\r\n"+
		"%s", account.Email, job.To, job.Subject, job.Message))

	// Autentikasi menggunakan akun yang dipilih
	auth := smtp.PlainAuth("", account.Email, account.Password, host)

	// Eksekusi
	err := smtp.SendMail(address, auth, account.Email, []string{job.To}, msg)
	if err != nil {
		return err
	}

	return nil
}
