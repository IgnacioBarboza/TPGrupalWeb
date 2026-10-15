package Film

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
type FilmHandler struct {
	service *FilmService
}

// Constructor
func NewFilmHandler(service *FilmService) *FilmHandler {
	return &FilmHandler{service: service}
}

func (h *FilmHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/films", h.handleCollection)
	mux.HandleFunc("/films/", h.handleFilm)
}

// writeJSON para escritura de JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error al codificar respuesta: %v", err)
	}
}

// Casos de uso: GET /films/{id}, PUT /films/{id}, DELETE /films/{id}
func (h *FilmHandler) handleFilm(w http.ResponseWriter, r *http.Request) {
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
		h.createFilm(w, r)
	case http.MethodGet:
		h.getFilm(w, r, id)
	case http.MethodPut:
		h.updateFilm(w, r, id)
	case http.MethodDelete:
		h.deleteFilm(w, r, id)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// POST /films  y  GET /films
func (h *FilmHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createMultipleFilms(w, r)
	case http.MethodGet:
		h.getFilms(w, r)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

func extractID(path string) (int32, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Obtener las dos partes de la ruta: "films" y "id"
	if len(parts) != 2 || parts[0] != "films" {
		return 0, errors.New("ruta inválida")
	}

	// Convertir "id" a int32 (segundos dos campos, base 10 y 32 bits)
	id, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil || id <= 0 {
		return 0, errors.New("id inválido")
	}
	return int32(id), nil
}

// GET /<entidades>/{id}: Debe obtener el ID de la URL, buscar el registro con Get... y devolverlo. Si no existe, debe devolver un 404.
func (h *FilmHandler) getFilm(w http.ResponseWriter, r *http.Request, id int32) {
	film, err := h.service.GetFilm(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, film)
}

// PUT /<entidades>/{id}: Debe recibir datos JSON, validarlos, y actualizar el registro correspondiente usando
func (h *FilmHandler) updateFilm(w http.ResponseWriter, r *http.Request, id int32) {
	var receivedfilm db.Film
	if err := json.NewDecoder(r.Body).Decode(&receivedfilm); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if id != receivedfilm.Filmid {
		http.Error(w, "El ID del film en la URL no coincide con el ID en el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	uptadedfilm, err := h.service.Update(r.Context(), receivedfilm)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, uptadedfilm)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *FilmHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrFilmNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidFilmID),
		errors.Is(err, ErrInvalidFilmTitleEmpty),
		errors.Is(err, ErrInvalidFilmDuration):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}

// GET /<entidades>: Debe usar el método List... de sqlc para obtener todos los registros y devolverlos como un array JSON.
func (h *FilmHandler) getFilms(w http.ResponseWriter, r *http.Request) {
	films, err := h.service.GetFilms(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, films)
}

// POST /<entidades>: Debe recibir datos en formato JSON, validarlos (ej. que el título no esté vacío), usar el método Create...
func (h *FilmHandler) createMultipleFilms(w http.ResponseWriter, r *http.Request) {
	var films []db.Film
	if err := json.NewDecoder(r.Body).Decode(&films); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateMultipleFilms(r.Context(), films)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created) // 201 (statusCreated)
}

// DELETE /<entidades>/{id}: Debe eliminar el registro usando Delete... y devolver un estado 204
func (h *FilmHandler) deleteFilm(w http.ResponseWriter, r *http.Request, id int32) {
	err := h.service.Delete(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FilmHandler) createFilm(w http.ResponseWriter, r *http.Request) {
	var receivedfilm db.Film
	if err := json.NewDecoder(r.Body).Decode(&receivedfilm); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	createdfilm, err := h.service.Create(r.Context(), receivedfilm)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createdfilm) // 201 (statusCreated)
}
