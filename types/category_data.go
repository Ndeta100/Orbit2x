package types

// GetDNSToolsCategory returns DNS & IP Tools category data
func GetDNSToolsCategory() CategoryData {
	return CategoryData{
		Name:        "DNS & IP Tools",
		Description: "Domain lookups, IP analysis, network diagnostics and DNS record management tools",
		Icon:        "M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9",
		ToolCount:   "5 tools available",
		SearchHint:  "DNS tools",
		Color:       "blue",
		Tools: []CategoryTool{
			{
				Name:        "DNS Lookup",
				Description: "Comprehensive DNS record lookup for any domain including A, AAAA, MX, NS, TXT, SOA and WHOIS information",
				URL:         "/lookup",
				Icon:        "M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9",
				Tags:        []string{"DNS", "Domain", "WHOIS"},
				IsPopular:   true,
			},
			{
				Name:        "My IP Address",
				Description: "View your current public IP address and detailed connection information",
				URL:         "/myip",
				Icon:        "M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7",
				Tags:        []string{"IP", "Network", "Location"},
				IsPopular:   true,
			},
			// Add more tools as you build them
		},
	}
}

// GetDeveloperToolsCategory returns Developer Tools category data
func GetDeveloperToolsCategory() CategoryData {
	return CategoryData{
		Name:        "Developer Tools",
		Description: "Code utilities, formatters, validators, generators and development aids",
		Icon:        "M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4",
		ToolCount:   "8 tools available",
		SearchHint:  "development tools",
		Color:       "green",
		Tools: []CategoryTool{
			{
				Name:        "JSON Formatter",
				Description: "Format, validate and beautify JSON data with syntax highlighting",
				URL:         "/formatter",
				Icon:        "M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4",
				Tags:        []string{"JSON", "Format", "Validate"},
				IsPopular:   true,
			},
			{
				Name:        "Text Encoder",
				Description: "Encode and decode text in various formats including Base64, URL encoding and more",
				URL:         "/encoder",
				Icon:        "M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z",
				Tags:        []string{"Encode", "Base64", "URL"},
			},
			{
				Name:        "File Converter",
				Description: "Convert between different file formats including CSV to JSON and vice versa",
				URL:         "/converter",
				Icon:        "M8 7H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h6a2 2 0 002-2V9a2 2 0 00-2-2h-6a2 2 0 00-2 2v10a2 2 0 002 2zm8-12V7a2 2 0 00-2-2h-2a2 2 0 00-2 2v8a2 2 0 002 2h2a2 2 0 002-2z",
				Tags:        []string{"Convert", "CSV", "JSON"},
			},
			// Add more tools as you build them
		},
	}
}

// GetDesignerToolsCategory returns Designer Tools category data (empty for now)
func GetDesignerToolsCategory() CategoryData {
	return CategoryData{
		Name:        "Designer Tools",
		Description: "Color palettes, gradients, typography tools and design utilities",
		Icon:        "M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4z",
		ToolCount:   "2 tools available",
		SearchHint:  "design tools",
		Color:       "purple",
		Tools: []CategoryTool{
			{
				Name:        "Color Converter",
				Description: "Convert colors between different formats (HEX, RGB, HSL) and generate palettes",
				URL:         "/color",
				Icon:        "M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4z",
				Tags:        []string{"Color", "HEX", "RGB"},
				IsNew:       true,
			},
		},
	}
}

// GetWebmasterToolsCategory returns Webmaster Tools category data (empty for now)
func GetWebmasterToolsCategory() CategoryData {
	return CategoryData{
		Name:        "Webmaster Tools",
		Description: "SEO analysis, site optimization, performance monitoring and webmaster utilities",
		Icon:        "M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z",
		ToolCount:   "Tools coming soon",
		SearchHint:  "SEO and webmaster tools",
		Color:       "orange",
		Tools:       []CategoryTool{}, // Empty - will show coming soon
	}
}

// GetNetworkToolsCategory returns Network Tools category data (empty for now)
func GetNetworkToolsCategory() CategoryData {
	return CategoryData{
		Name:        "Network Tools",
		Description: "Ping, traceroute, port scanning, bandwidth testing and network diagnostics",
		Icon:        "M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h6a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h6a2 2 0 002-2v-4a2 2 0 00-2-2m8 0V9a2 2 0 012-2h2a2 2 0 012 2v3m0 0v6a2 2 0 01-2 2h-2a2 2 0 01-2-2v-6m0 0h4",
		ToolCount:   "Tools coming soon",
		SearchHint:  "network diagnostic tools",
		Color:       "indigo",
		Tools:       []CategoryTool{}, // Empty - will show coming soon
	}
}

// GetSecurityToolsCategory returns Cybersecurity Tools category data
func GetSecurityToolsCategory() CategoryData {
	return CategoryData{
		Name:        "Cybersecurity Tools",
		Description: "Vulnerability scanning, password testing, encryption and security analysis tools",
		Icon:        "M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z",
		ToolCount:   "2 tools available",
		SearchHint:  "security tools",
		Color:       "red",
		Tools: []CategoryTool{
			{
				Name:        "SSL Checker",
				Description: "Verify SSL certificates, check expiration dates and analyze security configurations",
				URL:         "/ssl",
				Icon:        "M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z",
				Tags:        []string{"SSL", "Certificate", "Security"},
				IsPopular:   true,
			},
			{
				Name:        "Hash Generator",
				Description: "Generate MD5, SHA1, SHA256 and other hash functions for text and files",
				URL:         "/hash",
				Icon:        "M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z",
				Tags:        []string{"Hash", "MD5", "SHA256"},
			},
		},
	}
}

// GetProductivityToolsCategory returns Productivity Tools category data
func GetProductivityToolsCategory() CategoryData {
	return CategoryData{
		Name:        "Productivity Tools",
		Description: "Text processing, file converters, calculators and productivity utilities",
		Icon:        "M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2",
		ToolCount:   "3 tools available",
		SearchHint:  "productivity tools",
		Color:       "teal",
		Tools: []CategoryTool{
			{
				Name:        "User Agent Parser",
				Description: "Parse and analyze user agent strings to identify browsers, devices and operating systems",
				URL:         "/useragent",
				Icon:        "M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2",
				Tags:        []string{"User Agent", "Browser", "Device"},
			},
			{
				Name:        "Image to Base64",
				Description: "Convert images to Base64 encoding for embedding in HTML, CSS or applications",
				URL:         "/imagebase64",
				Icon:        "M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z",
				Tags:        []string{"Image", "Base64", "Convert"},
			},
			{
				Name:        "Subnet Calculator",
				Description: "Calculate network addresses, subnet masks and CIDR notation for network planning",
				URL:         "/subnet",
				Icon:        "M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z",
				Tags:        []string{"Network", "Subnet", "CIDR"},
			},
		},
	}
}

// GetGamingToolsCategory returns Gaming Tools category data (empty for now)
func GetGamingToolsCategory() CategoryData {
	return CategoryData{
		Name:        "Gaming Tools",
		Description: "Game utilities, generators, calculators and tools for gamers",
		Icon:        "M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 011-1h1a2 2 0 100-4H7a1 1 0 01-1-1V7a1 1 0 011-1h3a1 1 0 001-1V4z",
		ToolCount:   "Tools coming soon",
		SearchHint:  "gaming tools",
		Color:       "pink",
		Tools:       []CategoryTool{}, // Empty - will show coming soon
	}
}

// GetMoreCategoriesCategory returns More Categories data (empty for now)
func GetMoreCategoriesCategory() CategoryData {
	return CategoryData{
		Name:        "More Categories",
		Description: "Finance, health, education and other specialized tools",
		Icon:        "M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10",
		ToolCount:   "Categories expanding",
		SearchHint:  "specialized tools",
		Color:       "gray",
		Tools:       []CategoryTool{}, // Empty - will show coming soon
	}
}
