package Status

import (
	dto "WebGrupal/internal"
	"errors"
	"log"
	"net/http"
)

// Handler
type StatusHandler struct {
	service *StatusService
}

// Constructor
func NewStatusHandler(service *StatusService) *StatusHandler {
	return &StatusHandler{service: service}
}

// POST /statuses: recibe el formulario, crea el registro y redirige a "/" (PRG).
func (h *StatusHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	status := dto.Status{
		Statusname: r.FormValue("statusname"),
	}

	_, err = h.service.Create(r.Context(), status)
	if err != nil {
		h.writeError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *StatusHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidStatusNameEmpty):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}
