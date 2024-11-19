package archive

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleArchiveInformation(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleArchiveInformation: request received")

	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Parsing multipart form data")
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Failed to parse multipart form data: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file from request: %v", err)
		http.Error(w, "Failed to get file from request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Received file: %s", header.Filename)

	archiveInfo, err := GiveInfoArchive(file, header)
	if err != nil {
		log.Printf("Error processing archive %s: %v", header.Filename, err)
		http.Error(w, "Error processing archive: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Archive %s processed successfully", header.Filename)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(archiveInfo)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}

	log.Println("Response sent successfully")
}
func HandleCreateArchive(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleCreateArchive: request received")

	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Failed to parse form data: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		log.Println("No files provided")
		http.Error(w, "No files provided", http.StatusBadRequest)
		return
	}

	log.Printf("Number of files received: %d", len(files))

	archiveBytes, err := CreateArchive(files)
	if err != nil {
		log.Printf("Error creating archive: %v", err)
		http.Error(w, "Error creating archive: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=archive.zip")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(archiveBytes)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Error writing response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Archive created and response sent successfully")
}
