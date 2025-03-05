package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *APIImpl) respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		fmt.Printf("Could not encode JSON body: %s\n", err.Error())
	}
}

func (a *APIImpl) respondError(w http.ResponseWriter, status int, err error, msg string) {
	type response struct {
		Error string `json:"error"`
	}
	a.respond(w, status, response{Error: msg})
}

func parseRequest(r *http.Request, dst interface{}) error {
	if err := decoder.Decode(dst, r.URL.Query()); err != nil {
		return fmt.Errorf("failed to decode query parameters: %w", err)
	}

	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form body: %w", err)
	}
	if err := decoder.Decode(dst, r.Form); err != nil {
		return fmt.Errorf("failed to decode form body: %w", err)
	}

	return nil
}
