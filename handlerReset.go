package main

import "net/http"

func (cfg *apiConfig) handleMetricsReset(w http.ResponseWriter, r *http.Request) {
	// Resets the access count and deletes all users. Only authorized on dev platforms.
	if cfg.platform == "dev" {
		cfg.fileserverHits.Store(0)
		cfg.db.ResetUser(r.Context())
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(http.StatusText(http.StatusForbidden)))
}
