package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/Ndeta100/orbit2x/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type ViewData struct {
	Title string
}

func main() {
	if err := loadEnv(); err != nil {
		// Don't fatal on missing .env in production
		log.Println("No .env file found, using environment variables")
	}
	router := chi.NewMux()

	//Tools
	router.Get("/", handlers.Make(handlers.HandleHomeIndex))
	router.Get("/lookup", handlers.Make(handlers.HandleDNSLookupIndex))
	router.Post("/lookup", handlers.Make(handlers.HandleDNSLookup))
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
	//User agent
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
	//qr_code generator
	router.Get("/qr", handlers.Make(handlers.ShowQRTool))
	router.Post("/qr/generate", handlers.Make(handlers.GenerateQR))
	router.Post("/qr/preview", handlers.Make(handlers.PreviewQR))

	//lorem route
	router.Get("/lorem", handlers.Make(handlers.HandleLoremMainPage))
	router.Get("/lorem/limits", handlers.Make(handlers.HandleLoremLimits))
	router.Post("/lorem/generate", handlers.Make(handlers.HandleLoremGenerate))
	router.Post("/lorem/api/generate", handlers.Make(handlers.HandleLoremAPIGenerate))

	//Word counter
	router.Get("/word-counter", handlers.Make(handlers.ShowWordCounterPage))
	router.Post("/word-count", handlers.Make(handlers.CountWords))

	//home page routes and others
	router.Get("/privacy", handlers.Make(handlers.PrivacyHandler))
	router.Get("/about", handlers.Make(handlers.AboutHandler))

	//Faq route
	router.Get("/faq", handlers.Make(handlers.FAQHandler))
	router.Get("/contact", handlers.Make(handlers.ContactHandler))

	// tools routes
	router.Get("/tools", handlers.Make(handlers.ToolsHandler))
	//http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
	//	if r.URL.Path == "/categories/" {
	//		http.Redirect(w, r, "/categories", http.StatusMovedPermanently)
	//		return
	//	}
	//	handlers.CategoryDetailHandler(w, r)
	//})
	// Google site verification route
	router.Get("/google34ece54340b6c8d6.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("google-site-verification: google34ece54340b6c8d6.html"))
	})

	router.Get("/dns-tools", handlers.Make(handlers.HandleDNSToolsCategory))
	router.Get("/dev-tools", handlers.Make(handlers.HandleDeveloperToolsCategory))
	router.Get("/design-tools", handlers.Make(handlers.HandleDesignerToolsCategory))
	router.Get("/webmaster-tools", handlers.Make(handlers.HandleWebmasterToolsCategory))
	router.Get("/network-tools", handlers.Make(handlers.HandleNetworkToolsCategory))
	router.Get("/security-tools", handlers.Make(handlers.HandleSecurityToolsCategory))
	router.Get("/productivity-tools", handlers.Make(handlers.HandleProductivityToolsCategory))
	router.Get("/gaming-tools", handlers.Make(handlers.HandleGamingToolsCategory))
	router.Get("/more-tools", handlers.Make(handlers.HandleMoreCategoriesCategory))

	// This will catch ALL 404s - any route not defined above
	router.NotFound(handlers.Make(handlers.HandleComingSoon))

	port := os.Getenv("HTTP_LISTEN_ADR")
	if port == "" {
		port = os.Getenv("HTTP_LISTEN_ADR")
		if port == "" {
			port = ":8080"
		}
	}

	// Make sure port starts with ":"
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	slog.Info("Application is running", "port", port)
	log.Fatal(http.ListenAndServe(port, router))

}

func loadEnv() error {
	return godotenv.Load()
}
