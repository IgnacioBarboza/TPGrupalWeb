package FilmCast

import (
	dto "WebGrupal/internal"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Handler
type FilmCastHandler struct {
	service *FilmCastService
}

// Constructor
func NewFilmCastHandler(service *FilmCastService) *FilmCastHandler {
	return &FilmCastHandler{service: service}
}

func (h *FilmCastHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/filmcasts", h.handleCollection)
	mux.HandleFunc("/filmcasts/", h.handleFilmCast)
}

// writeJSON para escritura de JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error al codificar respuesta: %v", err)
	}
}

// Casos de uso: GET /filmcasts/{filmID}/{castMemberID}, DELETE /filmcasts/{filmID}/{castMemberID}
// FilmCast es una tabla de union sin campos propios, por lo que no tiene operacion PUT.
func (h *FilmCastHandler) handleFilmCast(w http.ResponseWriter, r *http.Request) {
	var filmID, castMemberID int32
	if r.Method != http.MethodPost {
		var err error
		filmID, castMemberID, err = extractIDs(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	switch r.Method {
	case http.MethodPost:
		h.createFilmCast(w, r)
	case http.MethodGet:
		h.getFilmCast(w, r, filmID, castMemberID)
	case http.MethodDelete:
		h.deleteFilmCast(w, r, filmID, castMemberID)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// POST /filmcasts  y  GET /filmcasts
func (h *FilmCastHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createMultipleFilmCasts(w, r)
	case http.MethodGet:
		h.getFilmCasts(w, r)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// extractIDs obtiene filmID y castMemberID de una ruta del tipo /filmcasts/{filmID}/{castMemberID}
func extractIDs(path string) (int32, int32, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Obtener las tres partes de la ruta: "filmcasts", "filmID" y "castMemberID"
	if len(parts) != 3 || parts[0] != "filmcasts" {
		return 0, 0, errors.New("ruta inválida, se esperaba /filmcasts/{filmID}/{castMemberID}")
	}

	filmID, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil || filmID <= 0 {
		return 0, 0, errors.New("filmID inválido")
	}

	castMemberID, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil || castMemberID <= 0 {
		return 0, 0, errors.New("castMemberID inválido")
	}

	return int32(filmID), int32(castMemberID), nil
}

// GET /filmcasts/{filmID}/{castMemberID}: Debe obtener los IDs de la URL, buscar el registro y devolverlo. Si no existe, debe devolver un 404.
func (h *FilmCastHandler) getFilmCast(w http.ResponseWriter, r *http.Request, filmID int32, castMemberID int32) {
	filmCast, err := h.service.GetFilmCast(r.Context(), filmID, castMemberID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, filmCast)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *FilmCastHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrFilmCastNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidFilmCast),
		errors.Is(err, ErrInvalidFilmCastIDs):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}

// GET /filmcasts: Debe usar el método GetFilmCasts de sqlc para obtener todos los registros y devolverlos como un array JSON.
func (h *FilmCastHandler) getFilmCasts(w http.ResponseWriter, r *http.Request) {
	filmCasts, err := h.service.GetFilmCasts(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, filmCasts)
}

// POST /filmcasts: Debe recibir datos en formato JSON, validarlos y usar el método CreateMultipleFilmCasts.
func (h *FilmCastHandler) createMultipleFilmCasts(w http.ResponseWriter, r *http.Request) {
	var filmCasts []dto.Filmcast
	if err := json.NewDecoder(r.Body).Decode(&filmCasts); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateMultipleFilmCasts(r.Context(), filmCasts)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created) // 201 (statusCreated)
}

// DELETE /filmcasts/{filmID}/{castMemberID}: Debe eliminar el registro y devolver un estado 204.
func (h *FilmCastHandler) deleteFilmCast(w http.ResponseWriter, r *http.Request, filmID int32, castMemberID int32) {
	err := h.service.Delete(r.Context(), filmID, castMemberID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FilmCastHandler) createFilmCast(w http.ResponseWriter, r *http.Request) {
	var receivedFilmCast dto.Filmcast
	if err := json.NewDecoder(r.Body).Decode(&receivedFilmCast); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	createdFilmCast, err := h.service.Create(r.Context(), receivedFilmCast)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createdFilmCast) // 201 (statusCreated)
}
