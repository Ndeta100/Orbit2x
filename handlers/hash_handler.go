// handlers/hash_handler.go
package handlers

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"hash/crc32"
	"io"
	"net/http"
	"strings"

	"github.com/Ndeta100/orbit2x/views/hashgen" // Adjust to your actual path
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// HandleHashIndex renders the Hash Generator page
func HandleHashIndex(w http.ResponseWriter, r *http.Request) error {
	return hashgen.Index().Render(r.Context(), w)
}

// HandleGenerateHash generates various hash values for the provided input
func HandleGenerateHash(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return hashgen.Results(hashgen.HashResult{
			Error: "Failed to parse form data",
		}).Render(r.Context(), w)
	}

	// Get text input from form
	text := r.FormValue("text")
	if text == "" {
		return hashgen.Results(hashgen.HashResult{
			Error: "Input text is required",
		}).Render(r.Context(), w)
	}

	// Selected hash algorithms
	selectedAlgorithms := r.Form["algorithms"]
	if len(selectedAlgorithms) == 0 {
		// If no algorithm is selected, use all
		selectedAlgorithms = []string{"md5", "sha1", "sha256", "sha512", "sha3-256", "sha3-512", "blake2b", "ripemd160", "crc32", "bcrypt"}
	}

	// Prepare result
	result := hashgen.HashResult{
		InputText:  text,
		HashValues: make(map[string]string),
	}

	// Generate hashes for selected algorithms
	data := []byte(text)
	for _, algo := range selectedAlgorithms {
		switch algo {
		case "md4":
			result.HashValues["MD4"] = generateHash(md4.New(), data)
		case "md5":
			result.HashValues["MD5"] = generateHash(md5.New(), data)
		case "sha1":
			result.HashValues["SHA-1"] = generateHash(sha1.New(), data)
		case "sha224":
			result.HashValues["SHA-224"] = generateHash(sha256.New224(), data)
		case "sha256":
			result.HashValues["SHA-256"] = generateHash(sha256.New(), data)
		case "sha384":
			result.HashValues["SHA-384"] = generateHash(sha512.New384(), data)
		case "sha512":
			result.HashValues["SHA-512"] = generateHash(sha512.New(), data)
		case "sha3-224":
			result.HashValues["SHA3-224"] = generateHash(sha3.New224(), data)
		case "sha3-256":
			result.HashValues["SHA3-256"] = generateHash(sha3.New256(), data)
		case "sha3-384":
			result.HashValues["SHA3-384"] = generateHash(sha3.New384(), data)
		case "sha3-512":
			result.HashValues["SHA3-512"] = generateHash(sha3.New512(), data)
		case "blake2b":
			blake2bHash, _ := blake2b.New512(nil)
			result.HashValues["BLAKE2b-512"] = generateHash(blake2bHash, data)
		case "ripemd160":
			result.HashValues["RIPEMD-160"] = generateHash(ripemd160.New(), data)
		case "crc32":
			crc32Hash := crc32.NewIEEE()
			crc32Hash.Write(data)
			result.HashValues["CRC32"] = hex.EncodeToString(crc32Hash.Sum(nil))
		case "bcrypt":
			bcryptHash, err := bcrypt.GenerateFromPassword(data, bcrypt.DefaultCost)
			if err != nil {
				result.HashValues["bcrypt"] = "Error generating bcrypt hash: " + err.Error()
			} else {
				result.HashValues["bcrypt"] = string(bcryptHash)
			}
		}
	}

	// Render the result
	return hashgen.Results(result).Render(r.Context(), w)
}

// HandleFileHash generates hash values for a file
func HandleFileHash(w http.ResponseWriter, r *http.Request) error {
	// Set max upload size - 10MB
	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		return hashgen.Results(hashgen.HashResult{
			Error: "The uploaded file is too large. Maximum size is 10MB.",
		}).Render(r.Context(), w)
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		return hashgen.Results(hashgen.HashResult{
			Error: "Failed to get file from form: " + err.Error(),
		}).Render(r.Context(), w)
	}
	defer file.Close()

	// Selected hash algorithms
	selectedAlgorithms := r.Form["algorithms"]
	if len(selectedAlgorithms) == 0 {
		// If no algorithm is selected, use common ones for files
		selectedAlgorithms = []string{"md5", "sha1", "sha256", "sha512"}
	}

	// Prepare result
	result := hashgen.HashResult{
		InputText:  "File: " + header.Filename,
		HashValues: make(map[string]string),
		IsFile:     true,
		FileName:   header.Filename,
		FileSize:   formatFileSize(header.Size),
	}

	// Generate file hashes for selected algorithms
	for _, algo := range selectedAlgorithms {
		// Seek back to start of file for each hash calculation
		file.Seek(0, io.SeekStart)

		switch algo {
		case "md4":
			result.HashValues["MD4"] = generateFileHash(md4.New(), file)
		case "md5":
			result.HashValues["MD5"] = generateFileHash(md5.New(), file)
		case "sha1":
			result.HashValues["SHA-1"] = generateFileHash(sha1.New(), file)
		case "sha224":
			result.HashValues["SHA-224"] = generateFileHash(sha256.New224(), file)
		case "sha256":
			result.HashValues["SHA-256"] = generateFileHash(sha256.New(), file)
		case "sha384":
			result.HashValues["SHA-384"] = generateFileHash(sha512.New384(), file)
		case "sha512":
			result.HashValues["SHA-512"] = generateFileHash(sha512.New(), file)
		case "sha3-224":
			result.HashValues["SHA3-224"] = generateFileHash(sha3.New224(), file)
		case "sha3-256":
			result.HashValues["SHA3-256"] = generateFileHash(sha3.New256(), file)
		case "sha3-384":
			result.HashValues["SHA3-384"] = generateFileHash(sha3.New384(), file)
		case "sha3-512":
			result.HashValues["SHA3-512"] = generateFileHash(sha3.New512(), file)
		case "blake2b":
			blake2bHash, _ := blake2b.New512(nil)
			result.HashValues["BLAKE2b-512"] = generateFileHash(blake2bHash, file)
		case "ripemd160":
			result.HashValues["RIPEMD-160"] = generateFileHash(ripemd160.New(), file)
		case "crc32":
			crc32Hash := crc32.NewIEEE()
			io.Copy(crc32Hash, file)
			result.HashValues["CRC32"] = hex.EncodeToString(crc32Hash.Sum(nil))
		}
		// Skip bcrypt for files as it's designed for passwords and has input size limitations
	}

	// Render the result
	return hashgen.Results(result).Render(r.Context(), w)
}

// HandleVerifyHash verifies if a hash matches the input
func HandleVerifyHash(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return hashgen.VerifyResult(hashgen.HashVerifyResult{
			Error: "Failed to parse form data",
		}).Render(r.Context(), w)
	}

	// Get input data from form
	input := r.FormValue("input")
	hashForm := r.FormValue("hash")
	algorithm := r.FormValue("algorithm")

	if input == "" || hashForm == "" || algorithm == "" {
		return hashgen.VerifyResult(hashgen.HashVerifyResult{
			Error: "Input, hash, and algorithm are all required",
		}).Render(r.Context(), w)
	}

	// Normalize hash input (remove spaces, make lowercase)
	hashForm = strings.ToLower(strings.ReplaceAll(hashForm, " ", ""))

	// Calculate the hash of the input
	var calculatedHash string
	data := []byte(input)

	switch algorithm {
	case "md5":
		calculatedHash = generateHash(md5.New(), data)
	case "sha1":
		calculatedHash = generateHash(sha1.New(), data)
	case "sha256":
		calculatedHash = generateHash(sha256.New(), data)
	case "sha512":
		calculatedHash = generateHash(sha512.New(), data)
	case "bcrypt":
		// Special case for bcrypt as it uses a different verification method
		err := bcrypt.CompareHashAndPassword([]byte(hashForm), data)
		result := hashgen.HashVerifyResult{
			Input:     input,
			Hash:      hashForm,
			Algorithm: algorithm,
			Matches:   err == nil,
		}
		return hashgen.VerifyResult(result).Render(r.Context(), w)
	default:
		return hashgen.VerifyResult(hashgen.HashVerifyResult{
			Error: "Unsupported algorithm for verification",
		}).Render(r.Context(), w)
	}

	// For all other algorithms, compare the calculated hash with the provided hash
	result := hashgen.HashVerifyResult{
		Input:          input,
		Hash:           hashForm,
		Algorithm:      algorithm,
		Matches:        strings.EqualFold(calculatedHash, hashForm),
		CalculatedHash: calculatedHash,
	}

	// Render the result
	return hashgen.VerifyResult(result).Render(r.Context(), w)
}

// generateHash generates a hash for text input
func generateHash(h hash.Hash, data []byte) string {
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// generateFileHash generates a hash for a file
func generateFileHash(h hash.Hash, file io.Reader) string {
	io.Copy(h, file)
	return hex.EncodeToString(h.Sum(nil))
}
