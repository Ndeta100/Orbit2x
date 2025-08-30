package handlers

import (
	"net/http"
	"strings"

	"github.com/Ndeta100/orbit2x/views/useragent"
	ua "github.com/mileusna/useragent"
)

// HandleUserAgentIndex renders the User Agent Parser page
func HandleUserAgentIndex(w http.ResponseWriter, r *http.Request) error {
	return useragent.UserAgentPage().Render(r.Context(), w)
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
	userAgentString := strings.TrimSpace(r.FormValue("userAgent"))
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
		BotName:   getBotName(parsedUA),
		FullDetails: map[string]string{
			"Browser":          getValueOrEmpty(parsedUA.Name),
			"Browser Version":  getValueOrEmpty(parsedUA.Version),
			"Operating System": getValueOrEmpty(parsedUA.OS),
			"OS Version":       getValueOrEmpty(parsedUA.OSVersion),
			"Device Type":      getDeviceType(parsedUA),
			"Is Mobile":        formatBool(parsedUA.Mobile),
			"Is Tablet":        formatBool(parsedUA.Tablet),
			"Is Desktop":       formatBool(parsedUA.Desktop),
			"Is Bot":           formatBool(parsedUA.Bot),
		},
	}

	// Add browser-specific detection flags
	addBrowserDetection(&result, parsedUA)

	// Add OS-specific detection flags
	addOSDetection(&result, parsedUA)

	// Add additional technical details if available
	addTechnicalDetails(&result, parsedUA)

	// Render the results
	return useragent.Results(result).Render(r.Context(), w)
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

	// Parse the user agent using the same logic as HandleUserAgentParse
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
		BotName:   getBotName(parsedUA),
		FullDetails: map[string]string{
			"Browser":          getValueOrEmpty(parsedUA.Name),
			"Browser Version":  getValueOrEmpty(parsedUA.Version),
			"Operating System": getValueOrEmpty(parsedUA.OS),
			"OS Version":       getValueOrEmpty(parsedUA.OSVersion),
			"Device Type":      getDeviceType(parsedUA),
			"Is Mobile":        formatBool(parsedUA.Mobile),
			"Is Tablet":        formatBool(parsedUA.Tablet),
			"Is Desktop":       formatBool(parsedUA.Desktop),
			"Is Bot":           formatBool(parsedUA.Bot),
		},
	}

	// Add browser-specific detection flags
	addBrowserDetection(&result, parsedUA)

	// Add OS-specific detection flags
	addOSDetection(&result, parsedUA)

	// Add additional technical details
	addTechnicalDetails(&result, parsedUA)

	// Render the results
	return useragent.Results(result).Render(r.Context(), w)
}

// Helper functions

// getDeviceType returns a user-friendly device type description
func getDeviceType(ua ua.UserAgent) string {
	if ua.Bot {
		return "Bot/Crawler"
	}
	if ua.Mobile {
		return "Mobile"
	}
	if ua.Tablet {
		return "Tablet"
	}
	if ua.Desktop {
		return "Desktop"
	}
	return "Unknown"
}

// getBotName extracts bot name if it's a bot (simplified approach)
func getBotName(ua ua.UserAgent) string {
	if !ua.Bot {
		return ""
	}

	// Try to extract bot name from the user agent string
	lowerUA := strings.ToLower(ua.String)
	botNames := []string{"googlebot", "bingbot", "slurp", "duckduckbot", "baiduspider", "yandexbot", "facebookexternalhit", "twitterbot", "linkedinbot", "whatsapp", "telegrambot"}

	for _, bot := range botNames {
		if strings.Contains(lowerUA, bot) {
			return strings.Title(bot)
		}
	}

	return "Unknown Bot"
}

// formatBool converts a boolean to Yes/No string
func formatBool(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

// getValueOrEmpty returns the value if not empty, otherwise returns "Unknown"
func getValueOrEmpty(value string) string {
	if strings.TrimSpace(value) == "" {
		return "Unknown"
	}
	return value
}

// addBrowserDetection adds browser-specific detection flags
func addBrowserDetection(result *useragent.ParsingResult, ua ua.UserAgent) {
	browserName := strings.ToLower(ua.Name)

	// Check for specific browsers using string contains for reliability
	if strings.Contains(browserName, "chrome") && !strings.Contains(browserName, "edge") {
		result.FullDetails["Is Chrome"] = "Yes"
	}
	if strings.Contains(browserName, "firefox") {
		result.FullDetails["Is Firefox"] = "Yes"
	}
	if strings.Contains(browserName, "safari") && !strings.Contains(browserName, "chrome") {
		result.FullDetails["Is Safari"] = "Yes"
	}
	if strings.Contains(browserName, "opera") {
		result.FullDetails["Is Opera"] = "Yes"
	}
	if strings.Contains(browserName, "edge") {
		result.FullDetails["Is Edge"] = "Yes"
	}
	if strings.Contains(browserName, "internet explorer") || strings.Contains(browserName, "msie") {
		result.FullDetails["Is Internet Explorer"] = "Yes"
	}
	if strings.Contains(browserName, "brave") {
		result.FullDetails["Is Brave"] = "Yes"
	}
	if strings.Contains(browserName, "vivaldi") {
		result.FullDetails["Is Vivaldi"] = "Yes"
	}
}

// addOSDetection adds OS-specific detection flags
func addOSDetection(result *useragent.ParsingResult, ua ua.UserAgent) {
	osName := strings.ToLower(ua.OS)

	// Check for specific operating systems
	if strings.Contains(osName, "windows") {
		result.FullDetails["Is Windows"] = "Yes"

		// Add Windows version detection
		if strings.Contains(osName, "windows 10") || strings.Contains(osName, "nt 10.0") {
			result.FullDetails["Windows Version"] = "Windows 10/11"
		} else if strings.Contains(osName, "windows 8") || strings.Contains(osName, "nt 6.2") || strings.Contains(osName, "nt 6.3") {
			result.FullDetails["Windows Version"] = "Windows 8/8.1"
		} else if strings.Contains(osName, "windows 7") || strings.Contains(osName, "nt 6.1") {
			result.FullDetails["Windows Version"] = "Windows 7"
		}
	}

	if strings.Contains(osName, "mac") || strings.Contains(osName, "darwin") {
		result.FullDetails["Is macOS"] = "Yes"
	}

	if strings.Contains(osName, "linux") && !strings.Contains(osName, "android") {
		result.FullDetails["Is Linux"] = "Yes"
	}

	if strings.Contains(osName, "android") {
		result.FullDetails["Is Android"] = "Yes"
	}

	if strings.Contains(osName, "ios") || strings.Contains(osName, "iphone") || strings.Contains(osName, "ipad") {
		result.FullDetails["Is iOS"] = "Yes"
	}

	if strings.Contains(osName, "ubuntu") {
		result.FullDetails["Is Ubuntu"] = "Yes"
	}

	if strings.Contains(osName, "debian") {
		result.FullDetails["Is Debian"] = "Yes"
	}

	if strings.Contains(osName, "centos") {
		result.FullDetails["Is CentOS"] = "Yes"
	}
}

// addTechnicalDetails adds additional technical information
func addTechnicalDetails(result *useragent.ParsingResult, ua ua.UserAgent) {
	// Add rendering engine information based on common patterns
	uaString := strings.ToLower(ua.String)

	if strings.Contains(uaString, "webkit") {
		result.FullDetails["Rendering Engine"] = "WebKit"
	} else if strings.Contains(uaString, "gecko") && !strings.Contains(uaString, "like gecko") {
		result.FullDetails["Rendering Engine"] = "Gecko"
	} else if strings.Contains(uaString, "trident") {
		result.FullDetails["Rendering Engine"] = "Trident"
	} else if strings.Contains(uaString, "edge") {
		result.FullDetails["Rendering Engine"] = "EdgeHTML/Chromium"
	}

	// Add architecture information
	if strings.Contains(uaString, "x86_64") || strings.Contains(uaString, "win64") || strings.Contains(uaString, "wow64") {
		result.FullDetails["Architecture"] = "64-bit"
	} else if strings.Contains(uaString, "x86") || strings.Contains(uaString, "i386") || strings.Contains(uaString, "win32") {
		result.FullDetails["Architecture"] = "32-bit"
	} else if strings.Contains(uaString, "arm64") || strings.Contains(uaString, "aarch64") {
		result.FullDetails["Architecture"] = "ARM64"
	} else if strings.Contains(uaString, "arm") {
		result.FullDetails["Architecture"] = "ARM"
	}

	// Add mobile-specific information
	if ua.Mobile {
		if strings.Contains(uaString, "mobile safari") && !strings.Contains(uaString, "chrome") {
			result.FullDetails["Mobile Browser"] = "Mobile Safari"
		} else if strings.Contains(uaString, "chrome") && strings.Contains(uaString, "mobile") {
			result.FullDetails["Mobile Browser"] = "Chrome Mobile"
		} else if strings.Contains(uaString, "firefox") && strings.Contains(uaString, "mobile") {
			result.FullDetails["Mobile Browser"] = "Firefox Mobile"
		}

		// Add device model information for mobile devices
		if strings.Contains(uaString, "iphone") {
			result.FullDetails["Device Model"] = extractiPhoneModel(uaString)
		} else if strings.Contains(uaString, "ipad") {
			result.FullDetails["Device Model"] = "iPad"
		} else if strings.Contains(uaString, "android") {
			model := extractAndroidModel(uaString)
			if model != "" {
				result.FullDetails["Device Model"] = model
			}
		}
	}

	// Add security information
	if strings.Contains(uaString, "https") || strings.Contains(uaString, "secure") {
		result.FullDetails["Security"] = "HTTPS Capable"
	}

	// Add language information if available
	if strings.Contains(uaString, "en-us") {
		result.FullDetails["Language"] = "English (US)"
	} else if strings.Contains(uaString, "en-gb") {
		result.FullDetails["Language"] = "English (UK)"
	} else if strings.Contains(uaString, "en") {
		result.FullDetails["Language"] = "English"
	}
}

// extractiPhoneModel attempts to extract iPhone model from user agent
func extractiPhoneModel(uaString string) string {
	if strings.Contains(uaString, "iphone os 17") {
		return "iPhone (iOS 17)"
	} else if strings.Contains(uaString, "iphone os 16") {
		return "iPhone (iOS 16)"
	} else if strings.Contains(uaString, "iphone os 15") {
		return "iPhone (iOS 15)"
	} else if strings.Contains(uaString, "iphone os 14") {
		return "iPhone (iOS 14)"
	}
	return "iPhone"
}

// extractAndroidModel attempts to extract Android device model
func extractAndroidModel(uaString string) string {
	// Look for common Android device patterns
	if strings.Contains(uaString, "sm-") {
		// Samsung Galaxy pattern
		start := strings.Index(uaString, "sm-")
		if start != -1 {
			end := strings.Index(uaString[start:], ")")
			if end != -1 {
				return strings.ToUpper(uaString[start : start+end])
			}
		}
	}

	if strings.Contains(uaString, "pixel") {
		return "Google Pixel"
	}

	// Add more device patterns as needed
	return ""
}
