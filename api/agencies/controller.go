package agencies

import (
	"encoding/json"
	"net/http"
	"way2go/domain/entities"
	"way2go/infraestructure/database"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()

	var agencies []entities.Agency
	db.Select("id, name, raw_id, url, phone_number").Find(&agencies)

	var res []map[string]interface{}
	for _, agency := range agencies {
		res = append(res, map[string]interface{}{
			"id":          agency.ID,
			"name":        agency.Name,
			"raw_id":      agency.RawId,
			"url":         agency.Url,
			"phone_number": agency.PhoneNumber,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

type createAgencyRequest struct {
	Name        string  `json:"name"`
	RawId       *string `json:"raw_id"`
	Url         *string `json:"url"`
	PhoneNumber *string `json:"phone_number"`
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var req createAgencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid request"})
		return
	}

	agency := entities.Agency{
		Name:        req.Name,
		RawId:       req.RawId,
		Url:         req.Url,
		PhoneNumber: req.PhoneNumber,
	}

	db := database.GetDB()
	var existingAgency entities.Agency
	db.Where("name = ?", agency.Name).First(&existingAgency)
	if existingAgency.ID != 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "agency already exists"})
		return
	}
	db.Create(&agency)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "created"})
}
