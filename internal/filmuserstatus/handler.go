package FilmUserStatus

import (
	dto "WebGrupal/internal"
	"errors"
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *FilmUserStatusHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidFilmUserStatusIDs):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}
