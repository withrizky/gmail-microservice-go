package model

// Struktur Payload dari User
type EmailPayload struct {
	To      string `json:"to" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// Struktur Konfigurasi Akun Gmail
type GmailAccount struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
