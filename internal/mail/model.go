package mail

// EmailRequest представляет структуру запроса для отправки файла по электронной почте.
type EmailRequest struct {
	Filename   string   `json:"filename"`   // Имя файла
	Recipients []string `json:"recipients"` // Список получателей
	FileBytes  []byte   `json:"-"`          // Содержимое файла (не экспортируется в JSON)
}
