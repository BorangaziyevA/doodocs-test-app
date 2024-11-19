package tests

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testProject/internal/archive"
	"testing"
)

func TestHandleArchiveInformation(t *testing.T) {
	file, err := os.Open("testFiles/archive.zip")
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
	part, err := writer.CreateFormFile("file", "archive.zip")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/archive/information", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(archive.HandleArchiveInformation)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Body.Len() == 0 {
		t.Errorf("handler returned empty body")
	}

	expectedSubstring := `"filename":"archive.zip"`
	if !bytes.Contains(rr.Body.Bytes(), []byte(expectedSubstring)) {
		t.Errorf("handler response body does not contain expected data: %v", expectedSubstring)
	}
}
func TestHandleCreateArchive(t *testing.T) {
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

	part, err := writer.CreateFormFile("files[]", "testPng.png")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/archive/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Content-Length", fmt.Sprintf("%d", body.Len()))
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "GoTestClient/1.0")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(archive.HandleCreateArchive)
	handler.ServeHTTP(rr, req)
	fmt.Println(req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentDisposition := rr.Header().Get("Content-Disposition")
	if contentDisposition != "attachment; filename=archive.zip" {
		t.Errorf("handler returned unexpected Content-Disposition: got %v want %v", contentDisposition, "attachment; filename=archive.zip")
	}
}
