package handlers

import (
	"net/http"

	"github.com/Ndeta100/orbit2x/internal/word_counter"
	word_counter_view "github.com/Ndeta100/orbit2x/views/word_count"
)

// WordCountRequest represents the request payload for word counting
type WordCountRequest struct {
	Text string `json:"text"`
}

// WordCountResponse represents the response payload for word counting
type WordCountResponse struct {
	Words                int    `json:"words"`
	Sentences            int    `json:"sentences"`
	CharactersWithSpaces int    `json:"charactersWithSpaces"`
	CharactersNoSpaces   int    `json:"charactersNoSpaces"`
	Paragraphs           int    `json:"paragraphs"`
	ReadingTimeFormatted string `json:"readingTimeFormatted"`
}

// ShowWordCounterPage renders the word counter page
func ShowWordCounterPage(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	data := word_counter_view.WordCounterPageData{
		Title: "Free Word Counter - Count Words, Characters & Reading Time | Orbit2x",
	}

	// Render the word counter page
	err := word_counter_view.WordCounterPage(data).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil
	}
	return nil
}

// CountWords handles the word counting request and returns rendered HTML
func CountWords(w http.ResponseWriter, r *http.Request) error {
	// Set response content type for HTML
	w.Header().Set("Content-Type", "text/html")

	// Parse form data (HTMX sends form data by default)
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// Return error in HTML format for HTMX
		w.Write([]byte(`<div class="p-4 text-red-600">Invalid request format</div>`))
		return nil
	}

	// Get text from form field
	text := r.FormValue("text-input")

	// Validate input length (1MB limit for safety)
	if len(text) > 1000000 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`<div class="p-4 text-red-600">Text too long</div>`))
		return nil
	}

	// Get word count statistics
	stats := word_counter.CountAll(text)

	// Prepare stats for templ component
	wordStats := word_counter_view.WordCountStats{
		Words:                stats.Words,
		Sentences:            stats.Sentences,
		CharactersWithSpaces: stats.CharactersWithSpaces,
		CharactersNoSpaces:   stats.CharactersNoSpaces,
		Paragraphs:           stats.Paragraphs,
		ReadingTimeFormatted: stats.ReadingTimeFormatted,
	}

	// Render the statistics component
	w.WriteHeader(http.StatusOK)
	return word_counter_view.WordCounterStats(wordStats).Render(r.Context(), w)
}
