# Test project for Internship in Doodocs.kz

## Description
This test project includes REST API and 3 routes for working with `.zip` files

## Installation
```bash 
# Clone repository
git clone https://github.com/username/repository.git

# Move to project root directory
cd dodocs-test-app

# Download all required dependencies
go mod tidy
```

## Usage

```bash
# Launch project
go run cmd/main.go
```

## Configuration
In root directory create .env file. Example of .env file:
```
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your-email@example.com
SMTP_PASSWORD=your-password
```

## Features
Project includes 3 routes(endpoints) which works with .zip file:

#### POST `/api/archive/information` 
- route for getting information about `.zip `file. `.zip` file sent by `multipart/form-data` in field `file`

- Example of HTTP-request:
```
POST /api/archive/information HTTP/1.1

Content-Type: multipart/form-data; boundary=-{some-random-boundary}

-{some-random-boundary}
Content-Disposition: form-data; name="file"; filename="my_archive.zip"
Content-Type: application/zip

{Binary data of ZIP file}
-{some-random-boundary}-
```
- Example of HTTP-response:
```
HTTP/1.1 200 OK
Content-Type: application/json

{
    "filename": "my_archive.zip",
    "archive_size": 4102029.312,
    "total_size": 6836715.52,
    "total_files": 2,
    "files": [
        {
            "file_path": "photo.jpg",
            "size": 2516582.4,
            "mimetype": "image/jpeg"
        },
        {
            "file_path": "directory/document.docx",
            "size": 4320133.12,
            "mimetype": "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        }
    ]
}
```

#### POST `/api/archive/files`
- route for creating `.zip` file. Files sent by `multipart/form-data` in field `files[]` 
- Allowed file types:
  - `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
  - `application/xml`
  - `image/jpeg`
  - `image/png`
- HTTP-request:
```
POST /api/archive/files HTTP/1.1
  Content-Type: multipart/form-data; boundary=-{some-random-boundary}

-{some-random-boundary}
Content-Disposition: form-data; name="files[]"; filename="document.docx"
Content-Type: application/vnd.openxmlformats-officedocument.wordprocessingml.document

{Binary data of file}
-{some-random-boundary}
Content-Disposition: form-data; name="files[]"; filename="avatar.png"
Content-Type: image/png

{Binary data of file}
-{some-random-boundary}--
```

- HTTP-response:
```
HTTP/1.1 200 OK
Content-Type: application/zip

{Binary data of ZIP file}
```

#### POST `/api/mail/file`
- route for sending to mail allowed file. File sent by `multipart/form-data` in field `file`. List of recipients sent in `emails` field
- Allowed file types:
    - `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
    - `image/png`
- HTTP-request:
```
POST /api/mail/file HTTP/1.1
Content-Type: multipart/form-data; boundary=-{some-random-boundary}

-{some-random-boundary}
Content-Disposition: form-data; name="file"; filename="document.docx"
Content-Type: application/vnd.openxmlformats-officedocument.wordprocessingml.document

{Binary data of file}
-{some-random-boundary}
Content-Disposition: form-data; name="emails"

elonmusk@x.com,jeffbezos@amazon.com,zuckerberg@meta.com
-{some-random-boundary}--
```

- HTTP-response:
```
HTTP/1.1 200 OK
```

## Technologies Used
- [Go](https://golang.org/) — The main programming language for the project.
- [github.com/gabriel-vasile/mimetype](https://github.com/gabriel-vasile/mimetype) v1.4.6 — A library for detecting the MIME type of files.
- [github.com/joho/godotenv](https://github.com/joho/godotenv) v1.5.1 — A library to load environment variables from a `.env` file.
- [golang.org/x/net](https://pkg.go.dev/golang.org/x/net) v0.30.0 — Provides extended networking utilities for handling network protocols and connections.

## Contact
- [Telegram](https://t.me/boranggaziyev)
- [LinkedIn](https://www.linkedin.com/in/%D0%B0%D0%BB%D0%B8%D1%88%D0%B5%D1%80-%D0%B1%D0%BE%D1%80%D0%B0%D0%BD%D0%B3%D0%B0%D0%B7%D0%B8%D0%B5%D0%B2-a6a119298/)
- [Github](https://github.com/BorangaziyevA)