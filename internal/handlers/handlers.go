package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akozadaev/go_es_analytical_system/internal/models"
	"github.com/akozadaev/go_es_analytical_system/internal/storage"
	"github.com/gorilla/mux"
)

type Handlers struct {
	esStorage    *storage.ElasticsearchStorage
	pgStorage    *storage.PostgresStorage
}

func NewHandlers(esStorage *storage.ElasticsearchStorage, pgStorage *storage.PostgresStorage) *Handlers {
	return &Handlers{
		esStorage: esStorage,
		pgStorage: pgStorage,
	}
}

// RecommendLocations обрабатывает запрос на рекомендацию локаций
func (h *Handlers) RecommendLocations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RecommendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Region == "" || req.BusinessType == "" {
		http.Error(w, "Region and business_type are required", http.StatusBadRequest)
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	locations, err := h.esStorage.RecommendLocations(r.Context(), &req)
	if err != nil {
		log.Printf("Error recommending locations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Преобразуем указатели в значения для JSON
	locationValues := make([]models.Location, len(locations))
	for i, loc := range locations {
		locationValues[i] = *loc
	}

	response := models.RecommendResponse{
		Locations: locationValues,
		Total:     len(locationValues),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetLocation обрабатывает запрос на получение деталей локации
func (h *Handlers) GetLocation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Location ID is required", http.StatusBadRequest)
		return
	}

	location, err := h.esStorage.GetLocation(r.Context(), id)
	if err != nil {
		if err.Error() == "location not found" {
			http.Error(w, "Location not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting location: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(location); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetBusinessTypes обрабатывает запрос на получение списка типов бизнеса
func (h *Handlers) GetBusinessTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	businessTypes, err := h.pgStorage.GetBusinessTypes(r.Context())
	if err != nil {
		log.Printf("Error getting business types: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Преобразуем указатели в значения для JSON
	btValues := make([]models.BusinessType, len(businessTypes))
	for i, bt := range businessTypes {
		btValues[i] = *bt
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(btValues); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetRegions обрабатывает запрос на получение списка регионов
func (h *Handlers) GetRegions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	regions, err := h.pgStorage.GetRegions(r.Context())
	if err != nil {
		log.Printf("Error getting regions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Преобразуем указатели в значения для JSON
	regionValues := make([]models.Region, len(regions))
	for i, r := range regions {
		regionValues[i] = *r
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(regionValues); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HealthCheck обрабатывает запрос на проверку здоровья сервиса
func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

