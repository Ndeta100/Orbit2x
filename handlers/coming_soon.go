package handlers

import (
	"net/http"

	"github.com/Ndeta100/orbit2x/views/pages"
)

func HandleComingSoon(w http.ResponseWriter, r *http.Request) error {
	return pages.ComingSoon().Render(r.Context(), w)
}

// 404 fallback handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	pages.ComingSoon().Render(r.Context(), w)
}
