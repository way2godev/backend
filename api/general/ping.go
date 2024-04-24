package general

import (
	"encoding/json"
	"net/http"
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}