package mail

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func HandleSendEmail(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleSendEmail: request received")

	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Retrieving file from request")
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file from request: %v", err)
		http.Error(w, "Failed to get file from request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	log.Printf("File received: %s", header.Filename)

	log.Println("Reading file content")
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("File %s read successfully", header.Filename)

	log.Println("Retrieving email addresses from request")
	emails := r.FormValue("emails")
	if emails == "" {
		log.Println("No email addresses provided")
		http.Error(w, "No email addresses provided", http.StatusBadRequest)
		return
	}
	recipients := strings.Split(emails, ",")
	log.Printf("Recipients: %v", recipients)

	emailRequest := EmailRequest{
		Filename:   header.Filename,
		Recipients: recipients,
		FileBytes:  fileBytes,
	}

	log.Println("Sending email")
	err = SendEmail(emailRequest)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Email sent successfully")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Email sent successfully"))
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
