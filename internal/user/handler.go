package User

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
type UserHandler struct {
	service *UserService
}

// Constructor
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.handleCollection)
	mux.HandleFunc("/users/", h.handleUser)
}

// writeJSON para escritura de JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error al codificar respuesta: %v", err)
	}
}

// Casos de uso: GET /users/{id}, PUT /users/{id}, DELETE /users/{id}
func (h *UserHandler) handleUser(w http.ResponseWriter, r *http.Request) {
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
		h.createUser(w, r)
	case http.MethodGet:
		h.getUser(w, r, id)
	case http.MethodPut:
		h.updateUser(w, r, id)
	case http.MethodDelete:
		h.deleteUser(w, r, id)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// POST /users  y  GET /users
func (h *UserHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createMultipleUsers(w, r)
	case http.MethodGet:
		h.getUsers(w, r)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

func extractID(path string) (int32, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Obtener las dos partes de la ruta: "users" y "id"
	if len(parts) != 2 || parts[0] != "users" {
		return 0, errors.New("ruta inválida")
	}

	// Convertir "id" a int32 ( base 10 y 32 bits)
	id, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil || id <= 0 {
		return 0, errors.New("id inválido")
	}
	return int32(id), nil
}

// GET /users/{id}: Obtiene el ID de la URL, busca el registro y lo devuelve. Si no existe, devuelve un 404.
func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request, id int32) {
	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// PUT /users/{id}: Recibe datos JSON, los valida, y actualiza el registro correspondiente.
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request, id int32) {
	var receivedUser dto.User
	if err := json.NewDecoder(r.Body).Decode(&receivedUser); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if id != receivedUser.Userid {
		http.Error(w, "El ID del user en la URL no coincide con el ID en el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	updatedUser, err := h.service.Update(r.Context(), receivedUser)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updatedUser)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *UserHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUserNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidUserID),
		errors.Is(err, ErrInvalidUserNameEmpty),
		errors.Is(err, ErrInvalidUserPasswordEmpty),
		errors.Is(err, ErrInvalidUserGenre),
		errors.Is(err, ErrInvalidUser):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}

// GET /users: Usa GetUsers para obtener todos los registros y devolverlos como un array JSON.
func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// POST /users: Recibe datos en formato JSON (array), los valida y crea múltiples users.
func (h *UserHandler) createMultipleUsers(w http.ResponseWriter, r *http.Request) {
	var users []dto.User
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateMultipleUsers(r.Context(), users)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created) // 201 (statusCreated)
}

// DELETE /users/{id}: Elimina el registro y devuelve un estado 204.
func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request, id int32) {
	err := h.service.Delete(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var receivedUser dto.User
	if err := json.NewDecoder(r.Body).Decode(&receivedUser); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.Create(r.Context(), receivedUser)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createdUser) // 201 (statusCreated)
}
