package tests

import (
	"bytes"
	"github.com/joho/godotenv"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testProject/internal/mail"
	"testing"
)

func TestHandleSendEmail(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	file, err := os.Open("testFiles/testPng.png")
	if err != nil {
		t.Fatalf("Failed to open local file: %v", err)
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read file content: %v", err)
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "testPng.png")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}

	err = writer.WriteField("emails", "alish3383@gmail.com")
	if err != nil {
		t.Fatalf("Failed to write form field: %v", err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/mail/file", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(mail.HandleSendEmail)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Email sent successfully"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
