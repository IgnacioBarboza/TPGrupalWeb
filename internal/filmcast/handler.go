package FilmCast

import (
	dto "WebGrupal/internal"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// Handler
type FilmCastHandler struct {
	service *FilmCastService
}

// Constructor
func NewFilmCastHandler(service *FilmCastService) *FilmCastHandler {
	return &FilmCastHandler{service: service}
}

// POST /filmcasts: recibe el formulario, crea el registro y redirige a "/" (PRG).
func (h *FilmCastHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	filmID, err := strconv.ParseInt(r.FormValue("filmid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	castMemberID, err := strconv.ParseInt(r.FormValue("castmemberid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de miembro del elenco inválido", http.StatusBadRequest)
		return
	}

	filmCast := dto.Filmcast{
		Filmid:       int32(filmID),
		Castmemberid: int32(castMemberID),
	}

	_, err = h.service.Create(r.Context(), filmCast)
	if err != nil {
		h.writeError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *FilmCastHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidFilmCastIDs):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}
