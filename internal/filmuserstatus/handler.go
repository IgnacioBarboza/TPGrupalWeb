package FilmUserStatus

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
type FilmUserStatusHandler struct {
	service *FilmUserStatusService
}

// Constructor
func NewFilmUserStatusHandler(service *FilmUserStatusService) *FilmUserStatusHandler {
	return &FilmUserStatusHandler{service: service}
}

// GET /filmuserstatuses: pide los registros al service y renderiza solo la lista (views.FilmUserStatusList).
func (h *FilmUserStatusHandler) ListFilmUserStatuses(w http.ResponseWriter, r *http.Request) {
	filmUserStatuses, err := h.service.GetFilmUserStatuses(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}

	err = views.FilmUserStatusList(filmUserStatuses).Render(r.Context(), w)
	if err != nil {
		log.Printf("Error renderizando la lista de FilmUserStatus: %v", err)
	}
}

// POST /filmuserstatuses: recibe el formulario, crea el registro y redirige a "/" (PRG).
func (h *FilmUserStatusHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	filmID, err := strconv.ParseInt(r.FormValue("film_filmid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.FormValue("user_userid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	statusID, err := strconv.ParseInt(r.FormValue("status_statusid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de estado inválido", http.StatusBadRequest)
		return
	}

	filmUserStatus := dto.Filmuserstatus{
		FilmFilmid:     int32(filmID),
		UserUserid:     int32(userID),
		StatusStatusid: int32(statusID),
	}

	_, err = h.service.Create(r.Context(), filmUserStatus)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Si es una petición HTMX, renderiza solo la lista actualizada
	if r.Header.Get("HX-Request") == "true" {
		h.ListFilmUserStatuses(w, r)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DELETE /filmuserstatuses/{filmid}/{userid}: borra el estado de la película para el usuario y responde 200 vacío,
// así HTMX reemplaza el elemento de la lista por nada.
func (h *FilmUserStatusHandler) Delete(w http.ResponseWriter, r *http.Request) {
	filmID, err := strconv.ParseInt(r.PathValue("filmid"), 10, 32)
	if err != nil || filmID <= 0 {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("userid"), 10, 32)
	if err != nil || userID <= 0 {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), int32(filmID), int32(userID))
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GET /filmuserstatuses/{filmid}/{userid}: devuelve la fila en modo lectura (la usa el botón Cancelar).
func (h *FilmUserStatusHandler) Item(w http.ResponseWriter, r *http.Request) {
	filmID, okFilm := pathID(r, "filmid")
	userID, okUser := pathID(r, "userid")
	if !okFilm || !okUser {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	fus, err := h.service.Get(r.Context(), filmID, userID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.FilmUserStatusItem(fus))
}

// GET /filmuserstatuses/{filmid}/{userid}/edit: devuelve la fila convertida en formulario de edición.
func (h *FilmUserStatusHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	filmID, okFilm := pathID(r, "filmid")
	userID, okUser := pathID(r, "userid")
	if !okFilm || !okUser {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	fus, err := h.service.Get(r.Context(), filmID, userID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.FilmUserStatusEditItem(fus))
}

// PUT /filmuserstatuses/{filmid}/{userid}: recibe el formulario de edición, guarda los cambios y devuelve la fila actualizada.
func (h *FilmUserStatusHandler) Update(w http.ResponseWriter, r *http.Request) {
	filmID, okFilm := pathID(r, "filmid")
	userID, okUser := pathID(r, "userid")
	if !okFilm || !okUser {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	statusID, err := strconv.ParseInt(r.FormValue("status_statusid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de estado inválido", http.StatusBadRequest)
		return
	}

	fus := dto.Filmuserstatus{
		FilmFilmid:     filmID,
		UserUserid:     userID,
		StatusStatusid: int32(statusID),
	}

	if err := h.service.Update(r.Context(), fus); err != nil {
		h.writeError(w, err)
		return
	}

	// Se relee de la base: un UPDATE sobre un ID inexistente no da error, el Get sí (404).
	updated, err := h.service.Get(r.Context(), filmID, userID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.FilmUserStatusItem(updated))
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *FilmUserStatusHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidFilmUserStatusIDs):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrFilmUserStatusNotFound):
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
