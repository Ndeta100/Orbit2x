// handlers/imagebase64_handler.go
package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Ndeta100/orbit2x/views/imagebase64" // Adjust to your actual path
)

// HandleImageBase64Index renders the Image to Base64 converter page
func HandleImageBase64Index(w http.ResponseWriter, r *http.Request) error {
	return imagebase64.Index().Render(r.Context(), w)
}

// HandleImageBase64Convert converts an uploaded image to a Base64 data URL
func HandleImageBase64Convert(w http.ResponseWriter, r *http.Request) error {
	// Set max upload size - 5MB
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)
	if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		return imagebase64.Results(imagebase64.ConversionResult{
			Error: "The uploaded file is too large. Maximum size is 5MB.",
		}).Render(r.Context(), w)
	}

	// Get file from form
	file, header, err := r.FormFile("image")
	if err != nil {
		return imagebase64.Results(imagebase64.ConversionResult{
			Error: "Failed to get file from form: " + err.Error(),
		}).Render(r.Context(), w)
	}
	defer file.Close()

	// Check if the file is an image
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return imagebase64.Results(imagebase64.ConversionResult{
			Error: "Uploaded file is not an image. Please upload an image file.",
		}).Render(r.Context(), w)
	}

	// Read the image file
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return imagebase64.Results(imagebase64.ConversionResult{
			Error: "Failed to read image file: " + err.Error(),
		}).Render(r.Context(), w)
	}

	// Convert to base64
	base64Data := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Format as data URL
	// If the content type is not specified, try to detect it from the file extension
	if contentType == "" || contentType == "application/octet-stream" {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		case ".svg":
			contentType = "image/svg+xml"
		case ".webp":
			contentType = "image/webp"
		default:
			contentType = "image/png" // Default to PNG
		}
	}

	dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)

	// Generate HTML tag examples
	imgTag := fmt.Sprintf(`<img src="%s" alt="Base64 image">`, dataURL)
	cssBackground := fmt.Sprintf(`background-image: url('%s');`, dataURL)

	// Create the result
	result := imagebase64.ConversionResult{
		FileName:      header.Filename,
		FileSize:      formatFileSize(header.Size),
		ContentType:   contentType,
		DataURL:       dataURL,
		ImgTag:        imgTag,
		CSSBackground: cssBackground,
	}

	// Render the result
	return imagebase64.Results(result).Render(r.Context(), w)
}

// formatFileSize converts file size in bytes to a human-readable format
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
