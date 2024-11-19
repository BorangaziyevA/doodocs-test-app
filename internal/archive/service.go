package archive

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"log"
	"mime/multipart"
)

func GiveInfoArchive(file multipart.File, header *multipart.FileHeader) (*ArchiveInfo, error) {
	log.Println("Started GiveInfoArchive function")
	defer log.Println("Ended GiveInfoArchive function")

	mime, err := mimetype.DetectReader(file)
	if err != nil {
		log.Printf("Cannot detect file type for %s: %v", header.Filename, err)
		return nil, errors.New("cannot detect file")
	}

	log.Printf("Detected MIME type for %s: %s", header.Filename, mime.String())

	if mime.String() != "application/zip" {
		log.Printf("File %s is not a valid ZIP archive", header.Filename)
		return nil, errors.New("file is not a valid ZIP archive")
	}

	fileInfo, err := header.Open()
	if err != nil {
		log.Printf("Failed to open file header for %s: %v", header.Filename, err)
		return nil, errors.New("failed to open file header")
	}
	defer fileInfo.Close()

	archive, err := zip.NewReader(fileInfo, header.Size)
	if err != nil {
		log.Printf("Failed to create ZIP reader for %s: %v", header.Filename, err)
		return nil, errors.New("failed to create ZIP reader")
	}

	var files []FileInfo
	totalSize := float64(0)
	for _, f := range archive.File {
		log.Printf("Processing file in archive: %s", f.Name)
		fileSize := float64(f.UncompressedSize64)
		totalSize += fileSize

		mimeType := detectMimeType(f)
		log.Printf("Detected MIME type for %s: %s", f.Name, mimeType)

		files = append(files, FileInfo{
			FilePath: f.Name,
			Size:     fileSize,
			MimeType: mimeType,
		})
	}

	archiveInfo := &ArchiveInfo{
		Filename:    header.Filename,
		ArchiveSize: float64(header.Size),
		TotalSize:   totalSize,
		TotalFiles:  len(files),
		Files:       files,
	}

	return archiveInfo, nil
}

func CreateArchive(files []*multipart.FileHeader) ([]byte, error) {
	log.Println("Started CreateArchive function")
	defer log.Println("Ended CreateArchive function")

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	allowedMIMETypes := map[string]bool{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/xml": true,
		"image/jpeg":      true,
		"image/png":       true,
	}

	for _, fileHeader := range files {
		log.Printf("Processing file: %s", fileHeader.Filename)
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Failed to open file %s: %v", fileHeader.Filename, err)
			return nil, errors.New("failed to open file: " + fileHeader.Filename)
		}
		defer file.Close()

		mimeType, err := mimetype.DetectReader(file)
		if err != nil || !allowedMIMETypes[mimeType.String()] {
			log.Printf("Invalid file type: %s (detected type: %s)", fileHeader.Filename, mimeType.String())
			return nil, errors.New("invalid file type: " + fileHeader.Filename)
		}

		if _, err = file.Seek(0, io.SeekStart); err != nil {
			log.Printf("Failed to reset file pointer for: %s", fileHeader.Filename)
			return nil, errors.New("failed to reset file pointer for: " + fileHeader.Filename)
		}

		writer, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			log.Printf("Failed to create entry in archive for: %s", fileHeader.Filename)
			return nil, errors.New("failed to create entry in archive for: " + fileHeader.Filename)
		}

		if _, err = io.Copy(writer, file); err != nil {
			log.Printf("Failed to write file to archive: %s", fileHeader.Filename)
			return nil, errors.New("failed to write file to archive: " + fileHeader.Filename)
		}

		log.Printf("File %s successfully added to archive", fileHeader.Filename)
	}

	if err := zipWriter.Close(); err != nil {
		log.Printf("Failed to close ZIP writer")
		return nil, errors.New("failed to close ZIP writer")
	}

	log.Println("Archive created successfully")
	return buf.Bytes(), nil
}

func detectMimeType(file *zip.File) string {
	log.Printf("Detecting MIME type for file: %s", file.Name)
	rc, err := file.Open()
	if err != nil {
		log.Printf("Failed to open file %s: %v", file.Name, err)
		return "unknown"
	}
	defer rc.Close()

	mime, err := mimetype.DetectReader(rc)
	if err != nil {
		log.Printf("Failed to detect MIME type for file %s: %v", file.Name, err)
		return "unknown"
	}

	log.Printf("Detected MIME type for file %s: %s", file.Name, mime.String())
	return mime.String()
}
