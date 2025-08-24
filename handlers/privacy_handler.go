package handlers

import (
	"net/http"

	"github.com/Ndeta100/orbit2x/views/privacy"
)

// PrivacyHandler handles the privacy policy page
func PrivacyHandler(w http.ResponseWriter, r *http.Request) error {
	// Set security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	// Render the privacy page
	component := privacy.Privacy()
	return component.Render(r.Context(), w)
}
