package Genre

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
type GenreHandler struct {
	service *GenreService
}

// Constructor
func NewGenreHandler(service *GenreService) *GenreHandler {
	return &GenreHandler{service: service}
}

// GET /genres: pide los registros al service y renderiza solo la lista (views.GenreList).
func (h *GenreHandler) ListGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.service.GetGenres(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}

	err = views.GenreList(genres).Render(r.Context(), w)
	if err != nil {
		log.Printf("Error renderizando la lista de géneros: %v", err)
	}
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

	// Si es una petición HTMX, renderiza solo la lista actualizada
	if r.Header.Get("HX-Request") == "true" {
		h.ListGenres(w, r)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DELETE /genres/{id}: borra el género y responde 200 vacío,
// así HTMX reemplaza el elemento de la lista por nada.
func (h *GenreHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("HX-Trigger", `{"genreDeleted": {"target": "body"}}`)
	w.WriteHeader(http.StatusOK)
}

// GET /genres/{id}: devuelve la fila en modo lectura (la usa el botón Cancelar).
func (h *GenreHandler) Item(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	genre, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.GenreItem(genre))
}

// GET /genres/{id}/edit: devuelve la fila convertida en formulario de edición.
func (h *GenreHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	genre, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.GenreEditItem(genre))
}

// PUT /genres/{id}: recibe el formulario de edición, guarda los cambios y devuelve la fila actualizada.
func (h *GenreHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	genre := dto.Genre{
		Genreid:   id,
		Genrename: r.FormValue("genrename"),
	}

	if err := h.service.Update(r.Context(), genre); err != nil {
		h.writeError(w, err)
		return
	}

	// Se relee de la base: un UPDATE sobre un ID inexistente no da error, el Get sí (404).
	updated, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.GenreItem(updated))
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *GenreHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidGenreNameEmpty):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrGenreNotFound):
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
