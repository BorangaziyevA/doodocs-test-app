package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"testProject/internal/archive"
	"testProject/internal/mail"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, make sure environment variables are set")
	}

	http.HandleFunc("/api/archive/information", archive.HandleArchiveInformation)
	http.HandleFunc("/api/archive/files", archive.HandleCreateArchive)
	http.HandleFunc("/api/mail/file", mail.HandleSendEmail)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
