package Film

import (
	dto "WebGrupal/internal"
	views "WebGrupal/views"
	"errors"
	"github.com/a-h/templ"
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

// GET /films: pide los registros al service y renderiza solo la lista (views.FilmList).
func (h *FilmHandler) ListFilms(w http.ResponseWriter, r *http.Request) {
	films, err := h.service.GetFilms(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}

	err = views.FilmList(films).Render(r.Context(), w)
	if err != nil {
		log.Printf("Error renderizando la lista de películas: %v", err)
	}
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

	// Si es una petición HTMX, renderiza solo la lista de películas actualizada
	if r.Header.Get("HX-Request") == "true" {
		h.ListFilms(w, r)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DELETE /films/{id}: borra la película y responde 200 vacío,
// así HTMX reemplaza el elemento de la lista por nada.
func (h *FilmHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), int32(id))
	if err != nil {
		h.writeError(w, err)
		return
	}

	// ON DELETE CASCADE también borró filas de otras tablas: se avisa para que esas listas se vuelvan a pedir.
	// El evento va a body porque el botón que hizo el pedido ya no está en la página cuando llega la respuesta.
	w.Header().Set("HX-Trigger", `{"filmDeleted": {"target": "body"}}`)
	w.WriteHeader(http.StatusOK)
}

// GET /films/{id}: devuelve la fila en modo lectura (la usa el botón Cancelar).
func (h *FilmHandler) Item(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	film, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.FilmItem(film))
}

// GET /films/{id}/edit: devuelve la fila convertida en formulario de edición.
func (h *FilmHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	film, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.FilmEditItem(film))
}

// PUT /films/{id}: recibe el formulario de edición, guarda los cambios y devuelve la fila actualizada.
func (h *FilmHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
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
		Filmid:          id,
		Filmtitle:       r.FormValue("filmtitle"),
		Filmreleasedate: releaseDate,
		Filmduration:    int32(duration),
	}

	if err := h.service.Update(r.Context(), film); err != nil {
		h.writeError(w, err)
		return
	}

	// Se relee de la base: un UPDATE sobre un ID inexistente no da error, el Get sí (404).
	updated, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.FilmItem(updated))
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *FilmHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidFilmTitleEmpty),
		errors.Is(err, ErrInvalidFilmDuration):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrFilmNotFound):
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
