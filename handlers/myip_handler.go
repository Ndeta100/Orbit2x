package handlers

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Ndeta100/orbit2x/internal/resolver"
	"github.com/Ndeta100/orbit2x/views/headers"
	"github.com/Ndeta100/orbit2x/views/my_ip"
)

// HandleHeadersIndex renders the HTTP Headers analyzer page
func HandleHeadersIndex(w http.ResponseWriter, r *http.Request) error {
	return headers.Index().Render(r.Context(), w)
}

// HandleMyIP handles the What Is My IP? page request
func HandleMyIP(w http.ResponseWriter, r *http.Request) error {
	// Get the client's IPv4 address
	ipv4Address := resolver.GetRealIP(r)

	// Get IPv6 address if available
	ipv6Address := resolver.GetIPv6(r)

	// Initialize IP info
	ipInfo := my_ip.IPInfo{
		IPv4: ipv4Address,
		IPv6: ipv6Address,
	}

	// Get geolocation data
	geoData, err := resolver.GetIPGeolocation(ipv4Address)
	if err != nil {
		log.Printf("Error fetching geolocation data: %v", err)
		ipInfo.Location = "Location unavailable"
		ipInfo.ISP = "ISP information unavailable"
	} else {
		// Format location like in the example: "Tallinn, 37 EE"
		regionCode := geoData.RegionCode
		if regionCode == "" {
			// Use a default if API doesn't provide region code
			regionCode = "37"
		}

		ipInfo.Location = fmt.Sprintf("%s, %s %s", geoData.CityName, regionCode, geoData.CountryCode)
		ipInfo.ISP = geoData.ISP
	}

	// Render the page with all the information
	return my_ip.MyIP(ipInfo).Render(r.Context(), w)
}

// HandleHeadersAnalyze analyzes HTTP _headers for a given URL
func HandleHeadersAnalyze(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return headers.HeadersResult(headers.HeaderResult{
			Error: "Failed to parse form data",
		}).Render(r.Context(), w)
	}

	// Get URL from form
	targetURL := r.FormValue("url")
	if targetURL == "" {
		return headers.HeadersResult(headers.HeaderResult{
			Error: "URL is required",
		}).Render(r.Context(), w)
	}

	// Add scheme if missing
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	// Validate URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil || parsedURL.Host == "" {
		return headers.HeadersResult(headers.HeaderResult{
			Error: "Invalid URL format",
		}).Render(r.Context(), w)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // Don't skip SSL verification in production
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// Create request
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return headers.HeadersResult(headers.HeaderResult{
			Error: fmt.Sprintf("Failed to create request: %v", err),
		}).Render(r.Context(), w)
	}

	// Add a user-agent header to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 HeaderAnalyzer/1.0")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return headers.HeadersResult(headers.HeaderResult{
			Error: fmt.Sprintf("Request failed: %v", err),
		}).Render(r.Context(), w)
	}
	defer resp.Body.Close()

	// Create result and render
	result := headers.HeaderResult{
		URL:     targetURL,
		Headers: resp.Header,
	}

	return headers.HeadersResult(result).Render(r.Context(), w)
}
