package handlers

import (
	"net/http"

	"github.com/Ndeta100/orbit2x/types"
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
