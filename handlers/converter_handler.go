// handlers/converter_handler.go
package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ndeta100/orbit2x/views/converter" // Adjust to your actual path
)

// HandleConverterIndex renders the CSV & JSON converter page
func HandleConverterIndex(w http.ResponseWriter, r *http.Request) error {
	return converter.Index().Render(r.Context(), w)
}

// HandleCSVToJSON converts CSV data to JSON
func HandleCSVToJSON(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return converter.Results(converter.ConversionResult{
			Error:        "Failed to parse form data",
			SourceFormat: "CSV",
			TargetFormat: "JSON",
		}).Render(r.Context(), w)
	}

	// Get text and options from form
	text := r.FormValue("text")
	if text == "" {
		return converter.Results(converter.ConversionResult{
			Error:        "CSV text is required",
			SourceFormat: "CSV",
			TargetFormat: "JSON",
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

	hasHeaders := r.FormValue("_headers") != ""

	// Parse CSV
	reader := csv.NewReader(strings.NewReader(text))
	records, err := reader.ReadAll()
	if err != nil {
		return converter.Results(converter.ConversionResult{
			Error:        "Invalid CSV format: " + err.Error(),
			OriginalText: text,
			SourceFormat: "CSV",
			TargetFormat: "JSON",
		}).Render(r.Context(), w)
	}

	if len(records) == 0 {
		return converter.Results(converter.ConversionResult{
			Error:        "CSV has no data",
			OriginalText: text,
			SourceFormat: "CSV",
			TargetFormat: "JSON",
		}).Render(r.Context(), w)
	}

	// Convert to JSON
	var result []interface{}

	if hasHeaders && len(records) > 1 {
		// First row contains _headers
		headers := records[0]

		// Process data rows
		for i := 1; i < len(records); i++ {
			row := records[i]
			item := make(map[string]interface{})

			for j := 0; j < len(headers) && j < len(row); j++ {
				// Try to convert to number or boolean if possible
				value := row[j]

				// Check if it's a number
				if numVal, err := strconv.ParseFloat(value, 64); err == nil {
					item[headers[j]] = numVal
				} else if value == "true" || value == "false" {
					// Check if it's a boolean
					item[headers[j]] = value == "true"
				} else {
					// It's a string
					item[headers[j]] = value
				}
			}

			result = append(result, item)
		}
	} else {
		// No _headers, convert each row to an array
		for _, row := range records {
			rowData := make([]interface{}, len(row))

			for i, value := range row {
				// Try to convert to number or boolean if possible
				if numVal, err := strconv.ParseFloat(value, 64); err == nil {
					rowData[i] = numVal
				} else if value == "true" || value == "false" {
					rowData[i] = value == "true"
				} else {
					rowData[i] = value
				}
			}

			result = append(result, rowData)
		}
	}

	// Format JSON with proper indentation
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", generateIndent(indent))
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(result); err != nil {
		return converter.Results(converter.ConversionResult{
			Error:        "Failed to format JSON: " + err.Error(),
			OriginalText: text,
			SourceFormat: "CSV",
			TargetFormat: "JSON",
		}).Render(r.Context(), w)
	}

	// Create result
	convResult := converter.ConversionResult{
		OriginalText:  text,
		ConvertedText: buf.String(),
		SourceFormat:  "CSV",
		TargetFormat:  "JSON",
	}

	// Render the result
	return converter.Results(convResult).Render(r.Context(), w)
}

// HandleJSONToCSV converts JSON data to CSV
func HandleJSONToCSV(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return converter.Results(converter.ConversionResult{
			Error:        "Failed to parse form data",
			SourceFormat: "JSON",
			TargetFormat: "CSV",
		}).Render(r.Context(), w)
	}

	// Get text and options from form
	text := r.FormValue("text")
	if text == "" {
		return converter.Results(converter.ConversionResult{
			Error:        "JSON text is required",
			SourceFormat: "JSON",
			TargetFormat: "CSV",
		}).Render(r.Context(), w)
	}

	delimiter := r.FormValue("delimiter")
	if delimiter == "" {
		delimiter = "," // Default delimiter
	} else if delimiter == "\\t" {
		delimiter = "\t" // Handle tab character
	}

	includeHeaders := r.FormValue("includeHeaders") != ""

	// Parse JSON
	var jsonData []interface{}
	err := json.Unmarshal([]byte(text), &jsonData)

	if err != nil {
		// Try parsing as a single object
		var jsonObj map[string]interface{}
		err = json.Unmarshal([]byte(text), &jsonObj)

		if err != nil {
			return converter.Results(converter.ConversionResult{
				Error:        "Invalid JSON format: " + err.Error(),
				OriginalText: text,
				SourceFormat: "JSON",
				TargetFormat: "CSV",
			}).Render(r.Context(), w)
		}

		// Convert single object to array
		jsonData = []interface{}{jsonObj}
	}

	if len(jsonData) == 0 {
		return converter.Results(converter.ConversionResult{
			Error:        "JSON has no data",
			OriginalText: text,
			SourceFormat: "JSON",
			TargetFormat: "CSV",
		}).Render(r.Context(), w)
	}

	// Convert to CSV
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = []rune(delimiter)[0]

	// For array of objects, extract _headers
	//if firstItem, ok := jsonData[0].(map[string]interface{}); ok {
	if _, ok := jsonData[0].(map[string]interface{}); ok {
		// Array of objects

		// Get all possible _headers
		headers := make(map[string]bool)
		for _, item := range jsonData {
			if obj, ok := item.(map[string]interface{}); ok {
				for key := range obj {
					headers[key] = true
				}
			}
		}

		// Convert _headers map to slice
		headerSlice := make([]string, 0, len(headers))
		for key := range headers {
			headerSlice = append(headerSlice, key)
		}

		// Write _headers
		if includeHeaders {
			if err := writer.Write(headerSlice); err != nil {
				return converter.Results(converter.ConversionResult{
					Error:        "Failed to write CSV _headers: " + err.Error(),
					OriginalText: text,
					SourceFormat: "JSON",
					TargetFormat: "CSV",
				}).Render(r.Context(), w)
			}
		}

		// Write data rows
		for _, item := range jsonData {
			if obj, ok := item.(map[string]interface{}); ok {
				row := make([]string, len(headerSlice))

				for i, key := range headerSlice {
					if val, exists := obj[key]; exists {
						row[i] = formatJSONValue(val)
					}
				}

				if err := writer.Write(row); err != nil {
					return converter.Results(converter.ConversionResult{
						Error:        "Failed to write CSV row: " + err.Error(),
						OriginalText: text,
						SourceFormat: "JSON",
						TargetFormat: "CSV",
					}).Render(r.Context(), w)
				}
			}
		}
	} else {
		// Array of arrays
		for _, item := range jsonData {
			if arr, ok := item.([]interface{}); ok {
				row := make([]string, len(arr))

				for i, val := range arr {
					row[i] = formatJSONValue(val)
				}

				if err := writer.Write(row); err != nil {
					return converter.Results(converter.ConversionResult{
						Error:        "Failed to write CSV row: " + err.Error(),
						OriginalText: text,
						SourceFormat: "JSON",
						TargetFormat: "CSV",
					}).Render(r.Context(), w)
				}
			} else {
				// Handle simple values
				if err := writer.Write([]string{formatJSONValue(item)}); err != nil {
					return converter.Results(converter.ConversionResult{
						Error:        "Failed to write CSV row: " + err.Error(),
						OriginalText: text,
						SourceFormat: "JSON",
						TargetFormat: "CSV",
					}).Render(r.Context(), w)
				}
			}
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return converter.Results(converter.ConversionResult{
			Error:        "Failed to write CSV: " + err.Error(),
			OriginalText: text,
			SourceFormat: "JSON",
			TargetFormat: "CSV",
		}).Render(r.Context(), w)
	}

	// Create result
	convResult := converter.ConversionResult{
		OriginalText:  text,
		ConvertedText: buf.String(),
		SourceFormat:  "JSON",
		TargetFormat:  "CSV",
	}

	// Render the result
	return converter.Results(convResult).Render(r.Context(), w)
}

// formatJSONValue converts a JSON value to a string
func formatJSONValue(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return ""
	case bool:
		return strconv.FormatBool(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	case []interface{}:
		// Convert array to JSON string
		bytes, _ := json.Marshal(v)
		return string(bytes)
	case map[string]interface{}:
		// Convert object to JSON string
		bytes, _ := json.Marshal(v)
		return string(bytes)
	default:
		// Convert any other type to string
		return fmt.Sprintf("%v", v)
	}
}
