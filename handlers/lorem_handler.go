package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ndeta100/orbit2x/internal/lorem"
	"github.com/Ndeta100/orbit2x/views/lorem_ipsum"
)

// HandleLoremMainPage HandleMainPage serves the main Lorem Ipsum generator page
func HandleLoremMainPage(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	setSecurityHeaders(w)

	data := lorem_ipsum.LoremPageData{
		Title: "Lorem Ipsum Generator - Professional Placeholder Text | Orbit2x",
	}

	// Render the page using your templ component
	component := lorem_ipsum.LoremGeneratorPage(data)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering Lorem Ipsum page: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil
	}
	return nil
}

// HandleLoremGenerate processes Lorem Ipsum generation requests (HTMX)
func HandleLoremGenerate(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		RenderLoremError(w, r, "Invalid form data")
		return nil
	}

	// Extract and validate form data
	req, err := ParseLoremRequest(r)
	if err != nil {
		log.Printf("Error parsing Lorem request: %v", err)
		RenderLoremError(w, r, err.Error())
		return nil
	}

	// Generate Lorem Ipsum text
	response := lorem.Generate(req)

	// Log for monitoring
	if response.Success {
		log.Printf("Generated %s: %d count, %d words, %d chars",
			req.Type, req.Count, response.WordCount, response.CharCount)
	} else {
		log.Printf("Generation failed: %s", response.Message)
	}

	// Render the result component
	component := lorem_ipsum.LoremResult(*response)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering Lorem result: %v", err)
		RenderLoremError(w, r, "Failed to render result")
		return nil
	}
	return nil
}

// HandleLoremAPIGenerate processes API requests and returns JSON
func HandleLoremAPIGenerate(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	// Set CORS headers for API usage
	setCORSHeaders(w)

	var req lorem.LoremRequest

	// Handle both JSON and form data
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			SendLoremErrorJSON(w, "Invalid JSON data", http.StatusBadRequest)
			return nil
		}
	} else {
		// Parse as form data
		if err := r.ParseForm(); err != nil {
			SendLoremErrorJSON(w, "Invalid form data", http.StatusBadRequest)
			return nil
		}

		var err error
		req, err = ParseLoremRequest(r)
		if err != nil {
			SendLoremErrorJSON(w, err.Error(), http.StatusBadRequest)
			return nil
		}
	}

	// Generate Lorem Ipsum
	response := lorem.Generate(req)

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		SendLoremErrorJSON(w, "Internal server error", http.StatusInternalServerError)
		return nil
	}
	return nil
}

// HandleLoremLimits returns the generation limits for frontend validation
func HandleLoremLimits(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	setCORSHeaders(w)

	limits := lorem.GetLimits()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(limits); err != nil {
		log.Printf("Error encoding limits JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil
	}
	return nil
}

// ParseLoremRequest extracts and validates Lorem Ipsum request from form data
func ParseLoremRequest(r *http.Request) (lorem.LoremRequest, error) {
	var req lorem.LoremRequest

	// Parse type
	req.Type = strings.TrimSpace(r.FormValue("type"))
	if req.Type == "" {
		req.Type = "paragraphs" // Default to paragraphs
	}

	// Validate type
	validTypes := map[string]bool{
		"words":      true,
		"sentences":  true,
		"paragraphs": true,
	}
	if !validTypes[req.Type] {
		return req, fmt.Errorf("invalid type: must be 'words', 'sentences', or 'paragraphs'")
	}

	// Parse count
	countStr := strings.TrimSpace(r.FormValue("count"))
	if countStr == "" {
		req.Count = 3 // Default count
	} else {
		var err error
		req.Count, err = strconv.Atoi(countStr)
		if err != nil {
			return req, fmt.Errorf("invalid count: must be a number")
		}
	}

	// Parse start_with option
	req.StartWith = r.FormValue("start_with") == "true"

	return req, nil
}

// RenderLoremError renders an error component
func RenderLoremError(w http.ResponseWriter, r *http.Request, message string) {
	component := lorem_ipsum.LoremError(message)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering error component: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// SendLoremErrorJSON sends a JSON error response
func SendLoremErrorJSON(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := lorem.LoremResponse{
		Success: false,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

// setSecurityHeaders sets common security headers
func setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
}

// setCORSHeaders sets CORS headers for API endpoints
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
