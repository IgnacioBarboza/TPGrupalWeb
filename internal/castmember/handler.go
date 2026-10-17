package CastMember

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
type CastMemberHandler struct {
	service *CastMemberService
}

// Constructor
func NewCastMemberHandler(service *CastMemberService) *CastMemberHandler {
	return &CastMemberHandler{service: service}
}

func (h *CastMemberHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/castmembers", h.handleCollection)
	mux.HandleFunc("/castmembers/", h.handleCastMember)
}

// writeJSON para escritura de JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error al codificar respuesta: %v", err)
	}
}

// Casos de uso: GET /castmembers/{id}, PUT /castmembers/{id}, DELETE /castmembers/{id}
func (h *CastMemberHandler) handleCastMember(w http.ResponseWriter, r *http.Request) {
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
		h.createCastMember(w, r)
	case http.MethodGet:
		h.getCastMember(w, r, id)
	case http.MethodPut:
		h.updateCastMember(w, r, id)
	case http.MethodDelete:
		h.deleteCastMember(w, r, id)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// POST /castmembers  y  GET /castmembers
func (h *CastMemberHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createMultipleCastMembers(w, r)
	case http.MethodGet:
		h.getCastMembers(w, r)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

func extractID(path string) (int32, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Obtener las dos partes de la ruta: "castmembers" y "id"
	if len(parts) != 2 || parts[0] != "castmembers" {
		return 0, errors.New("ruta inválida")
	}

	// Convertir "id" a int32 (segundos dos campos, base 10 y 32 bits)
	id, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil || id <= 0 {
		return 0, errors.New("id inválido")
	}
	return int32(id), nil
}

// GET /castmembers/{id}: Obtiene el ID de la URL, busca el registro y lo devuelve. Si no existe, devuelve un 404.
func (h *CastMemberHandler) getCastMember(w http.ResponseWriter, r *http.Request, id int32) {
	castMember, err := h.service.GetCastMember(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, castMember)
}

// PUT /castmembers/{id}: Recibe datos JSON, los valida, y actualiza el registro correspondiente.
func (h *CastMemberHandler) updateCastMember(w http.ResponseWriter, r *http.Request, id int32) {
	var receivedCastMember db.Castmember
	if err := json.NewDecoder(r.Body).Decode(&receivedCastMember); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if id != receivedCastMember.Castmemberid {
		http.Error(w, "El ID del castMember en la URL no coincide con el ID en el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	updatedCastMember, err := h.service.Update(r.Context(), receivedCastMember)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updatedCastMember)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *CastMemberHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrCastMemberNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidCastMemberID),
		errors.Is(err, ErrInvalidCastMemberNameEmpty),
		errors.Is(err, ErrInvalidCastMemberRole),
		errors.Is(err, ErrInvalidCastMemberGender),
		errors.Is(err, ErrInvalidCastMember):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}

// GET /castmembers: Usa GetCastMembers para obtener todos los registros y devolverlos como un array JSON.
func (h *CastMemberHandler) getCastMembers(w http.ResponseWriter, r *http.Request) {
	castMembers, err := h.service.GetCastMembers(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, castMembers)
}

// POST /castmembers: Recibe datos en formato JSON (array), los valida y crea múltiples castMembers.
func (h *CastMemberHandler) createMultipleCastMembers(w http.ResponseWriter, r *http.Request) {
	var castMembers []db.Castmember
	if err := json.NewDecoder(r.Body).Decode(&castMembers); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateMultipleCastMembers(r.Context(), castMembers)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created) // 201 (statusCreated)
}

// DELETE /castmembers/{id}: Elimina el registro y devuelve un estado 204.
func (h *CastMemberHandler) deleteCastMember(w http.ResponseWriter, r *http.Request, id int32) {
	err := h.service.Delete(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CastMemberHandler) createCastMember(w http.ResponseWriter, r *http.Request) {
	var receivedCastMember db.Castmember
	if err := json.NewDecoder(r.Body).Decode(&receivedCastMember); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	createdCastMember, err := h.service.Create(r.Context(), receivedCastMember)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createdCastMember) // 201 (statusCreated)
}

