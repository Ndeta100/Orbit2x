package main

import (
	"github.com/Ndeta100/orbit2x/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type ViewData struct {
	Title string
}

func main() {
	if err := loadEnv(); err != nil {
		log.Fatal("Error loading .env file")
	}
	router := chi.NewMux()
	router.Get("/", handlers.Make(handlers.HandleHomeIndex))
	router.Post("/lookup", handlers.HandleDNSLookup)
	router.Get("/myip", handlers.Make(handlers.HandleMyIP))
	router.Get("/headers", handlers.Make(handlers.HandleHeadersIndex))
	router.Post("/headers/analyze", handlers.Make(handlers.HandleHeadersAnalyze))
	router.Get("/ssl", handlers.Make(handlers.HandleSSLIndex))
	router.Post("/ssl/check", handlers.Make(handlers.HandleSSLCheck))
	router.Get("/subnet", handlers.Make(handlers.HandleSubnetIndex))
	router.Post("/subnet/calculate-cidr", handlers.Make(handlers.HandleSubnetCalculateCIDR))
	router.Post("/subnet/calculate-mask", handlers.Make(handlers.HandleSubnetCalculateMask))
	router.Get("/encoder", handlers.Make(handlers.HandleEncoderIndex))
	router.Post("/encoder/encode", handlers.Make(handlers.HandleEncoderEncode))
	router.Post("/encoder/decode", handlers.Make(handlers.HandleEncoderDecode))
	router.Get("/formatter", handlers.Make(handlers.HandleFormatterIndex))
	router.Post("/formatter/json", handlers.Make(handlers.HandleJSONFormat))
	router.Post("/formatter/yaml", handlers.Make(handlers.HandleYAMLFormat))
	router.Get("/converter", handlers.Make(handlers.HandleConverterIndex))
	router.Post("/converter/csv-to-json", handlers.Make(handlers.HandleCSVToJSON))
	router.Post("/converter/json-to-csv", handlers.Make(handlers.HandleJSONToCSV))
	router.Get("/useragent", handlers.Make(handlers.HandleUserAgentIndex))
	router.Post("/useragent/parse", handlers.Make(handlers.HandleUserAgentParse))
	router.Get("/useragent/detect", handlers.Make(handlers.HandleDetectUserAgent))
	router.Get("/imagebase64", handlers.Make(handlers.HandleImageBase64Index))
	router.Post("/imagebase64/convert", handlers.Make(handlers.HandleImageBase64Convert))
	router.Get("/hash", handlers.Make(handlers.HandleHashIndex))
	router.Post("/hash/generate", handlers.Make(handlers.HandleGenerateHash))
	router.Post("/hash/file", handlers.Make(handlers.HandleFileHash))
	router.Get("/color", handlers.Make(handlers.HandleColorIndex))
	router.Post("/color/convert", handlers.Make(handlers.HandleColorConvert))
	router.Get("/color/random", handlers.Make(handlers.HandleRandomColor))
	port := os.Getenv("HTTP_LISTEN_ADR")
	slog.Info("Application is running on port %s", port)
	log.Fatal(http.ListenAndServe(port, router))

}

func loadEnv() error {
	return godotenv.Load()
}
