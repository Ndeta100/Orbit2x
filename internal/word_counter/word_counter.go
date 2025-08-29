package word_counter

import (
	"regexp"
	"strings"
	"time"
	"unicode"
)

// Default words per minute for reading time calculation
const DefaultWordsPerMinute = 200

// WordCounterStats represents all the statistics for the word counter
type WordCounterStats struct {
	Words                int           `json:"words"`
	Sentences            int           `json:"sentences"`
	CharactersWithSpaces int           `json:"charactersWithSpaces"`
	CharactersNoSpaces   int           `json:"charactersNoSpaces"`
	Paragraphs           int           `json:"paragraphs"`
	ReadingTime          time.Duration `json:"readingTime"`
	ReadingTimeFormatted string        `json:"readingTimeFormatted"`
}

// CountAll returns comprehensive statistics for the given text
func CountAll(text string) WordCounterStats {
	return WordCounterStats{
		Words:                CountWords(text),
		Sentences:            CountSentences(text),
		CharactersWithSpaces: CountCharactersWithSpaces(text),
		CharactersNoSpaces:   CountCharactersNoSpaces(text),
		Paragraphs:           CountParagraphs(text),
		ReadingTime:          CalculateReadingTime(text),
		ReadingTimeFormatted: FormatReadingTime(CalculateReadingTime(text)),
	}
}

// CountWords counts the number of words in the text
func CountWords(text string) int {
	if strings.TrimSpace(text) == "" {
		return 0
	}

	// Split by whitespace and filter out empty strings
	words := strings.Fields(text)
	validWords := 0

	for _, word := range words {
		// Remove punctuation and check if there are actual letters/numbers
		cleanedWord := cleanWord(word)
		if cleanedWord != "" {
			validWords++
		}
	}

	return validWords
}

// CountSentences counts the number of sentences in the text
func CountSentences(text string) int {
	if strings.TrimSpace(text) == "" {
		return 0
	}

	// Regex to match sentence endings (.!?) followed by whitespace or end of string
	// This handles cases like "Hello world. How are you?" and "What?! Really..."
	sentenceRegex := regexp.MustCompile(`[.!?]+(?:\s|$)`)
	matches := sentenceRegex.FindAllString(text, -1)

	count := len(matches)

	// If text doesn't end with sentence punctuation but has content, count it as one sentence
	trimmed := strings.TrimSpace(text)
	if count == 0 && trimmed != "" {
		return 1
	}

	// Handle case where text ends without proper punctuation
	lastChar := ""
	if len(trimmed) > 0 {
		lastChar = string(trimmed[len(trimmed)-1])
	}

	if count > 0 && lastChar != "." && lastChar != "!" && lastChar != "?" {
		// Check if there's content after the last sentence marker
		lastMatch := sentenceRegex.FindAllStringIndex(text, -1)
		if len(lastMatch) > 0 {
			lastMatchEnd := lastMatch[len(lastMatch)-1][1]
			if lastMatchEnd < len(text) && strings.TrimSpace(text[lastMatchEnd:]) != "" {
				count++
			}
		}
	}

	return count
}

// CountCharactersWithSpaces counts all characters including spaces
func CountCharactersWithSpaces(text string) int {
	return len(text)
}

// CountCharactersNoSpaces counts all characters excluding whitespace
func CountCharactersNoSpaces(text string) int {
	count := 0
	for _, char := range text {
		if !unicode.IsSpace(char) {
			count++
		}
	}
	return count
}

// CountParagraphs counts the number of paragraphs in the text
func CountParagraphs(text string) int {
	if strings.TrimSpace(text) == "" {
		return 0
	}

	// Split by double newlines or more to identify paragraph breaks
	paragraphRegex := regexp.MustCompile(`\n\s*\n`)
	paragraphs := paragraphRegex.Split(text, -1)

	// Filter out empty paragraphs
	validParagraphs := 0
	for _, paragraph := range paragraphs {
		if strings.TrimSpace(paragraph) != "" {
			validParagraphs++
		}
	}

	// If no paragraph breaks found but text exists, it's one paragraph
	if validParagraphs == 0 && strings.TrimSpace(text) != "" {
		return 1
	}

	return validParagraphs
}

// CalculateReadingTime estimates reading time based on word count
func CalculateReadingTime(text string) time.Duration {
	wordCount := CountWords(text)
	if wordCount == 0 {
		return 0
	}

	minutes := float64(wordCount) / float64(DefaultWordsPerMinute)
	seconds := minutes * 60

	return time.Duration(seconds) * time.Second
}

// CalculateReadingTimeCustom estimates reading time with custom words per minute
func CalculateReadingTimeCustom(text string, wordsPerMinute int) time.Duration {
	wordCount := CountWords(text)
	if wordCount == 0 || wordsPerMinute <= 0 {
		return 0
	}

	minutes := float64(wordCount) / float64(wordsPerMinute)
	seconds := minutes * 60

	return time.Duration(seconds) * time.Second
}

// FormatReadingTime formats the reading time duration into a human-readable string
func FormatReadingTime(duration time.Duration) string {
	if duration == 0 {
		return "0m 0s"
	}

	totalSeconds := int(duration.Seconds())

	if totalSeconds < 60 {
		return "0m 0s" // Less than a minute shows as 0m 0s like in Grammarly
	}

	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	if seconds == 0 {
		return formatMinutes(minutes)
	}

	return formatMinutes(minutes) + " " + formatSeconds(seconds)
}

// cleanWord removes punctuation from a word and returns the cleaned version
func cleanWord(word string) string {
	var result strings.Builder

	for _, char := range word {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// formatMinutes formats minutes for display
func formatMinutes(minutes int) string {
	if minutes < 10 {
		return strings.Replace("Xm", "X", string(rune('0'+minutes)), 1)
	}
	return strings.Replace("XXm", "XX", formatTwoDigits(minutes), 1)
}

// formatSeconds formats seconds for display
func formatSeconds(seconds int) string {
	if seconds < 10 {
		return strings.Replace("Xs", "X", string(rune('0'+seconds)), 1)
	}
	return strings.Replace("XXs", "XX", formatTwoDigits(seconds), 1)
}

// formatTwoDigits formats a number as two digits
func formatTwoDigits(num int) string {
	if num < 10 {
		return "0" + string(rune('0'+num))
	}
	return string(rune('0'+num/10)) + string(rune('0'+num%10))
}

// Advanced functionality for production use

// CountWordsAdvanced provides more sophisticated word counting with language detection
func CountWordsAdvanced(text string, options WordCountOptions) int {
	if options.ExcludePunctuation {
		// Remove all punctuation before counting
		reg := regexp.MustCompile(`[^\p{L}\p{N}\s]+`)
		text = reg.ReplaceAllString(text, " ")
	}

	if options.ExcludeNumbers {
		// Remove standalone numbers
		reg := regexp.MustCompile(`\b\d+\b`)
		text = reg.ReplaceAllString(text, " ")
	}

	return CountWords(text)
}

// WordCountOptions provides configuration for advanced word counting
type WordCountOptions struct {
	ExcludePunctuation bool
	ExcludeNumbers     bool
	MinWordLength      int
	CaseSensitive      bool
}

// GetWordFrequency returns a map of words and their frequencies
func GetWordFrequency(text string, options WordCountOptions) map[string]int {
	frequency := make(map[string]int)

	if strings.TrimSpace(text) == "" {
		return frequency
	}

	words := strings.Fields(text)

	for _, word := range words {
		cleanedWord := cleanWord(word)

		if len(cleanedWord) < options.MinWordLength {
			continue
		}

		if !options.CaseSensitive {
			cleanedWord = strings.ToLower(cleanedWord)
		}

		if cleanedWord != "" {
			frequency[cleanedWord]++
		}
	}

	return frequency
}

// IsEmpty checks if the text is effectively empty (only whitespace)
func IsEmpty(text string) bool {
	return strings.TrimSpace(text) == ""
}

// GetTextSummary returns a summary of all text statistics
func GetTextSummary(text string) map[string]interface{} {
	stats := CountAll(text)

	return map[string]interface{}{
		"stats": stats,
		"metadata": map[string]interface{}{
			"isEmpty":                 IsEmpty(text),
			"averageWordsPerSentence": getAverageWordsPerSentence(stats),
			"averageCharsPerWord":     getAverageCharsPerWord(stats),
			"textLength":              len(text),
			"complexity":              getTextComplexity(stats),
		},
	}
}

// Helper functions for summary statistics
func getAverageWordsPerSentence(stats WordCounterStats) float64 {
	if stats.Sentences == 0 {
		return 0
	}
	return float64(stats.Words) / float64(stats.Sentences)
}

func getAverageCharsPerWord(stats WordCounterStats) float64 {
	if stats.Words == 0 {
		return 0
	}
	return float64(stats.CharactersNoSpaces) / float64(stats.Words)
}

func getTextComplexity(stats WordCounterStats) string {
	avgWordsPerSentence := getAverageWordsPerSentence(stats)
	avgCharsPerWord := getAverageCharsPerWord(stats)

	// Simple complexity scoring
	complexityScore := avgWordsPerSentence*0.6 + avgCharsPerWord*0.4

	switch {
	case complexityScore < 10:
		return "Simple"
	case complexityScore < 15:
		return "Medium"
	default:
		return "Complex"
	}
}
