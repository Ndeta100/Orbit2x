// handlers/formatter_handler.go
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Ndeta100/orbit2x/views/formatter" // Adjust to your actual path
	"gopkg.in/yaml.v3"                            // You'll need to add this to your dependencies
)

// HandleFormatterIndex renders the JSON & YAML formatter page
func HandleFormatterIndex(w http.ResponseWriter, r *http.Request) error {
	return formatter.Index().Render(r.Context(), w)
}

// HandleJSONFormat formats and validates JSON
func HandleJSONFormat(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return formatter.Results(formatter.FormatterResult{
			Error:  "Failed to parse form data",
			Format: "JSON",
		}).Render(r.Context(), w)
	}

	// Get text and indentation from form
	text := r.FormValue("text")
	if text == "" {
		return formatter.Results(formatter.FormatterResult{
			Error:  "JSON text is required",
			Format: "JSON",
		}).Render(r.Context(), w)
	}

	indentStr := r.FormValue("indent")
	indent := 4 // Default indentation
	if indentStr != "" {
		var err error
		indent, err = strconv.Atoi(indentStr)
		if err != nil {
			indent = 4 // Default to 4 if there's an error
		}
	}

	// Parse JSON to validate it
	var data interface{}
	err := json.Unmarshal([]byte(text), &data)
	if err != nil {
		return formatter.Results(formatter.FormatterResult{
			Error:        "Invalid JSON: " + err.Error(),
			OriginalText: text,
			Format:       "JSON",
		}).Render(r.Context(), w)
	}

	// Format JSON with proper indentation
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", generateIndent(indent))
	encoder.SetEscapeHTML(false) // Avoid escaping HTML characters
	if err := encoder.Encode(data); err != nil {
		return formatter.Results(formatter.FormatterResult{
			Error:        "Failed to format JSON: " + err.Error(),
			OriginalText: text,
			Format:       "JSON",
		}).Render(r.Context(), w)
	}

	// Create result
	result := formatter.FormatterResult{
		OriginalText:  text,
		FormattedText: buf.String(),
		Format:        "JSON",
	}

	// Render the result
	return formatter.Results(result).Render(r.Context(), w)
}

// HandleYAMLFormat formats and validates YAML
func HandleYAMLFormat(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return formatter.Results(formatter.FormatterResult{
			Error:  "Failed to parse form data",
			Format: "YAML",
		}).Render(r.Context(), w)
	}

	// Get text and indentation from form
	text := r.FormValue("text")
	if text == "" {
		return formatter.Results(formatter.FormatterResult{
			Error:  "YAML text is required",
			Format: "YAML",
		}).Render(r.Context(), w)
	}

	indentStr := r.FormValue("indent")
	indent := 2 // Default indentation for YAML
	if indentStr != "" {
		var err error
		indent, err = strconv.Atoi(indentStr)
		if err != nil {
			indent = 2 // Default to 2 if there's an error
		}
	}

	// Parse YAML to validate it
	var data interface{}
	err := yaml.Unmarshal([]byte(text), &data)
	if err != nil {
		return formatter.Results(formatter.FormatterResult{
			Error:        "Invalid YAML: " + err.Error(),
			OriginalText: text,
			Format:       "YAML",
		}).Render(r.Context(), w)
	}

	// Format YAML with proper indentation
	yamlEncoder := yaml.NewEncoder(&bytes.Buffer{})
	yamlEncoder.SetIndent(indent)
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(indent)
	if err := encoder.Encode(data); err != nil {
		return formatter.Results(formatter.FormatterResult{
			Error:        "Failed to format YAML: " + err.Error(),
			OriginalText: text,
			Format:       "YAML",
		}).Render(r.Context(), w)
	}

	// Create result
	result := formatter.FormatterResult{
		OriginalText:  text,
		FormattedText: buf.String(),
		Format:        "YAML",
	}

	// Render the result
	return formatter.Results(result).Render(r.Context(), w)
}

// generateIndent creates a string with the specified number of spaces
func generateIndent(count int) string {
	indent := ""
	for i := 0; i < count; i++ {
		indent += " "
	}
	return indent
}
