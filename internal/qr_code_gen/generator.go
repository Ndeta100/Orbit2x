package qr_code_gen

import (
	"encoding/base64"
	"fmt"
	"image/color"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// QRRequest represents user input for QR generation
type QRRequest struct {
	Type string `json:"type" form:"type"` // text, url, email, phone, sms
	Data string `json:"data" form:"data"` // user's content
}

// QRResponse represents the generated QR code
type QRResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	ImageData string `json:"image_data,omitempty"` // base64 data URL
	QRType    string `json:"qr_type,omitempty"`
	Content   string `json:"content,omitempty"`
}

// Generate creates a QR code from user input - main function
func Generate(req QRRequest) *QRResponse {
	// Validate input
	if strings.TrimSpace(req.Data) == "" {
		return &QRResponse{
			Success: false,
			Message: "Please enter some data to generate QR code",
		}
	}

	// Prepare data based on type
	processedData, err := prepareData(req)
	if err != nil {
		return &QRResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	// Generate QR code
	imageData, err := createQRCode(processedData)
	if err != nil {
		return &QRResponse{
			Success: false,
			Message: "Failed to generate QR code. Please try again.",
		}
	}

	return &QRResponse{
		Success:   true,
		Message:   "QR code generated successfully!",
		ImageData: imageData,
		QRType:    req.Type,
		Content:   req.Data,
	}
}

// prepareData formats user input based on QR type
func prepareData(req QRRequest) (string, error) {
	data := strings.TrimSpace(req.Data)

	switch req.Type {
	case "url":
		return prepareURL(data)
	case "email":
		return prepareEmail(data)
	case "phone":
		return preparePhone(data)
	case "sms":
		return prepareSMS(data)
	case "text":
		fallthrough
	default:
		if len(data) > 1000 {
			return "", fmt.Errorf("text is too long (max 1000 characters)")
		}
		return data, nil
	}
}

// prepareURL validates and formats URL - FIXED: Keep original URL format
func prepareURL(data string) (string, error) {
	if data == "" {
		return "", fmt.Errorf("please enter a website URL")
	}

	// Clean the input
	data = strings.TrimSpace(data)

	// For validation, add https:// if missing
	validationURL := data
	if !strings.HasPrefix(data, "http://") && !strings.HasPrefix(data, "https://") {
		validationURL = "https://" + data
	}

	// Validate URL format
	if _, err := url.Parse(validationURL); err != nil {
		return "", fmt.Errorf("please enter a valid URL (e.g., example.com or https://example.com)")
	}

	// FIXED: Return the original input or with https:// - both work in QR readers
	if strings.HasPrefix(data, "http://") || strings.HasPrefix(data, "https://") {
		return data, nil
	}

	// Add https:// for better compatibility
	return "https://" + data, nil
}

// prepareEmail formats email - FIXED: Provide both options
func prepareEmail(data string) (string, error) {
	if data == "" {
		return "", fmt.Errorf("please enter an email address")
	}

	// Basic email validation
	if !strings.Contains(data, "@") || !strings.Contains(data, ".") {
		return "", fmt.Errorf("please enter a valid email address")
	}

	// FIXED: Try without mailto: first for better compatibility
	// Most modern QR readers can detect email addresses without the mailto: prefix
	// But we'll keep mailto: for better app integration
	return fmt.Sprintf("mailto:%s", data), nil
}

// preparePhone formats phone number - FIXED: Provide cleaner format
func preparePhone(data string) (string, error) {
	if data == "" {
		return "", fmt.Errorf("please enter a phone number")
	}

	// Clean phone number
	cleaned := strings.ReplaceAll(data, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")

	if len(cleaned) < 5 {
		return "", fmt.Errorf("please enter a valid phone number")
	}

	// FIXED: Try without tel: prefix first for universal compatibility
	// Many QR readers work better with plain phone numbers
	return cleaned, nil
}

// prepareSMS formats SMS (expects "phone|message" format from frontend)
func prepareSMS(data string) (string, error) {
	if data == "" {
		return "", fmt.Errorf("please enter a phone number")
	}

	parts := strings.Split(data, "|")
	phone := strings.TrimSpace(parts[0])

	if phone == "" {
		return "", fmt.Errorf("please enter a phone number")
	}

	// Clean phone
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	if len(parts) > 1 && strings.TrimSpace(parts[1]) != "" {
		message := strings.TrimSpace(parts[1])
		return fmt.Sprintf("sms:%s?body=%s", phone, url.QueryEscape(message)), nil
	}

	return fmt.Sprintf("sms:%s", phone), nil
}

// createQRCode generates the actual QR code image
func createQRCode(data string) (string, error) {
	// Create QR code with higher error correction for better scanning
	qrc, err := qrcode.NewWith(data,
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart), // FIXED: Using ErrorCorrectionQuart (25% recovery)
	)
	if err != nil {
		return "", err
	}

	// Create temporary file
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("qr_temp_%d.png", len(data)))
	defer os.Remove(tmpFile)

	// Create writer with optimized settings for scanning
	writer, err := standard.New(tmpFile,
		standard.WithFgColor(color.RGBA{R: 0, G: 0, B: 0, A: 255}),       // Black
		standard.WithBgColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}), // White
		standard.WithQRWidth(10),    // IMPROVED: Slightly smaller for better balance
		standard.WithBorderWidth(2), // IMPROVED: Smaller border for more QR content space
	)
	if err != nil {
		return "", err
	}

	// Save QR code
	if err := qrc.Save(writer); err != nil {
		return "", err
	}

	// Read file and convert to base64
	imageBytes, err := os.ReadFile(tmpFile)
	if err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(imageBytes)
	return fmt.Sprintf("data:image/png;base64,%s", base64String), nil
}
