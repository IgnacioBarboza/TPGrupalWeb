package Review

import (
	dto "WebGrupal/internal"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Handler
type ReviewHandler struct {
	service *ReviewService
}

// Constructor
func NewReviewHandler(service *ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/reviews", h.handleCollection)
	mux.HandleFunc("/reviews/", h.handleReview)
}

// writeJSON para escritura de JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error al codificar respuesta: %v", err)
	}
}

// Casos de uso: GET /reviews/{filmID}/{userID}, PUT /reviews/{filmID}/{userID}, DELETE /reviews/{filmID}/{userID}
func (h *ReviewHandler) handleReview(w http.ResponseWriter, r *http.Request) {
	var filmID int32
	var userID int32
	if r.Method != http.MethodPost {
		var err error
		filmID, userID, err = extractIDs(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	switch r.Method {
	case http.MethodPost:
		h.createReview(w, r)
	case http.MethodGet:
		h.getReview(w, r, filmID, userID)
	case http.MethodPut:
		h.updateReview(w, r, filmID, userID)
	case http.MethodDelete:
		h.deleteReview(w, r, filmID, userID)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// POST /reviews  y  GET /reviews
func (h *ReviewHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createMultipleReviews(w, r)
	case http.MethodGet:
		h.getReviews(w, r)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// extractIDs obtiene filmID y userID de una ruta del tipo /reviews/{filmID}/{userID}
func extractIDs(path string) (int32, int32, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Obtener las tres partes de la ruta: "reviews", "filmID" y "userID"
	if len(parts) != 3 || parts[0] != "reviews" {
		return 0, 0, errors.New("ruta inválida, se esperaba /reviews/{filmID}/{userID}")
	}

	filmID, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil || filmID <= 0 {
		return 0, 0, errors.New("filmID inválido")
	}

	userID, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil || userID <= 0 {
		return 0, 0, errors.New("userID inválido")
	}

	return int32(filmID), int32(userID), nil
}

// GET /reviews/{filmID}/{userID}: Obtiene el ID de la URL, busca el registro y lo devuelve. Si no existe, devuelve un 404.
func (h *ReviewHandler) getReview(w http.ResponseWriter, r *http.Request, filmID int32, userID int32) {
	review, err := h.service.GetReview(r.Context(), filmID, userID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, review)
}

// PUT /reviews/{filmID}/{userID}: Recibe datos JSON, los valida, y actualiza el registro correspondiente.
func (h *ReviewHandler) updateReview(w http.ResponseWriter, r *http.Request, filmID int32, userID int32) {
	var receivedReview dto.Review
	if err := json.NewDecoder(r.Body).Decode(&receivedReview); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if filmID != receivedReview.Reviewfilmid || userID != receivedReview.Reviewuserid {
		http.Error(w, "Los IDs de la URL no coinciden con los IDs en el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	updatedReview, err := h.service.Update(r.Context(), receivedReview)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updatedReview)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *ReviewHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrReviewNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidReviewStars),
		errors.Is(err, ErrInvalidReviewDescriptionEmpty),
		errors.Is(err, ErrInvalidReview):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}

// GET /reviews: Usa GetReviews para obtener todos los registros y devolverlos como un array JSON.
func (h *ReviewHandler) getReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.service.GetReviews(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, reviews)
}

// POST /reviews: Recibe datos en formato JSON (array), los valida y crea múltiples reviews.
func (h *ReviewHandler) createMultipleReviews(w http.ResponseWriter, r *http.Request) {
	var reviews []dto.Review
	if err := json.NewDecoder(r.Body).Decode(&reviews); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateMultipleReviews(r.Context(), reviews)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created) // 201 (statusCreated)
}

// DELETE /reviews/{filmID}/{userID}: Elimina el registro y devuelve un estado 204.
func (h *ReviewHandler) deleteReview(w http.ResponseWriter, r *http.Request, filmID int32, userID int32) {
	err := h.service.Delete(r.Context(), filmID, userID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ReviewHandler) createReview(w http.ResponseWriter, r *http.Request) {
	var receivedReview dto.Review
	if err := json.NewDecoder(r.Body).Decode(&receivedReview); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	createdReview, err := h.service.Create(r.Context(), receivedReview)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createdReview) // 201 (statusCreated)
}
