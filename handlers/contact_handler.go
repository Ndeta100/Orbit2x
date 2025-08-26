package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Ndeta100/orbit2x/views/contact"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type ContactSubmission struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

// ContactHandler handles GET requests to /contact
func ContactHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	component := contact.ContactPage()
	return component.Render(r.Context(), w)
}

//// ContactSubmitHandler handles POST requests to /contact/submit
//func ContactSubmitHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	// Handle form submission logic here
//	// For now, just return success
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(map[string]string{
//		"status":  "success",
//		"message": "Message sent successfully",
//	})
//}

// ContactSubmitHandler handles form submissions
func ContactSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	submission := ContactSubmission{
		Name:    strings.TrimSpace(r.FormValue("name")),
		Email:   strings.TrimSpace(r.FormValue("email")),
		Subject: r.FormValue("subject"),
		Message: strings.TrimSpace(r.FormValue("message")),
		Date:    time.Now().Format("2006-01-02 15:04:05"),
	}

	// Basic validation
	if submission.Name == "" || submission.Email == "" ||
		submission.Subject == "" || submission.Message == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Log submission (replace with database storage in production)
	log.Printf("Contact form submission: %+v", submission)

	// Send email notification (optional)
	if err := sendEmailNotification(submission); err != nil {
		log.Printf("Error sending email notification: %v", err)
		// Don't fail the request if email fails
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Message sent successfully",
	})
}

// Send email notification (configure SMTP settings)
func sendEmailNotification(submission ContactSubmission) error {
	// Configure these environment variables:
	smtpHost := os.Getenv("SMTP_HOST")    // e.g., "smtp.gmail.com"
	smtpPort := os.Getenv("SMTP_PORT")    // e.g., "587"
	smtpUser := os.Getenv("SMTP_USER")    // your email
	smtpPass := os.Getenv("SMTP_PASS")    // your password or app password
	toEmail := os.Getenv("CONTACT_EMAIL") // where to send notifications

	if smtpHost == "" || smtpUser == "" || toEmail == "" {
		return fmt.Errorf("SMTP configuration missing")
	}

	// Email content
	subject := fmt.Sprintf("Contact Form: %s", submission.Subject)
	body := fmt.Sprintf(`
New contact form submission from Orbit2x:

Name: %s
Email: %s
Subject: %s
Date: %s

Message:
%s

---
This is an automated notification from Orbit2x contact form.
`, submission.Name, submission.Email, submission.Subject, submission.Date, submission.Message)

	// Send email
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", toEmail, subject, body))

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{toEmail}, msg)
}
