package Status

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
type StatusHandler struct {
	service *StatusService
}

// Constructor
func NewStatusHandler(service *StatusService) *StatusHandler {
	return &StatusHandler{service: service}
}

func (h *StatusHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/statuses", h.handleCollection)
	mux.HandleFunc("/statuses/", h.handleStatus)
}

// writeJSON para escritura de JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error al codificar respuesta: %v", err)
	}
}

// Casos de uso: GET /statuses/{id}, PUT /statuses/{id}, DELETE /statuses/{id}
func (h *StatusHandler) handleStatus(w http.ResponseWriter, r *http.Request) {
	var id int32
	if r.Method != http.MethodPost {
		var err error
		id, err = extractID(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	switch r.Method {
	case http.MethodPost:
		h.createStatus(w, r)
	case http.MethodGet:
		h.getStatus(w, r, id)
	case http.MethodPut:
		h.updateStatus(w, r, id)
	case http.MethodDelete:
		h.deleteStatus(w, r, id)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// POST /statuses  y  GET /statuses
func (h *StatusHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createMultipleStatuses(w, r)
	case http.MethodGet:
		h.getStatuses(w, r)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

func extractID(path string) (int32, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Obtener las dos partes de la ruta: "statuses" y "id"
	if len(parts) != 2 || parts[0] != "statuses" {
		return 0, errors.New("ruta inválida")
	}

	// Convertir "id" a int32 (segundos dos campos, base 10 y 32 bits)
	id, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil || id <= 0 {
		return 0, errors.New("id inválido")
	}
	return int32(id), nil
}

// GET /statuses/{id}: Obtiene el ID de la URL, busca el registro y lo devuelve. Si no existe, devuelve un 404.
func (h *StatusHandler) getStatus(w http.ResponseWriter, r *http.Request, id int32) {
	status, err := h.service.GetStatus(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, status)
}

// PUT /statuses/{id}: Recibe datos JSON, los valida, y actualiza el registro correspondiente.
func (h *StatusHandler) updateStatus(w http.ResponseWriter, r *http.Request, id int32) {
	var receivedStatus db.Status
	if err := json.NewDecoder(r.Body).Decode(&receivedStatus); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if id != receivedStatus.Statusid {
		http.Error(w, "El ID del status en la URL no coincide con el ID en el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	updatedStatus, err := h.service.Update(r.Context(), receivedStatus)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updatedStatus)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *StatusHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrStatusNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidStatusID),
		errors.Is(err, ErrInvalidStatusNameEmpty),
		errors.Is(err, ErrInvalidStatus):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}

// GET /statuses: Usa GetStatuses para obtener todos los registros y devolverlos como un array JSON.
func (h *StatusHandler) getStatuses(w http.ResponseWriter, r *http.Request) {
	statuses, err := h.service.GetStatuses(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, statuses)
}

// POST /statuses: Recibe datos en formato JSON (array), los valida y crea múltiples statuses.
func (h *StatusHandler) createMultipleStatuses(w http.ResponseWriter, r *http.Request) {
	var statuses []db.Status
	if err := json.NewDecoder(r.Body).Decode(&statuses); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateMultipleStatuses(r.Context(), statuses)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created) // 201 (statusCreated)
}

// DELETE /statuses/{id}: Elimina el registro y devuelve un estado 204.
func (h *StatusHandler) deleteStatus(w http.ResponseWriter, r *http.Request, id int32) {
	err := h.service.Delete(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *StatusHandler) createStatus(w http.ResponseWriter, r *http.Request) {
	var receivedStatus db.Status
	if err := json.NewDecoder(r.Body).Decode(&receivedStatus); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	createdStatus, err := h.service.Create(r.Context(), receivedStatus)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createdStatus) // 201 (statusCreated)
}

