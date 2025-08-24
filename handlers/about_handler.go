package handlers

import (
	"github.com/Ndeta100/orbit2x/views/about"
	"net/http"
)

// AboutHandler handles the about page
func AboutHandler(w http.ResponseWriter, r *http.Request) error {
	// Set security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	// Render the about page
	component := about.About()
	return component.Render(r.Context(), w)
}
