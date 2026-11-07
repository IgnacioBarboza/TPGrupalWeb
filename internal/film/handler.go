package Film

import (
	dto "WebGrupal/internal"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Handler
type FilmHandler struct {
	service *FilmService
}

// Constructor
func NewFilmHandler(service *FilmService) *FilmHandler {
	return &FilmHandler{service: service}
}

// POST /films: recibe el formulario, crea el registro y redirige a "/" (PRG).
func (h *FilmHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	releaseDate, err := time.Parse("2006-01-02", r.FormValue("filmreleasedate"))
	if err != nil {
		http.Error(w, "Fecha de estreno inválida", http.StatusBadRequest)
		return
	}

	duration, err := strconv.ParseInt(r.FormValue("filmduration"), 10, 32)
	if err != nil {
		http.Error(w, "Duración inválida", http.StatusBadRequest)
		return
	}

	film := dto.Film{
		Filmtitle:       r.FormValue("filmtitle"),
		Filmreleasedate: releaseDate,
		Filmduration:    int32(duration),
	}

	_, err = h.service.Create(r.Context(), film)
	if err != nil {
		h.writeError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *FilmHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidFilmTitleEmpty),
		errors.Is(err, ErrInvalidFilmDuration):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}
