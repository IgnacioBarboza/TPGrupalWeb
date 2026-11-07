package User

import (
	dto "WebGrupal/internal"
	"errors"
	"log"
	"net/http"
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *UserHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidUserNameEmpty),
		errors.Is(err, ErrInvalidUserPasswordEmpty),
		errors.Is(err, ErrInvalidUserGenre):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}
