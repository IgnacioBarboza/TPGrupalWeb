package CastMember

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
type CastMemberHandler struct {
	service *CastMemberService
}

// Constructor
func NewCastMemberHandler(service *CastMemberService) *CastMemberHandler {
	return &CastMemberHandler{service: service}
}

// GET /castmembers: pide los registros al service y renderiza solo la lista (views.CastMemberList).
func (h *CastMemberHandler) ListCastMembers(w http.ResponseWriter, r *http.Request) {
	castMembers, err := h.service.GetCastMembers(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}

	err = views.CastMemberList(castMembers).Render(r.Context(), w)
	if err != nil {
		log.Printf("Error renderizando la lista de miembros del elenco: %v", err)
	}
}

// POST /castmembers: recibe el formulario, crea el registro y redirige a "/" (PRG).
func (h *CastMemberHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	castMember := dto.Castmember{
		Castmembername:   r.FormValue("castmembername"),
		Castmemberrole:   r.FormValue("castmemberrole"),
		Castmembergender: r.FormValue("castmembergender"),
	}

	_, err = h.service.Create(r.Context(), castMember)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Si es una petición HTMX, renderiza solo la lista actualizada
	if r.Header.Get("HX-Request") == "true" {
		h.ListCastMembers(w, r)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DELETE /castmembers/{id}: borra el miembro del elenco y responde 200 vacío,
// así HTMX reemplaza el elemento de la lista por nada.
func (h *CastMemberHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("HX-Trigger", `{"castMemberDeleted": {"target": "body"}}`)
	w.WriteHeader(http.StatusOK)
}

// GET /castmembers/{id}: devuelve la fila en modo lectura (la usa el botón Cancelar).
func (h *CastMemberHandler) Item(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	castMember, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.CastMemberItem(castMember))
}

// GET /castmembers/{id}/edit: devuelve la fila convertida en formulario de edición.
func (h *CastMemberHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	castMember, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.CastMemberEditItem(castMember))
}

// PUT /castmembers/{id}: recibe el formulario de edición, guarda los cambios y devuelve la fila actualizada.
func (h *CastMemberHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	castMember := dto.Castmember{
		Castmemberid:     id,
		Castmembername:   r.FormValue("castmembername"),
		Castmemberrole:   r.FormValue("castmemberrole"),
		Castmembergender: r.FormValue("castmembergender"),
	}

	if err := h.service.Update(r.Context(), castMember); err != nil {
		h.writeError(w, err)
		return
	}

	// Se relee de la base: un UPDATE sobre un ID inexistente no da error, el Get sí (404).
	updated, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.CastMemberItem(updated))
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *CastMemberHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidCastMemberNameEmpty),
		errors.Is(err, ErrInvalidCastMemberRole),
		errors.Is(err, ErrInvalidCastMemberGender):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrCastMemberNotFound):
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
