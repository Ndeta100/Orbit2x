package lorem

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// LoremRequest represents user input for Lorem Ipsum generation
type LoremRequest struct {
	Type      string `json:"type" form:"type"`             // words, sentences, paragraphs
	Count     int    `json:"count" form:"count"`           // number to generate
	StartWith bool   `json:"start_with" form:"start_with"` // start with "Lorem ipsum dolor sit amet..."
}

// LoremResponse represents the generated Lorem Ipsum text
type LoremResponse struct {
	Success   bool   `json:"success"`
	Content   string `json:"content,omitempty"`
	Type      string `json:"type,omitempty"`
	Count     int    `json:"count,omitempty"`
	WordCount int    `json:"word_count,omitempty"`
	CharCount int    `json:"char_count,omitempty"`
	Message   string `json:"message,omitempty"`
}

// Configuration constants
const (
	MaxWords                 = 10000
	MaxSentences             = 1000
	MaxParagraphs            = 500
	MinWordsPerSentence      = 4
	MaxWordsPerSentence      = 18
	MinSentencesPerParagraph = 3
	MaxSentencesPerParagraph = 8
)

// Classic Lorem Ipsum word bank - authentic Latin-based words
var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et",
	"dolore", "magna", "aliqua", "enim", "ad", "minim", "veniam", "quis",
	"nostrud", "exercitation", "ullamco", "laboris", "nisi", "aliquip", "ex", "ea",
	"commodo", "consequat", "duis", "aute", "irure", "in", "reprehenderit", "voluptate",
	"velit", "esse", "cillum", "fugiat", "nulla", "pariatur", "excepteur", "sint",
	"occaecat", "cupidatat", "non", "proident", "sunt", "culpa", "qui", "officia",
	"deserunt", "mollit", "anim", "id", "est", "laborum", "at", "vero",
	"eos", "accusamus", "accusantium", "doloremque", "laudantium", "totam", "rem", "aperiam",
	"eaque", "ipsa", "quae", "ab", "illo", "inventore", "veritatis", "et",
	"quasi", "architecto", "beatae", "vitae", "dicta", "sunt", "explicabo", "nemo",
	"ipsam", "voluptatem", "quia", "voluptas", "aspernatur", "aut", "odit", "fugit",
	"sed", "quia", "consequuntur", "magni", "dolores", "ratione", "sequi", "neque",
	"porro", "quisquam", "dolorem", "adipisci", "numquam", "eius", "modi", "tempora",
	"incidunt", "magnam", "quaerat", "voluptatem", "fuga", "harum", "quidem", "rerum",
	"facilis", "expedita", "distinctio", "nam", "libero", "tempore", "cum", "soluta",
	"nobis", "eleifend", "option", "congue", "nihil", "imperdiet", "doming", "placerat",
	"facer", "possim", "assum", "typi", "non", "habent", "claritatem", "insitam",
	"processus", "dynamicus", "sequitur", "mutationem", "consuetudium", "lectorum", "mirum", "notare",
	"quam", "littera", "gothica", "putamus", "parum", "claram", "anteposuerit", "litterarum",
	"formas", "humanitatis", "seacula", "quarta", "decima", "quinta", "decima", "eodem",
	"modo", "typi", "qui", "nunc", "nobis", "videntur", "parum", "clari",
	"fiant", "sollemnes", "futurum", "claritas", "processus", "dynamicus", "qui", "sequitur",
	"mutationem", "consuetudium", "lectorum", "investigationes", "demonstraverunt", "lectores", "legere", "me",
	"lius", "quod", "ii", "legunt", "saepius", "claritas", "est", "etiam",
}

// Starting phrase for traditional Lorem Ipsum
const loremStart = "Lorem ipsum dolor sit amet, consectetur adipiscing elit"

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Generate creates Lorem Ipsum text based on user requirements
func Generate(req LoremRequest) *LoremResponse {
	// Validate input
	if err := validateRequest(req); err != nil {
		return &LoremResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	var content string
	var wordCount int

	switch req.Type {
	case "words":
		content, wordCount = generateWords(req.Count, req.StartWith)
	case "sentences":
		content, wordCount = generateSentences(req.Count, req.StartWith)
	case "paragraphs":
		content, wordCount = generateParagraphs(req.Count, req.StartWith)
	default:
		content, wordCount = generateParagraphs(req.Count, req.StartWith)
	}

	return &LoremResponse{
		Success:   true,
		Content:   content,
		Type:      req.Type,
		Count:     req.Count,
		WordCount: wordCount,
		CharCount: len(content),
	}
}

// validateRequest ensures the request parameters are within acceptable limits
func validateRequest(req LoremRequest) error {
	if req.Count <= 0 {
		return fmt.Errorf("count must be greater than 0")
	}

	switch req.Type {
	case "words":
		if req.Count > MaxWords {
			return fmt.Errorf("maximum %d words allowed", MaxWords)
		}
	case "sentences":
		if req.Count > MaxSentences {
			return fmt.Errorf("maximum %d sentences allowed", MaxSentences)
		}
	case "paragraphs":
		if req.Count > MaxParagraphs {
			return fmt.Errorf("maximum %d paragraphs allowed", MaxParagraphs)
		}
	default:
		// Default to paragraphs, but still validate
		if req.Count > MaxParagraphs {
			return fmt.Errorf("maximum %d paragraphs allowed", MaxParagraphs)
		}
	}

	return nil
}

// generateWords creates the specified number of Lorem Ipsum words
func generateWords(count int, startWith bool) (string, int) {
	words := make([]string, 0, count)
	wordCount := 0

	// Add a traditional start if requested
	if startWith {
		startWords := strings.Fields(loremStart)
		if count <= len(startWords) {
			return strings.Join(startWords[:count], " "), count
		}
		words = append(words, startWords...)
		wordCount = len(startWords)
		count -= wordCount
	}

	// Generate remaining words
	for i := 0; i < count; i++ {
		words = append(words, getRandomWord())
	}

	return strings.Join(words, " "), wordCount + count
}

// generateSentences creates the specified number of Lorem Ipsum sentences
func generateSentences(count int, startWith bool) (string, int) {
	sentences := make([]string, 0, count)
	totalWordCount := 0

	for i := 0; i < count; i++ {
		var sentence string
		var wordCount int

		if i == 0 && startWith {
			// The first sentence starts with traditional Lorem Ipsum
			sentence, wordCount = createSentenceWithStart()
		} else {
			sentence, wordCount = createSentence()
		}

		sentences = append(sentences, sentence)
		totalWordCount += wordCount
	}

	return strings.Join(sentences, " "), totalWordCount
}

// generateParagraphs creates the specified number of Lorem Ipsum paragraphs
func generateParagraphs(count int, startWith bool) (string, int) {
	paragraphs := make([]string, 0, count)
	totalWordCount := 0

	for i := 0; i < count; i++ {
		paragraph, wordCount := createParagraph(i == 0 && startWith)
		paragraphs = append(paragraphs, paragraph)
		totalWordCount += wordCount
	}

	return strings.Join(paragraphs, "\n\n"), totalWordCount
}

// createSentence generates a single sentence with proper capitalization and punctuation
func createSentence() (string, int) {
	wordCount := rng.Intn(MaxWordsPerSentence-MinWordsPerSentence+1) + MinWordsPerSentence
	words := make([]string, wordCount)

	for i := 0; i < wordCount; i++ {
		word := getRandomWord()
		if i == 0 {
			word = strings.Title(word)
		}
		words[i] = word
	}

	sentence := strings.Join(words, " ")

	// Add occasional commas for realism (about 30% chance after word 3+)
	if wordCount > 5 && rng.Float32() < 0.3 {
		commaPos := rng.Intn(wordCount-3) + 2
		words[commaPos] += ","
		sentence = strings.Join(words, " ")
	}

	return sentence + ".", wordCount
}

// createSentenceWithStart creates a sentence that begins with traditional Lorem Ipsum
func createSentenceWithStart() (string, int) {
	startWords := strings.Fields(loremStart)
	additionalWords := rng.Intn(8) + 2 // Add 2-9 more words

	words := make([]string, len(startWords)+additionalWords)
	copy(words, startWords)

	for i := len(startWords); i < len(words); i++ {
		words[i] = getRandomWord()
	}

	sentence := strings.Join(words, " ")
	return sentence + ".", len(words)
}

// createParagraph generates a paragraph with multiple sentences
func createParagraph(startWithLorem bool) (string, int) {
	sentenceCount := rng.Intn(MaxSentencesPerParagraph-MinSentencesPerParagraph+1) + MinSentencesPerParagraph
	sentences := make([]string, sentenceCount)
	totalWordCount := 0

	for i := 0; i < sentenceCount; i++ {
		var sentence string
		var wordCount int

		if i == 0 && startWithLorem {
			sentence, wordCount = createSentenceWithStart()
		} else {
			sentence, wordCount = createSentence()
		}

		sentences[i] = sentence
		totalWordCount += wordCount
	}

	return strings.Join(sentences, " "), totalWordCount
}

// getRandomWord returns a random word from the Lorem Ipsum vocabulary
func getRandomWord() string {
	return loremWords[rng.Intn(len(loremWords))]
}

// GetWordBank returns the complete Lorem Ipsum word bank (useful for frontend preview)
func GetWordBank() []string {
	return loremWords
}

// GetLimits returns the maximum limits for each type (useful for frontend validation)
func GetLimits() map[string]int {
	return map[string]int{
		"words":      MaxWords,
		"sentences":  MaxSentences,
		"paragraphs": MaxParagraphs,
	}
}
