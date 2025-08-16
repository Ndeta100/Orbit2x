// handlers/color_handler.go
package handlers

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Ndeta100/orbit2x/views/colorconverter" // Adjust to your actual path
)

// HandleColorIndex renders the Color Converter page
func HandleColorIndex(w http.ResponseWriter, r *http.Request) error {
	return colorconverter.Index().Render(r.Context(), w)
}

// HandleColorConvert converts color formats
func HandleColorConvert(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return colorconverter.Results(colorconverter.ConversionResult{
			Error: "Failed to parse form data",
		}).Render(context.Background(), w)
	}

	// Get input and format
	colorInput := r.FormValue("color")
	fromFormat := r.FormValue("from")

	if colorInput == "" {
		return colorconverter.Results(colorconverter.ConversionResult{
			Error: "Color input is required",
		}).Render(context.Background(), w)
	}

	// Prepare result
	result := colorconverter.ConversionResult{
		InputColor: colorInput,
		InputType:  fromFormat,
	}

	// Process based on input format
	switch fromFormat {
	case "hex":
		// Parse HEX
		red, green, blue, alpha, err := parseHex(colorInput)
		if err != nil {
			result.Error = err.Error()
			return colorconverter.Results(result).Render(context.Background(), w)
		}

		// Convert to other formats
		result.HexColor = formatHex(red, green, blue, alpha)
		result.RGBColor = formatRGB(red, green, blue, alpha)
		result.HSLColor = formatHSL(red, green, blue, alpha)
		result.PreviewColor = result.HexColor

	case "rgb":
		// Parse RGB
		red, green, blue, alpha, err := parseRGB(colorInput)
		if err != nil {
			result.Error = err.Error()
			return colorconverter.Results(result).Render(context.Background(), w)
		}

		// Convert to other formats
		result.HexColor = formatHex(red, green, blue, alpha)
		result.RGBColor = formatRGB(red, green, blue, alpha)
		result.HSLColor = formatHSL(red, green, blue, alpha)
		result.PreviewColor = result.HexColor

	case "hsl":
		// Parse HSL
		h, s, l, alpha, err := parseHSL(colorInput)
		if err != nil {
			result.Error = err.Error()
			return colorconverter.Results(result).Render(context.Background(), w)
		}

		// Convert HSL to RGB
		red, green, blue := hslToRGB(h, s, l)

		// Convert to other formats
		result.HexColor = formatHex(red, green, blue, alpha)
		result.RGBColor = formatRGB(red, green, blue, alpha)
		result.HSLColor = formatHSL(red, green, blue, alpha)
		result.PreviewColor = result.HexColor

	default:
		result.Error = "Unsupported input format"
		return colorconverter.Results(result).Render(context.Background(), w)
	}

	// Render the result
	return colorconverter.Results(result).Render(context.Background(), w)
}

// parseHex parses a hexadecimal color string
func parseHex(hexColor string) (r, g, b, a int, err error) {
	// Remove # if present
	hexColor = strings.TrimPrefix(hexColor, "#")
	hexColor = strings.ToLower(hexColor)

	// Check for valid hex format
	hexPattern := regexp.MustCompile(`^([0-9a-f]{3}|[0-9a-f]{4}|[0-9a-f]{6}|[0-9a-f]{8})$`)
	if !hexPattern.MatchString(hexColor) {
		return 0, 0, 0, 255, fmt.Errorf("invalid hex color format: %s", hexColor)
	}

	// Parse based on length
	switch len(hexColor) {
	case 3: // #RGB
		rVal, _ := strconv.ParseInt(string(hexColor[0])+string(hexColor[0]), 16, 0)
		gVal, _ := strconv.ParseInt(string(hexColor[1])+string(hexColor[1]), 16, 0)
		bVal, _ := strconv.ParseInt(string(hexColor[2])+string(hexColor[2]), 16, 0)
		r, g, b, a = int(rVal), int(gVal), int(bVal), 255
	case 4: // #RGBA
		rVal, _ := strconv.ParseInt(string(hexColor[0])+string(hexColor[0]), 16, 0)
		gVal, _ := strconv.ParseInt(string(hexColor[1])+string(hexColor[1]), 16, 0)
		bVal, _ := strconv.ParseInt(string(hexColor[2])+string(hexColor[2]), 16, 0)
		aVal, _ := strconv.ParseInt(string(hexColor[3])+string(hexColor[3]), 16, 0)
		r, g, b, a = int(rVal), int(gVal), int(bVal), int(aVal)
	case 6: // #RRGGBB
		rVal, _ := strconv.ParseInt(hexColor[0:2], 16, 0)
		gVal, _ := strconv.ParseInt(hexColor[2:4], 16, 0)
		bVal, _ := strconv.ParseInt(hexColor[4:6], 16, 0)
		r, g, b, a = int(rVal), int(gVal), int(bVal), 255
	case 8: // #RRGGBBAA
		rVal, _ := strconv.ParseInt(hexColor[0:2], 16, 0)
		gVal, _ := strconv.ParseInt(hexColor[2:4], 16, 0)
		bVal, _ := strconv.ParseInt(hexColor[4:6], 16, 0)
		aVal, _ := strconv.ParseInt(hexColor[6:8], 16, 0)
		r, g, b, a = int(rVal), int(gVal), int(bVal), int(aVal)
	}

	return r, g, b, a, nil
}

// parseRGB parses an RGB/RGBA color string
func parseRGB(rgbColor string) (r, g, b, a int, err error) {
	// Default alpha
	a = 255

	// Check for RGB format
	rgbPattern := regexp.MustCompile(`rgba?\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)(?:\s*,\s*([0-9.]+))?\s*\)`)
	matches := rgbPattern.FindStringSubmatch(rgbColor)

	if matches == nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid RGB format: %s", rgbColor)
	}

	// Parse RGB values
	r, _ = strconv.Atoi(matches[1])
	g, _ = strconv.Atoi(matches[2])
	b, _ = strconv.Atoi(matches[3])

	// Check value ranges
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return 0, 0, 0, 0, fmt.Errorf("RGB values must be between 0 and 255")
	}

	// Parse alpha if present
	if len(matches) > 4 && matches[4] != "" {
		alpha, err := strconv.ParseFloat(matches[4], 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("invalid alpha value")
		}
		if alpha < 0 || alpha > 1 {
			return 0, 0, 0, 0, fmt.Errorf("alpha value must be between 0 and 1")
		}
		a = int(alpha * 255)
	}

	return r, g, b, a, nil
}

// parseHSL parses an HSL/HSLA color string
func parseHSL(hslColor string) (h, s, l float64, a int, err error) {
	// Default alpha
	a = 255

	// Check for HSL format
	hslPattern := regexp.MustCompile(`hsla?\(\s*(\d+)\s*,\s*(\d+)%\s*,\s*(\d+)%(?:\s*,\s*([0-9.]+))?\s*\)`)
	matches := hslPattern.FindStringSubmatch(hslColor)

	if matches == nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid HSL format: %s", hslColor)
	}

	// Parse HSL values
	hVal, _ := strconv.Atoi(matches[1])
	sVal, _ := strconv.Atoi(matches[2])
	lVal, _ := strconv.Atoi(matches[3])

	// Check value ranges
	if hVal < 0 || hVal > 360 {
		return 0, 0, 0, 0, fmt.Errorf("hue must be between 0 and 360")
	}
	if sVal < 0 || sVal > 100 {
		return 0, 0, 0, 0, fmt.Errorf("saturation must be between 0%% and 100%%")
	}
	if lVal < 0 || lVal > 100 {
		return 0, 0, 0, 0, fmt.Errorf("lightness must be between 0%% and 100%%")
	}

	// Convert to decimal
	h = float64(hVal)
	s = float64(sVal) / 100
	l = float64(lVal) / 100

	// Parse alpha if present
	if len(matches) > 4 && matches[4] != "" {
		alpha, err := strconv.ParseFloat(matches[4], 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("invalid alpha value")
		}
		if alpha < 0 || alpha > 1 {
			return 0, 0, 0, 0, fmt.Errorf("alpha value must be between 0 and 1")
		}
		a = int(alpha * 255)
	}

	return h, s, l, a, nil
}

// formatHex formats RGB values to HEX
func formatHex(r, g, b, a int) string {
	if a == 255 {
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", r, g, b, a)
}

// formatRGB formats RGB values to RGB/RGBA string
func formatRGB(r, g, b, a int) string {
	if a == 255 {
		return fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)
	}
	return fmt.Sprintf("rgba(%d, %d, %d, %.2f)", r, g, b, float64(a)/255)
}

// formatHSL formats RGB values to HSL/HSLA string
func formatHSL(r, g, b, a int) string {
	// Convert RGB to HSL
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	max := math.Max(math.Max(rf, gf), bf)
	min := math.Min(math.Min(rf, gf), bf)
	l := (max + min) / 2
	delta := max - min

	var h, s float64
	if delta == 0 {
		h = 0
		s = 0
	} else {
		s = delta / (1 - math.Abs(2*l-1))

		switch max {
		case rf:
			h = 60 * (((gf - bf) / delta) + float64(0))
			if gf < bf {
				h += 360
			}
		case gf:
			h = 60 * (((bf - rf) / delta) + float64(2))
		case bf:
			h = 60 * (((rf - gf) / delta) + float64(4))
		}
	}

	// Round values
	h = math.Round(h)
	s = math.Round(s * 100)
	l = math.Round(l * 100)

	if a == 255 {
		return fmt.Sprintf("hsl(%d, %d%%, %d%%)", int(h), int(s), int(l))
	}
	return fmt.Sprintf("hsla(%d, %d%%, %d%%, %.2f)", int(h), int(s), int(l), float64(a)/255)
}

// hslToRGB converts HSL to RGB
func hslToRGB(h, s, l float64) (r, g, b int) {
	var c = (1 - math.Abs(2*l-1)) * s
	var x = c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	var m = l - c/2
	var rf, gf, bf float64

	switch {
	case h < 60:
		rf, gf, bf = c, x, 0
	case h < 120:
		rf, gf, bf = x, c, 0
	case h < 180:
		rf, gf, bf = 0, c, x
	case h < 240:
		rf, gf, bf = 0, x, c
	case h < 300:
		rf, gf, bf = x, 0, c
	default:
		rf, gf, bf = c, 0, x
	}

	r = int(math.Round((rf + m) * 255))
	g = int(math.Round((gf + m) * 255))
	b = int(math.Round((bf + m) * 255))

	return r, g, b
}

// HandleRandomColor generates a random color
func HandleRandomColor(w http.ResponseWriter, r *http.Request) error {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Generate random RGB values
	red := rand.Intn(256)
	green := rand.Intn(256)
	blue := rand.Intn(256)

	// Create result
	result := colorconverter.ConversionResult{
		InputColor:   fmt.Sprintf("rgb(%d, %d, %d)", red, green, blue),
		InputType:    "rgb",
		HexColor:     formatHex(red, green, blue, 255),
		RGBColor:     formatRGB(red, green, blue, 255),
		HSLColor:     formatHSL(red, green, blue, 255),
		PreviewColor: formatHex(red, green, blue, 255),
	}

	// Render the result
	return colorconverter.Results(result).Render(context.Background(), w)
}
