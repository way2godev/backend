package cache

import (
	"encoding/json"
	"net/http"
	"way2go/jobs"
)

func handleUpdateStopsCache(w http.ResponseWriter, r *http.Request) {
	jobs.ForceUpdateStopsCache <- true
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Update cache job triggered"})
}
