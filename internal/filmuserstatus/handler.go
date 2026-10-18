package FilmUserStatus

import (
	db "WebGrupal/db/sqlc"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Handler
type FilmUserStatusHandler struct {
	service *FilmUserStatusService
}

// Constructor
func NewFilmUserStatusHandler(service *FilmUserStatusService) *FilmUserStatusHandler {
	return &FilmUserStatusHandler{service: service}
}

func (h *FilmUserStatusHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/filmuserstatuses", h.handleCollection)
	mux.HandleFunc("/filmuserstatuses/", h.handleFilmUserStatus)
}

// writeJSON para escritura de JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error al codificar respuesta: %v", err)
	}
}

// Casos de uso: GET /filmuserstatuses/{filmID}/{userID}, PUT /filmuserstatuses/{filmID}/{userID}, DELETE /filmuserstatuses/{filmID}/{userID}
func (h *FilmUserStatusHandler) handleFilmUserStatus(w http.ResponseWriter, r *http.Request) {
	var filmID, userID int32
	if r.Method != http.MethodPost {
		var err error
		filmID, userID, err = extractIDs(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	switch r.Method {
	case http.MethodPost:
		h.createFilmUserStatus(w, r)
	case http.MethodGet:
		h.getFilmUserStatus(w, r, filmID, userID)
	case http.MethodPut:
		h.updateFilmUserStatus(w, r, filmID, userID)
	case http.MethodDelete:
		h.deleteFilmUserStatus(w, r, filmID, userID)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// POST /filmuserstatuses  y  GET /filmuserstatuses
func (h *FilmUserStatusHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createMultipleFilmUserStatuses(w, r)
	case http.MethodGet:
		h.getFilmUserStatuses(w, r)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// extractIDs obtiene filmID y userID de una ruta del tipo /filmuserstatuses/{filmID}/{userID}
func extractIDs(path string) (int32, int32, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Obtener las tres partes de la ruta: "filmuserstatuses", "filmID" y "userID"
	if len(parts) != 3 || parts[0] != "filmuserstatuses" {
		return 0, 0, errors.New("ruta inválida, se esperaba /filmuserstatuses/{filmID}/{userID}")
	}

	filmID, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil || filmID <= 0 {
		return 0, 0, errors.New("filmID inválido")
	}

	userID, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil || userID <= 0 {
		return 0, 0, errors.New("userID inválido")
	}

	return int32(filmID), int32(userID), nil
}

// GET /filmuserstatuses/{filmID}/{userID}: Debe obtener los IDs de la URL, buscar el registro y devolverlo. Si no existe, debe devolver un 404.
func (h *FilmUserStatusHandler) getFilmUserStatus(w http.ResponseWriter, r *http.Request, filmID int32, userID int32) {
	fus, err := h.service.GetFilmUserStatus(r.Context(), filmID, userID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, fus)
}

// PUT /filmuserstatuses/{filmID}/{userID}: Debe recibir datos JSON, validarlos, y actualizar el registro correspondiente.
func (h *FilmUserStatusHandler) updateFilmUserStatus(w http.ResponseWriter, r *http.Request, filmID int32, userID int32) {
	var receivedFUS db.Filmuserstatus
	if err := json.NewDecoder(r.Body).Decode(&receivedFUS); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if filmID != receivedFUS.FilmFilmid || userID != receivedFUS.UserUserid {
		http.Error(w, "Los IDs de la URL no coinciden con los IDs en el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	updatedFUS, err := h.service.Update(r.Context(), receivedFUS)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updatedFUS)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *FilmUserStatusHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrFilmUserStatusNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidFilmUserStatus),
		errors.Is(err, ErrInvalidFilmUserStatusIDs):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}

// GET /filmuserstatuses: Debe usar el método GetFilmUserStatuses de sqlc para obtener todos los registros y devolverlos como un array JSON.
func (h *FilmUserStatusHandler) getFilmUserStatuses(w http.ResponseWriter, r *http.Request) {
	statuses, err := h.service.GetFilmUserStatuses(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, statuses)
}

// POST /filmuserstatuses: Debe recibir datos en formato JSON, validarlos y usar el método CreateMultipleFilmUserStatuses.
func (h *FilmUserStatusHandler) createMultipleFilmUserStatuses(w http.ResponseWriter, r *http.Request) {
	var statuses []db.Filmuserstatus
	if err := json.NewDecoder(r.Body).Decode(&statuses); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateMultipleFilmUserStatuses(r.Context(), statuses)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created) // 201 (statusCreated)
}

// DELETE /filmuserstatuses/{filmID}/{userID}: Debe eliminar el registro y devolver un estado 204.
func (h *FilmUserStatusHandler) deleteFilmUserStatus(w http.ResponseWriter, r *http.Request, filmID int32, userID int32) {
	err := h.service.Delete(r.Context(), filmID, userID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FilmUserStatusHandler) createFilmUserStatus(w http.ResponseWriter, r *http.Request) {
	var receivedFUS db.Filmuserstatus
	if err := json.NewDecoder(r.Body).Decode(&receivedFUS); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	createdFUS, err := h.service.Create(r.Context(), receivedFUS)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createdFUS) // 201 (statusCreated)
}

