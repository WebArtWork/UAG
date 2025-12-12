package docs

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterOpenAPIService adds a minimal OpenAPI endpoint to the API router.
// This is a placeholder implementation to satisfy the application's runtime
// dependency on an OpenAPI registration hook. The endpoint currently returns a
// 404 response, and can be replaced with generated documentation when
// available.
func RegisterOpenAPIService(_ string, router *mux.Router) {
	if router == nil {
		return
	}

	router.HandleFunc("/openapi", func(w http.ResponseWriter, _ *http.Request) {
		http.NotFound(w, nil)
	}).Methods(http.MethodGet)
}
