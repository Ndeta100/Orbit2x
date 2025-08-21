package handlers

import (
	"net/http"

	"github.com/Ndeta100/orbit2x/types"
	"github.com/Ndeta100/orbit2x/views/components"
)

// HandleDNSToolsCategory handles /dns-tools route
func HandleDNSToolsCategory(w http.ResponseWriter, r *http.Request) error {
	category := types.GetDNSToolsCategory()
	return components.CategoryPage(category).Render(r.Context(), w)
}

// HandleDeveloperToolsCategory handles /dev-tools route
func HandleDeveloperToolsCategory(w http.ResponseWriter, r *http.Request) error {
	category := types.GetDeveloperToolsCategory()
	return components.CategoryPage(category).Render(r.Context(), w)
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
