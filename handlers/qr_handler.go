package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ndeta100/orbit2x/internal/qr_code_gen"
	"github.com/Ndeta100/orbit2x/views/qr_code"
)

// ShowQRTool displays the main QR generator page
func ShowQRTool(w http.ResponseWriter, r *http.Request) error {
	data := qr_code.QRPageData{
		Title: "Free QR Code Generator - orbit2x.com",
	}

	component := qr_code.QRGeneratorPage(data)
	w.Header().Set("Content-Type", "text/html")
	return component.Render(r.Context(), w)
}

// GenerateQR handles full QR generation (HTMX)
func GenerateQR(w http.ResponseWriter, r *http.Request) error {
	var req qr_code_gen.QRRequest

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			component := qr_code.QRError("Invalid request data")
			return component.Render(r.Context(), w)
		}
	} else {
		if err := r.ParseForm(); err != nil {
			component := qr_code.QRError("Failed to parse form data")
			return component.Render(r.Context(), w)
		}
		req = qr_code_gen.QRRequest{
			Type: r.FormValue("type"),
			Data: r.FormValue("data"),
		}
	}

	// Generate QR code using the package function
	response := qr_code_gen.Generate(req)

	templateData := qr_code.QRResultData{
		Success:   response.Success,
		Message:   response.Message,
		ImageData: response.ImageData,
		QRType:    response.QRType,
		Content:   response.Content,
	}

	component := qr_code.QRResult(templateData)
	return component.Render(r.Context(), w)
}

// PreviewQR handles live preview generation (HTMX)
func PreviewQR(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		component := qr_code.QRInitialState()
		return component.Render(r.Context(), w)
	}

	req := qr_code_gen.QRRequest{
		Type: r.FormValue("type"),
		Data: r.FormValue("data"),
	}

	if strings.TrimSpace(req.Data) == "" {
		component := qr_code.QRInitialState()
		return component.Render(r.Context(), w)
	}

	// Generate QR code using the package function
	response := qr_code_gen.Generate(req)

	if response.Success {
		component := qr_code.QRPreview(response.ImageData)
		return component.Render(r.Context(), w)
	}

	component := qr_code.QRInitialState()
	return component.Render(r.Context(), w)
}
