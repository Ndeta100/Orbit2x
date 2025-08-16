// handlers/useragent_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/Ndeta100/orbit2x/views/useragent" // Adjust to your actual path
	ua "github.com/mileusna/useragent"
)

// HandleUserAgentIndex renders the User Agent Parser page
func HandleUserAgentIndex(w http.ResponseWriter, r *http.Request) error {
	return useragent.Index().Render(r.Context(), w)
}

// HandleUserAgentParse parses a user agent string and returns the details
func HandleUserAgentParse(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return useragent.Results(useragent.ParsingResult{
			Error: "Failed to parse form data",
		}).Render(r.Context(), w)
	}

	// Get user agent string from form
	userAgentString := r.FormValue("userAgent")
	if userAgentString == "" {
		return useragent.Results(useragent.ParsingResult{
			Error: "User agent string is required",
		}).Render(r.Context(), w)
	}

	// Parse the user agent string
	parsedUA := ua.Parse(userAgentString)

	// Build the detailed result
	result := useragent.ParsingResult{
		OriginalUA: userAgentString,
		Browser: useragent.BrowserInfo{
			Name:    parsedUA.Name,
			Version: parsedUA.Version,
		},
		Device: useragent.DeviceInfo{
			Type:      getDeviceType(parsedUA),
			OS:        parsedUA.OS,
			OSVersion: parsedUA.OSVersion,
		},
		IsMobile:  parsedUA.Mobile,
		IsTablet:  parsedUA.Tablet,
		IsDesktop: parsedUA.Desktop,
		IsBot:     parsedUA.Bot,
		BotName:   "", // The library doesn't have a BotName field - leaving it empty
		// Add additional information
		FullDetails: map[string]string{
			"Browser":          parsedUA.Name,
			"Browser Version":  parsedUA.Version,
			"Operating System": parsedUA.OS,
			"OS Version":       parsedUA.OSVersion,
			"Device Type":      getDeviceType(parsedUA),
			"Is Mobile":        formatBool(parsedUA.Mobile),
			"Is Tablet":        formatBool(parsedUA.Tablet),
			"Is Desktop":       formatBool(parsedUA.Desktop),
			"Is Bot":           formatBool(parsedUA.Bot),
		},
	}

	// Add specific browser checks
	// The library doesn't have methods like Chrome(), Firefox(), etc.
	// Instead, we'll check the browser name against known values
	if parsedUA.Name == ua.Chrome {
		result.FullDetails["Is Chrome"] = "Yes"
	}
	if parsedUA.Name == ua.Firefox {
		result.FullDetails["Is Firefox"] = "Yes"
	}
	if parsedUA.Name == ua.Safari || parsedUA.Name == ua.MobileSafari {
		result.FullDetails["Is Safari"] = "Yes"
	}
	if parsedUA.Name == ua.Opera || parsedUA.Name == ua.OperaMini {
		result.FullDetails["Is Opera"] = "Yes"
	}
	if parsedUA.Name == ua.InternetExplorer || parsedUA.Name == ua.Msie {
		result.FullDetails["Is Internet Explorer"] = "Yes"
	}
	if parsedUA.Name == ua.Edge {
		result.FullDetails["Is Edge"] = "Yes"
	}

	// Add specific OS checks
	// The library doesn't have methods like Windows(), Mac(), etc.
	// Instead, we'll check the OS name against known values
	if parsedUA.OS == ua.Windows || parsedUA.OS == ua.WindowsPhone || parsedUA.OS == ua.WindowsNT {
		result.FullDetails["Is Windows"] = "Yes"
	}
	if parsedUA.OS == ua.MacOS {
		result.FullDetails["Is macOS"] = "Yes"
	}
	if parsedUA.OS == ua.Linux {
		result.FullDetails["Is Linux"] = "Yes"
	}
	if parsedUA.OS == ua.Android {
		result.FullDetails["Is Android"] = "Yes"
	}
	if parsedUA.OS == ua.IOS {
		result.FullDetails["Is iOS"] = "Yes"
	}

	// Render the results
	return useragent.Results(result).Render(r.Context(), w)
}

// getDeviceType returns a user-friendly device type description
func getDeviceType(ua ua.UserAgent) string {
	if ua.Mobile {
		return "Mobile"
	}
	if ua.Tablet {
		return "Tablet"
	}
	if ua.Desktop {
		return "Desktop"
	}
	if ua.Bot {
		return "Bot"
	}
	return "Unknown"
}

// formatBool converts a boolean to Yes/No string
func formatBool(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

// HandleDetectUserAgent detects and parses the client's user agent
func HandleDetectUserAgent(w http.ResponseWriter, r *http.Request) error {
	// Get user agent from the request header
	userAgentString := r.Header.Get("User-Agent")

	// Check if the UA string is empty
	if strings.TrimSpace(userAgentString) == "" {
		return useragent.Results(useragent.ParsingResult{
			Error: "No User-Agent header found in your request",
		}).Render(r.Context(), w)
	}

	// Parse the user agent
	parsedUA := ua.Parse(userAgentString)

	// Build the result (same as in HandleUserAgentParse)
	result := useragent.ParsingResult{
		OriginalUA: userAgentString,
		Browser: useragent.BrowserInfo{
			Name:    parsedUA.Name,
			Version: parsedUA.Version,
		},
		Device: useragent.DeviceInfo{
			Type:      getDeviceType(parsedUA),
			OS:        parsedUA.OS,
			OSVersion: parsedUA.OSVersion,
		},
		IsMobile:  parsedUA.Mobile,
		IsTablet:  parsedUA.Tablet,
		IsDesktop: parsedUA.Desktop,
		IsBot:     parsedUA.Bot,
		BotName:   "", // The library doesn't have a BotName field - leaving it empty
		FullDetails: map[string]string{
			"Browser":          parsedUA.Name,
			"Browser Version":  parsedUA.Version,
			"Operating System": parsedUA.OS,
			"OS Version":       parsedUA.OSVersion,
			"Device Type":      getDeviceType(parsedUA),
			"Is Mobile":        formatBool(parsedUA.Mobile),
			"Is Tablet":        formatBool(parsedUA.Tablet),
			"Is Desktop":       formatBool(parsedUA.Desktop),
			"Is Bot":           formatBool(parsedUA.Bot),
		},
	}

	// Add specific browser checks
	// The library doesn't have methods like Chrome(), Firefox(), etc.
	// Instead, we'll check the browser name against known values
	if parsedUA.Name == ua.Chrome {
		result.FullDetails["Is Chrome"] = "Yes"
	}
	if parsedUA.Name == ua.Firefox {
		result.FullDetails["Is Firefox"] = "Yes"
	}
	if parsedUA.Name == ua.Safari || parsedUA.Name == ua.MobileSafari {
		result.FullDetails["Is Safari"] = "Yes"
	}
	if parsedUA.Name == ua.Opera || parsedUA.Name == ua.OperaMini {
		result.FullDetails["Is Opera"] = "Yes"
	}
	if parsedUA.Name == ua.InternetExplorer || parsedUA.Name == ua.Msie {
		result.FullDetails["Is Internet Explorer"] = "Yes"
	}
	if parsedUA.Name == ua.Edge {
		result.FullDetails["Is Edge"] = "Yes"
	}

	// Add specific OS checks
	// The library doesn't have methods like Windows(), Mac(), etc.
	// Instead, we'll check the OS name against known values
	if parsedUA.OS == ua.Windows || parsedUA.OS == ua.WindowsPhone || parsedUA.OS == ua.WindowsNT {
		result.FullDetails["Is Windows"] = "Yes"
	}
	if parsedUA.OS == ua.MacOS {
		result.FullDetails["Is macOS"] = "Yes"
	}
	if parsedUA.OS == ua.Linux {
		result.FullDetails["Is Linux"] = "Yes"
	}
	if parsedUA.OS == ua.Android {
		result.FullDetails["Is Android"] = "Yes"
	}
	if parsedUA.OS == ua.IOS {
		result.FullDetails["Is iOS"] = "Yes"
	}

	// Render the results
	return useragent.Results(result).Render(r.Context(), w)
}
