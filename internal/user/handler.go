package User

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
type UserHandler struct {
	service *UserService
}

// Constructor
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GET /users: pide los registros al service y renderiza solo la lista (views.UserList).
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}

	err = views.UserList(users).Render(r.Context(), w)
	if err != nil {
		log.Printf("Error renderizando la lista de usuarios: %v", err)
	}
}

// POST /users: recibe el formulario, crea el registro y redirige a "/" (PRG).
func (h *UserHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	birthday, err := time.Parse("2006-01-02", r.FormValue("userbirthday"))
	if err != nil {
		http.Error(w, "Fecha de nacimiento inválida", http.StatusBadRequest)
		return
	}

	user := dto.User{
		Username:           r.FormValue("username"),
		Userpasswordhashed: r.FormValue("userpasswordhashed"),
		Userbirthday:       birthday,
		Usergenre:          r.FormValue("usergenre"),
	}

	_, err = h.service.Create(r.Context(), user)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Si es una petición HTMX, renderiza solo la lista actualizada
	if r.Header.Get("HX-Request") == "true" {
		h.ListUsers(w, r)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DELETE /users/{id}: borra el usuario y responde 200 vacío,
// así HTMX reemplaza el elemento de la lista por nada.
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("HX-Trigger", `{"userDeleted": {"target": "body"}}`)
	w.WriteHeader(http.StatusOK)
}

// GET /users/{id}: devuelve la fila en modo lectura (la usa el botón Cancelar).
func (h *UserHandler) Item(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	user, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.UserItem(user))
}

// GET /users/{id}/edit: devuelve la fila convertida en formulario de edición.
func (h *UserHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	user, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.UserEditItem(user))
}

// PUT /users/{id}: recibe el formulario de edición, guarda los cambios y devuelve la fila actualizada.
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	birthday, err := time.Parse("2006-01-02", r.FormValue("userbirthday"))
	if err != nil {
		http.Error(w, "Fecha de nacimiento inválida", http.StatusBadRequest)
		return
	}

	// El form de edición no muestra la contraseña actual: si viene vacía, se conserva la guardada.
	password := r.FormValue("userpasswordhashed")
	if password == "" {
		current, err := h.service.Get(r.Context(), id)
		if err != nil {
			h.writeError(w, err)
			return
		}
		password = current.Userpasswordhashed
	}

	user := dto.User{
		Userid:             id,
		Username:           r.FormValue("username"),
		Userpasswordhashed: password,
		Userbirthday:       birthday,
		Usergenre:          r.FormValue("usergenre"),
	}

	if err := h.service.Update(r.Context(), user); err != nil {
		h.writeError(w, err)
		return
	}

	// Se relee de la base: un UPDATE sobre un ID inexistente no da error, el Get sí (404).
	updated, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.UserItem(updated))
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *UserHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidUserNameEmpty),
		errors.Is(err, ErrInvalidUserPasswordEmpty),
		errors.Is(err, ErrInvalidUserGenre):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrUserNotFound):
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
