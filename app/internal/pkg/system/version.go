package system

import (
	"encoding/json"
	"net/http"
)

type versionResponse struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

func Version(service, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(versionResponse{
			Service: service,
			Version: version,
		})
	}
}
