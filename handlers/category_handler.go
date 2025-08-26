package handlers

import (
	"net/http"
	"strings"

	"github.com/Ndeta100/orbit2x/types"
	"github.com/Ndeta100/orbit2x/views/category"
	"github.com/Ndeta100/orbit2x/views/components"
	"github.com/Ndeta100/orbit2x/views/faq"
)

func FAQHandler(w http.ResponseWriter, r *http.Request) error {
	// Set security and SEO headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	component := faq.FAQPage()
	return component.Render(r.Context(), w)
}

// CategoriesOverviewHandler handles the main /categories page
func CategoriesOverviewHandler(w http.ResponseWriter, r *http.Request) error {
	// Set headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	component := category.CategoriesPage()
	return component.Render(r.Context(), w)
}

// CategoryDetailHandler handles individual category pages /categories/{slug}
func CategoryDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL path
	path := strings.TrimPrefix(r.URL.Path, "/categories/")
	if path == "" {
		http.Redirect(w, r, "/categories", http.StatusMovedPermanently)
		return
	}

	// Get category data by slug
	category := getCategoryBySlug(path)
	if category == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	component := components.CategoryPage(*category)
	component.Render(r.Context(), w)
}

// Helper function to get category by slug
func getCategoryBySlug(slug string) *types.CategoryData {
	switch slug {
	case "dns-tools":
		cat := types.GetDNSToolsCategory()
		return &cat
	case "developer-tools", "dev-tools":
		cat := types.GetDeveloperToolsCategory()
		return &cat
	case "designer-tools", "design-tools":
		cat := types.GetDesignerToolsCategory()
		return &cat
	case "webmaster-tools":
		cat := types.GetWebmasterToolsCategory()
		return &cat
	case "network-tools":
		cat := types.GetNetworkToolsCategory()
		return &cat
	case "security-tools", "cybersecurity-tools":
		cat := types.GetSecurityToolsCategory()
		return &cat
	case "productivity-tools":
		cat := types.GetProductivityToolsCategory()
		return &cat
	case "gaming-tools":
		cat := types.GetGamingToolsCategory()
		return &cat
	case "more-tools", "more-categories":
		cat := types.GetMoreCategoriesCategory()
		return &cat
	default:
		return nil
	}
}

// HandleDNSToolsCategory handles /dns-tools route
func HandleDNSToolsCategory(w http.ResponseWriter, r *http.Request) error {
	category_1 := types.GetDNSToolsCategory()
	return components.CategoryPage(category_1).Render(r.Context(), w)
}

// HandleDeveloperToolsCategory handles /dev-tools route
func HandleDeveloperToolsCategory(w http.ResponseWriter, r *http.Request) error {
	categoryTool := types.GetDeveloperToolsCategory()
	return components.CategoryPage(categoryTool).Render(r.Context(), w)
}

// HandleDesignerToolsCategory handles /design-tools route
func HandleDesignerToolsCategory(w http.ResponseWriter, r *http.Request) error {
	category := types.GetDesignerToolsCategory()
	return components.CategoryPage(category).Render(r.Context(), w)
}

// HandleWebmasterToolsCategory handles /webmaster-tools route
func HandleWebmasterToolsCategory(w http.ResponseWriter, r *http.Request) error {
	category := types.GetWebmasterToolsCategory()
	return components.CategoryPage(category).Render(r.Context(), w)
}

// HandleNetworkToolsCategory handles /network-tools route
func HandleNetworkToolsCategory(w http.ResponseWriter, r *http.Request) error {
	category := types.GetNetworkToolsCategory()
	return components.CategoryPage(category).Render(r.Context(), w)
}

// HandleSecurityToolsCategory handles /security-tools route
func HandleSecurityToolsCategory(w http.ResponseWriter, r *http.Request) error {
	category := types.GetSecurityToolsCategory()
	return components.CategoryPage(category).Render(r.Context(), w)
}

// HandleProductivityToolsCategory handles /productivity-tools route
func HandleProductivityToolsCategory(w http.ResponseWriter, r *http.Request) error {
	category := types.GetProductivityToolsCategory()
	return components.CategoryPage(category).Render(r.Context(), w)
}

// HandleGamingToolsCategory handles /gaming-tools route
func HandleGamingToolsCategory(w http.ResponseWriter, r *http.Request) error {
	category := types.GetGamingToolsCategory()
	return components.CategoryPage(category).Render(r.Context(), w)
}

// HandleMoreCategoriesCategory handles /more-tools route
func HandleMoreCategoriesCategory(w http.ResponseWriter, r *http.Request) error {
	category := types.GetMoreCategoriesCategory()
	return components.CategoryPage(category).Render(r.Context(), w)
}
