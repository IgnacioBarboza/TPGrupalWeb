package Genre

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
type GenreHandler struct {
	service *GenreService
}

// Constructor
func NewGenreHandler(service *GenreService) *GenreHandler {
	return &GenreHandler{service: service}
}

func (h *GenreHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/genres", h.handleCollection)
	mux.HandleFunc("/genres/", h.handleGenre)
}

// writeJSON para escritura de JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error al codificar respuesta: %v", err)
	}
}

// Casos de uso: GET /genres/{id}, PUT /genres/{id}, DELETE /genres/{id}
func (h *GenreHandler) handleGenre(w http.ResponseWriter, r *http.Request) {
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
		h.createGenre(w, r)
	case http.MethodGet:
		h.getGenre(w, r, id)
	case http.MethodPut:
		h.updateGenre(w, r, id)
	case http.MethodDelete:
		h.deleteGenre(w, r, id)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// POST /genres  y  GET /genres
func (h *GenreHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createMultipleGenres(w, r)
	case http.MethodGet:
		h.getGenres(w, r)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

func extractID(path string) (int32, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Obtener las dos partes de la ruta: "genres" y "id"
	if len(parts) != 2 || parts[0] != "genres" {
		return 0, errors.New("ruta inválida")
	}

	// Convertir "id" a int32 (base 10 y 32 bits)
	id, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil || id <= 0 {
		return 0, errors.New("id inválido")
	}
	return int32(id), nil
}

// GET /genres/{id}: Debe obtener el ID de la URL, buscar el registro con GetGenre y devolverlo. Si no existe, debe devolver un 404.
func (h *GenreHandler) getGenre(w http.ResponseWriter, r *http.Request, id int32) {
	genre, err := h.service.GetGenre(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, genre)
}

// PUT /genres/{id}: Debe recibir datos JSON, validarlos, y actualizar el registro correspondiente.
func (h *GenreHandler) updateGenre(w http.ResponseWriter, r *http.Request, id int32) {
	var receivedGenre db.Genre
	if err := json.NewDecoder(r.Body).Decode(&receivedGenre); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if id != receivedGenre.Genreid {
		http.Error(w, "El ID del género en la URL no coincide con el ID en el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	updatedGenre, err := h.service.Update(r.Context(), receivedGenre)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updatedGenre)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *GenreHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrGenreNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidGenre),
		errors.Is(err, ErrInvalidGenreNameEmpty):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}

// GET /genres: Debe usar el método GetGenres de sqlc para obtener todos los registros y devolverlos como un array JSON.
func (h *GenreHandler) getGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.service.GetGenres(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, genres)
}

// POST /genres: Debe recibir datos en formato JSON, validarlos y usar el método CreateMultipleGenres.
func (h *GenreHandler) createMultipleGenres(w http.ResponseWriter, r *http.Request) {
	var genres []db.Genre
	if err := json.NewDecoder(r.Body).Decode(&genres); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateMultipleGenres(r.Context(), genres)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created) // 201 (statusCreated)
}

// DELETE /genres/{id}: Debe eliminar el registro usando DeleteGenre y devolver un estado 204.
func (h *GenreHandler) deleteGenre(w http.ResponseWriter, r *http.Request, id int32) {
	err := h.service.Delete(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GenreHandler) createGenre(w http.ResponseWriter, r *http.Request) {
	var receivedGenre db.Genre
	if err := json.NewDecoder(r.Body).Decode(&receivedGenre); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	createdGenre, err := h.service.Create(r.Context(), receivedGenre)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createdGenre) // 201 (statusCreated)
}

