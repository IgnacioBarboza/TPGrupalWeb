package FilmCast

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
type FilmCastHandler struct {
	service *FilmCastService
}

// Constructor
func NewFilmCastHandler(service *FilmCastService) *FilmCastHandler {
	return &FilmCastHandler{service: service}
}

// GET /filmcasts: pide los registros al service y renderiza solo la lista (views.FilmCastList).
func (h *FilmCastHandler) ListFilmCasts(w http.ResponseWriter, r *http.Request) {
	filmCasts, err := h.service.GetFilmCasts(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}

	err = views.FilmCastList(filmCasts).Render(r.Context(), w)
	if err != nil {
		log.Printf("Error renderizando la lista de FilmCast: %v", err)
	}
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

	// Si es una petición HTMX, renderiza solo la lista actualizada
	if r.Header.Get("HX-Request") == "true" {
		h.ListFilmCasts(w, r)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DELETE /filmcasts/{filmid}/{castmemberid}: borra la relación película-elenco y responde 200 vacío,
// así HTMX reemplaza el elemento de la lista por nada.
func (h *FilmCastHandler) Delete(w http.ResponseWriter, r *http.Request) {
	filmID, err := strconv.ParseInt(r.PathValue("filmid"), 10, 32)
	if err != nil || filmID <= 0 {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	castMemberID, err := strconv.ParseInt(r.PathValue("castmemberid"), 10, 32)
	if err != nil || castMemberID <= 0 {
		http.Error(w, "ID de miembro del elenco inválido", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), int32(filmID), int32(castMemberID))
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GET /filmcasts/{filmid}/{castmemberid}: devuelve la fila en modo lectura (la usa el botón Cancelar).
func (h *FilmCastHandler) Item(w http.ResponseWriter, r *http.Request) {
	filmID, okA := pathID(r, "filmid")
	castMemberID, okB := pathID(r, "castmemberid")
	if !okA || !okB {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	filmCast, err := h.service.Get(r.Context(), filmID, castMemberID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.FilmCastItem(filmCast))
}

// GET /filmcasts/{filmid}/{castmemberid}/edit: devuelve la fila convertida en formulario de edición.
func (h *FilmCastHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	filmID, okA := pathID(r, "filmid")
	castMemberID, okB := pathID(r, "castmemberid")
	if !okA || !okB {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	filmCast, err := h.service.Get(r.Context(), filmID, castMemberID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.FilmCastEditItem(filmCast))
}

// PUT /filmcasts/{filmid}/{castmemberid}: la URL trae el par actual y el formulario el nuevo; devuelve la fila con el par nuevo.
func (h *FilmCastHandler) Update(w http.ResponseWriter, r *http.Request) {
	filmID, okA := pathID(r, "filmid")
	castMemberID, okB := pathID(r, "castmemberid")
	if !okA || !okB {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	// Si el par actual no existe, el UPDATE no cambiaría nada sin dar error: se responde 404 antes.
	if _, err := h.service.Get(r.Context(), filmID, castMemberID); err != nil {
		h.writeError(w, err)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	newA, err := strconv.ParseInt(r.FormValue("filmid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	newB, err := strconv.ParseInt(r.FormValue("castmemberid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de miembro del elenco inválido", http.StatusBadRequest)
		return
	}

	filmCast := dto.Filmcast{Filmid: int32(newA), Castmemberid: int32(newB)}
	if err := h.service.Update(r.Context(), filmID, castMemberID, filmCast); err != nil {
		h.writeError(w, err)
		return
	}

	updated, err := h.service.Get(r.Context(), filmCast.Filmid, filmCast.Castmemberid)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.FilmCastItem(updated))
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *FilmCastHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidFilmCastIDs):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrFilmCastNotFound):
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
