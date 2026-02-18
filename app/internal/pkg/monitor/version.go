package monitor

import (
	"encoding/json"
	"net/http"
)

type VersionResponse struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

func Version(service, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(VersionResponse{
			Service: service,
			Version: version,
		})
	}
}
