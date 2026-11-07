package CastMember

import (
	dto "WebGrupal/internal"
	"errors"
	"log"
	"net/http"
)

// Handler
type CastMemberHandler struct {
	service *CastMemberService
}

// Constructor
func NewCastMemberHandler(service *CastMemberService) *CastMemberHandler {
	return &CastMemberHandler{service: service}
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *CastMemberHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidCastMemberNameEmpty),
		errors.Is(err, ErrInvalidCastMemberRole),
		errors.Is(err, ErrInvalidCastMemberGender):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}
