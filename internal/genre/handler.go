package Genre

import (
	dto "WebGrupal/internal"
	"errors"
	"log"
	"net/http"
)

// Handler
type GenreHandler struct {
	service *GenreService
}

// Constructor
func NewGenreHandler(service *GenreService) *GenreHandler {
	return &GenreHandler{service: service}
}

// POST /genres: recibe el formulario, crea el registro y redirige a "/" (PRG).
func (h *GenreHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	genre := dto.Genre{
		Genrename: r.FormValue("genrename"),
	}

	_, err = h.service.Create(r.Context(), genre)
	if err != nil {
		h.writeError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *GenreHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidGenreNameEmpty):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}
