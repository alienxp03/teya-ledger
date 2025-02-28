package api

import (
	"log/slog"
	"net/http"
	"sync"
)

// API provides the REST endpoints for the application.
type API struct {
	Logger *slog.Logger

	once sync.Once
	mux  *http.ServeMux
}

func (a *API) setupRoutes() {
	mux := http.NewServeMux()

	a.mux = mux
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.once.Do(a.setupRoutes)
	a.Logger.Info("Request received", "method", r.Method, "path", r.URL.Path)
	a.mux.ServeHTTP(w, r)
}
