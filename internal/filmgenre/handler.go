package FilmGenre

import (
	dto "WebGrupal/internal"
	views "WebGrupal/views"
	"errors"
	"github.com/a-h/templ"
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

// GET /filmgenres: pide los registros al service y renderiza solo la lista (views.FilmGenreList).
func (h *FilmGenreHandler) ListFilmGenres(w http.ResponseWriter, r *http.Request) {
	filmGenres, err := h.service.GetFilmGenres(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}

	err = views.FilmGenreList(filmGenres).Render(r.Context(), w)
	if err != nil {
		log.Printf("Error renderizando la lista de FilmGenre: %v", err)
	}
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

	// Si es una petición HTMX, renderiza solo la lista actualizada
	if r.Header.Get("HX-Request") == "true" {
		h.ListFilmGenres(w, r)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DELETE /filmgenres/{genreid}/{filmid}: borra la relación película-género y responde 200 vacío,
// así HTMX reemplaza el elemento de la lista por nada.
func (h *FilmGenreHandler) Delete(w http.ResponseWriter, r *http.Request) {
	genreID, err := strconv.ParseInt(r.PathValue("genreid"), 10, 32)
	if err != nil || genreID <= 0 {
		http.Error(w, "ID de género inválido", http.StatusBadRequest)
		return
	}

	filmID, err := strconv.ParseInt(r.PathValue("filmid"), 10, 32)
	if err != nil || filmID <= 0 {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), int32(genreID), int32(filmID))
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GET /filmgenres/{genreid}/{filmid}: devuelve la fila en modo lectura (la usa el botón Cancelar).
func (h *FilmGenreHandler) Item(w http.ResponseWriter, r *http.Request) {
	genreID, okA := pathID(r, "genreid")
	filmID, okB := pathID(r, "filmid")
	if !okA || !okB {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	filmGenre, err := h.service.Get(r.Context(), genreID, filmID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.FilmGenreItem(filmGenre))
}

// GET /filmgenres/{genreid}/{filmid}/edit: devuelve la fila convertida en formulario de edición.
func (h *FilmGenreHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	genreID, okA := pathID(r, "genreid")
	filmID, okB := pathID(r, "filmid")
	if !okA || !okB {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	filmGenre, err := h.service.Get(r.Context(), genreID, filmID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.FilmGenreEditItem(filmGenre))
}

// PUT /filmgenres/{genreid}/{filmid}: la URL trae el par actual y el formulario el nuevo; devuelve la fila con el par nuevo.
func (h *FilmGenreHandler) Update(w http.ResponseWriter, r *http.Request) {
	genreID, okA := pathID(r, "genreid")
	filmID, okB := pathID(r, "filmid")
	if !okA || !okB {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	// Si el par actual no existe, el UPDATE no cambiaría nada sin dar error: se responde 404 antes.
	if _, err := h.service.Get(r.Context(), genreID, filmID); err != nil {
		h.writeError(w, err)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	newA, err := strconv.ParseInt(r.FormValue("genre_genreid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de género inválido", http.StatusBadRequest)
		return
	}

	newB, err := strconv.ParseInt(r.FormValue("film_filmid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	filmGenre := dto.Filmgenre{GenreGenreid: int32(newA), FilmFilmid: int32(newB)}
	if err := h.service.Update(r.Context(), genreID, filmID, filmGenre); err != nil {
		h.writeError(w, err)
		return
	}

	updated, err := h.service.Get(r.Context(), filmGenre.GenreGenreid, filmGenre.FilmFilmid)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.FilmGenreItem(updated))
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *FilmGenreHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidFilmGenreIDs):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrFilmGenreNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}

// pathID lee un ID positivo del segmento {name} de la ruta.
func pathID(r *http.Request, name string) (int32, bool) {
	id, err := strconv.ParseInt(r.PathValue(name), 10, 32)
	return int32(id), err == nil && id > 0
}

// render escribe el componente; si falla solo se puede loguear, porque el código HTTP ya se envió.
func render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	if err := c.Render(r.Context(), w); err != nil {
		log.Printf("Error renderizando: %v", err)
	}
}
