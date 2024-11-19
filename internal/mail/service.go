package mail

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
)

func SendEmail(emailRequest EmailRequest) error {
	log.Println("SendEmail: starting email preparation")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPassword == "" {
		log.Println("SMTP configuration is missing or incomplete")
		return errors.New("SMTP configuration is missing or incomplete")
	}

	boundary := fmt.Sprintf("BOUNDARY_%d", time.Now().UnixNano())
	headers := make(map[string]string)
	headers["From"] = smtpUser
	headers["To"] = strings.Join(emailRequest.Recipients, ", ")
	headers["Subject"] = "File Delivery"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "multipart/mixed; boundary=" + boundary

	var body bytes.Buffer
	for key, value := range headers {
		body.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	body.WriteString("\r\n")

	log.Println("Adding email body content")
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	body.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	body.WriteString("Please find the attached file.\r\n\r\n")

	log.Printf("Adding attachment: %s", emailRequest.Filename)
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Type: application/octet-stream\r\n")
	body.WriteString("Content-Transfer-Encoding: base64\r\n")
	body.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", emailRequest.Filename))

	encoded := base64.StdEncoding.EncodeToString(emailRequest.FileBytes)
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		body.WriteString(encoded[i:end] + "\r\n")
	}

	body.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	log.Println("Sending email")
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		smtpUser,
		emailRequest.Recipients,
		body.Bytes(),
	)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
