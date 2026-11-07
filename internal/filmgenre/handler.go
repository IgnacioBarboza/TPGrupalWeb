package FilmGenre

import (
	dto "WebGrupal/internal"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// Handler
type FilmGenreHandler struct {
	service *FilmGenreService
}

// Constructor
func NewFilmGenreHandler(service *FilmGenreService) *FilmGenreHandler {
	return &FilmGenreHandler{service: service}
}

// POST /filmgenres: recibe el formulario, crea el registro y redirige a "/" (PRG).
func (h *FilmGenreHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	genreID, err := strconv.ParseInt(r.FormValue("genre_genreid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de género inválido", http.StatusBadRequest)
		return
	}

	filmID, err := strconv.ParseInt(r.FormValue("film_filmid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	filmGenre := dto.Filmgenre{
		GenreGenreid: int32(genreID),
		FilmFilmid:   int32(filmID),
	}

	_, err = h.service.Create(r.Context(), filmGenre)
	if err != nil {
		h.writeError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *FilmGenreHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidFilmGenreIDs):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}
