// handlers/encoder_handler.go
package handlers

import (
	"encoding/base64"
	"net/http"

	"github.com/Ndeta100/orbit2x/views/encoder" // Adjust to your actual path
)

// HandleEncoderIndex renders the text encoder/decoder page
func HandleEncoderIndex(w http.ResponseWriter, r *http.Request) error {
	return encoder.Index().Render(r.Context(), w)
}

// HandleEncoderEncode encodes text to Base64
func HandleEncoderEncode(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return encoder.Results(encoder.EncodingResult{
			Error: "Failed to parse form data",
			Mode:  "encode",
		}).Render(r.Context(), w)
	}

	// Get text from form
	text := r.FormValue("text")
	if text == "" {
		return encoder.Results(encoder.EncodingResult{
			Error: "Text is required",
			Mode:  "encode",
		}).Render(r.Context(), w)
	}

	// Encode to Base64
	encoded := base64.StdEncoding.EncodeToString([]byte(text))

	// Create result
	result := encoder.EncodingResult{
		OriginalText: text,
		EncodedText:  encoded,
		Mode:         "encode",
	}

	// Render the result
	return encoder.Results(result).Render(r.Context(), w)
}

// HandleEncoderDecode decodes Base64 to text
func HandleEncoderDecode(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return encoder.Results(encoder.EncodingResult{
			Error: "Failed to parse form data",
			Mode:  "decode",
		}).Render(r.Context(), w)
	}

	// Get text from form
	text := r.FormValue("text")
	if text == "" {
		return encoder.Results(encoder.EncodingResult{
			Error: "Base64 text is required",
			Mode:  "decode",
		}).Render(r.Context(), w)
	}

	// Decode from Base64
	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return encoder.Results(encoder.EncodingResult{
			Error:        "Invalid Base64 format: " + err.Error(),
			OriginalText: text,
			Mode:         "decode",
		}).Render(r.Context(), w)
	}

	// Create result
	result := encoder.EncodingResult{
		OriginalText: text,
		DecodedText:  string(decoded),
		Mode:         "decode",
	}

	// Render the result
	return encoder.Results(result).Render(r.Context(), w)
}
