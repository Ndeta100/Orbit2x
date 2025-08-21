package resolver

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Cache implementation
type geoIPCache struct {
	cache map[string]cacheEntry
	mu    sync.RWMutex
}

type cacheEntry struct {
	data      *geoIPResponse
	timestamp time.Time
}

// GeoIP API response structure
type geoIPResponse struct {
	IP          string `json:"query"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"country"`
	RegionCode  string `json:"region"`
	RegionName  string `json:"regionName"`
	CityName    string `json:"city"`
	ISP         string `json:"isp"`
}

// Global cache with 1-hour expiration
var (
	ipCache = &geoIPCache{
		cache: make(map[string]cacheEntry),
	}
	cacheDuration = 1 * time.Hour
)

// getRealIP determines the client's real IPv4 address
func GetRealIP(r *http.Request) string {
	// Check common _headers that might contain the real IP
	headersToCheck := []string{
		"CF-Connecting-IP", // Cloudflare
		"X-Forwarded-For",  // Standard proxy header
		"X-Real-IP",        // NGINX proxy
		"X-Client-IP",      // Apache
		"True-Client-IP",   // Akamai and Cloudflare
		"Forwarded-For",
		"X-Forwarded",
	}

	for _, header := range headersToCheck {
		if value := r.Header.Get(header); value != "" {
			// Headers might contain multiple IPs (e.g., X-Forwarded-For: client, proxy1, proxy2)
			// Take the first one, which is typically the client IP
			ips := strings.Split(value, ",")
			ip := strings.TrimSpace(ips[0])

			// Validate it's an IPv4 address
			parsedIP := net.ParseIP(ip)
			if parsedIP != nil && parsedIP.To4() != nil {
				return ip
			}
		}
	}

	// Fall back to RemoteAddr if no valid IP in _headers
	ipPort := r.RemoteAddr
	ip, _, err := net.SplitHostPort(ipPort)
	if err != nil {
		// If we can't split, just return as is
		return ipPort
	}

	return ip
}

// getIPv6 tries to detect the client's IPv6 address
func GetIPv6(r *http.Request) string {
	// Check X-Forwarded-For first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		for _, ipStr := range ips {
			ipStr = strings.TrimSpace(ipStr)
			ip := net.ParseIP(ipStr)
			if ip != nil && ip.To4() == nil {
				return ipStr // Found an IPv6 address
			}
		}
	}

	// Try to get from RemoteAddr
	ipStr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		ip := net.ParseIP(ipStr)
		if ip != nil && ip.To4() == nil {
			return ipStr // Found an IPv6 address
		}
	}

	// No IPv6 found
	return ""
}

// getIPGeolocation gets location and ISP data for an IP address
func GetIPGeolocation(ipAddress string) (*geoIPResponse, error) {
	// Check cache first
	ipCache.mu.RLock()
	if entry, found := ipCache.cache[ipAddress]; found {
		if time.Since(entry.timestamp) < cacheDuration {
			ipCache.mu.RUnlock()
			return entry.data, nil
		}
	}
	ipCache.mu.RUnlock()

	// Call IP-API for geolocation data
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode,region,regionName,city,isp,query", ipAddress)

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	// Read and parse the response
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // Limit to 1MB
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %w", err)
	}

	// Parse the JSON
	var result struct {
		Status      string `json:"status"`
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
		Region      string `json:"region"`
		RegionName  string `json:"regionName"`
		City        string `json:"city"`
		ISP         string `json:"isp"`
		Query       string `json:"query"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Check if the API returned an error
	if result.Status != "success" {
		return nil, fmt.Errorf("API returned error status")
	}

	// Create response object
	geoResponse := &geoIPResponse{
		IP:          result.Query,
		CountryCode: result.CountryCode,
		CountryName: result.Country,
		RegionCode:  result.Region,
		RegionName:  result.RegionName,
		CityName:    result.City,
		ISP:         result.ISP,
	}

	// Update cache
	ipCache.mu.Lock()
	ipCache.cache[ipAddress] = cacheEntry{
		data:      geoResponse,
		timestamp: time.Now(),
	}
	ipCache.mu.Unlock()

	return geoResponse, nil
}

// InitCache initializes the cache and starts a cleanup routine
func InitCache() {
	// Start a background goroutine to clean up the cache every 12 hours
	go func() {
		ticker := time.NewTicker(12 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			cleanupCache()
		}
	}()
}

// cleanupCache removes expired entries from the cache
func cleanupCache() {
	ipCache.mu.Lock()
	defer ipCache.mu.Unlock()

	for ip, entry := range ipCache.cache {
		if time.Since(entry.timestamp) > cacheDuration {
			delete(ipCache.cache, ip)
		}
	}
}
