package handlers

import (
	"github.com/Ndeta100/orbit2x/views/tools"
	"net/http"
)

func ToolsHandler(w http.ResponseWriter, r *http.Request) error {
	// Set SEO and performance headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// Render the tools page
	component := tools.ToolsPage()
	return component.Render(r.Context(), w)
}
