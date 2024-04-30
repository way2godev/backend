package stops

import (
	"encoding/json"
	"net/http"
	"strconv"
	"way2go/domain/entities"
)

func index(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	stops, total, err := entities.GetStops(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"stops": stops,
		"pagination": map[string]interface{}{
			"page": page,
			"total": total,
		},	
	}
	json.NewEncoder(w).Encode(response)
}