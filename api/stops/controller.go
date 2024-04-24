package stops

import (
	"encoding/json"
	"net/http"
	"way2go/domain/services"
)

func handleSearch(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("s")
	if search == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing search query"})
		return
	}
	results := services.StopSearchServiceInstance.Search(search)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
