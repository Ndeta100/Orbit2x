package handlers

import (
	"encoding/json"
	"errors"
	_ "html/template"
	"net/http"
	"strings"

	"github.com/Ndeta100/orbit2x/internal/resolver"
	"github.com/Ndeta100/orbit2x/views/dns"
	"github.com/Ndeta100/orbit2x/views/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}

func HandleDNSLookupIndex(w http.ResponseWriter, r *http.Request) error {
	return dns.DNSLookup().Render(r.Context(), w)
}

func HandleDNSLookup(w http.ResponseWriter, r *http.Request) error {
	var domain string

	// Handle both form and JSON data
	if r.Header.Get("Content-Type") == "application/json" {
		var data struct {
			Domain string `json:"domain"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return errors.New("invalid JSON")
		}
		domain = data.Domain
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return errors.New("failed to parse form")
		}
		domain = r.FormValue("domain")
	}

	if domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return errors.New("domain is required")
	}

	// Clean the domain
	domain = strings.TrimSpace(domain)
	domain = strings.ToLower(domain)
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "www.")
	domain = strings.Split(domain, "/")[0]

	// Create resolver and perform lookups
	dnsResolver := &resolver.DefaultResolver{}
	results := resolver.PerformAllLookups(dnsResolver, domain)

	// Render results using templ
	dns.DNSResults(domain, results).Render(r.Context(), w)
	return nil
}
