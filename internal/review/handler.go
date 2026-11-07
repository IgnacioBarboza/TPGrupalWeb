package Review

import (
	dto "WebGrupal/internal"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// Handler
type ReviewHandler struct {
	service *ReviewService
}

// Constructor
func NewReviewHandler(service *ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

// POST /reviews: recibe el formulario, crea el registro y redirige a "/" (PRG).
func (h *ReviewHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	filmID, err := strconv.ParseInt(r.FormValue("reviewfilmid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.FormValue("reviewuserid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	stars, err := strconv.ParseInt(r.FormValue("reviewstars"), 10, 16)
	if err != nil {
		http.Error(w, "Estrellas inválidas", http.StatusBadRequest)
		return
	}

	review := dto.Review{
		Reviewfilmid:      int32(filmID),
		Reviewuserid:      int32(userID),
		Reviewstars:       int16(stars),
		Reviewdescription: r.FormValue("reviewdescription"),
	}

	_, err = h.service.Create(r.Context(), review)
	if err != nil {
		h.writeError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *ReviewHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidReviewStars),
		errors.Is(err, ErrInvalidReviewDescriptionEmpty):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}
