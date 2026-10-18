package FilmGenre

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
type FilmGenreHandler struct {
	service *FilmGenreService
}

// Constructor
func NewFilmGenreHandler(service *FilmGenreService) *FilmGenreHandler {
	return &FilmGenreHandler{service: service}
}

func (h *FilmGenreHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/filmgenres", h.handleCollection)
	mux.HandleFunc("/filmgenres/", h.handleFilmGenre)
}

// writeJSON para escritura de JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error al codificar respuesta: %v", err)
	}
}

// Casos de uso: GET /filmgenres/{genreID}/{filmID}, DELETE /filmgenres/{genreID}/{filmID}
// FilmGenre es una tabla de union sin campos propios, por lo que no tiene operacion PUT.
func (h *FilmGenreHandler) handleFilmGenre(w http.ResponseWriter, r *http.Request) {
	var genreID, filmID int32
	if r.Method != http.MethodPost {
		var err error
		genreID, filmID, err = extractIDs(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	switch r.Method {
	case http.MethodPost:
		h.createFilmGenre(w, r)
	case http.MethodGet:
		h.getFilmGenre(w, r, genreID, filmID)
	case http.MethodDelete:
		h.deleteFilmGenre(w, r, genreID, filmID)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// POST /filmgenres  y  GET /filmgenres
func (h *FilmGenreHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createMultipleFilmGenres(w, r)
	case http.MethodGet:
		h.getFilmGenres(w, r)
	default:
		http.Error(w, ErrMetodoNoPermitido.Error(), http.StatusMethodNotAllowed)
	}
}

// extractIDs obtiene genreID y filmID de una ruta del tipo /filmgenres/{genreID}/{filmID}
func extractIDs(path string) (int32, int32, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Obtener las tres partes de la ruta: "filmgenres", "genreID" y "filmID"
	if len(parts) != 3 || parts[0] != "filmgenres" {
		return 0, 0, errors.New("ruta inválida, se esperaba /filmgenres/{genreID}/{filmID}")
	}

	genreID, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil || genreID <= 0 {
		return 0, 0, errors.New("genreID inválido")
	}

	filmID, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil || filmID <= 0 {
		return 0, 0, errors.New("filmID inválido")
	}

	return int32(genreID), int32(filmID), nil
}

// GET /filmgenres/{genreID}/{filmID}: Debe obtener los IDs de la URL, buscar el registro y devolverlo. Si no existe, debe devolver un 404.
func (h *FilmGenreHandler) getFilmGenre(w http.ResponseWriter, r *http.Request, genreID int32, filmID int32) {
	filmGenre, err := h.service.GetFilmGenre(r.Context(), genreID, filmID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, filmGenre)
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *FilmGenreHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrFilmGenreNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidFilmGenre),
		errors.Is(err, ErrInvalidFilmGenreIDs):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}

// GET /filmgenres: Debe usar el método GetFilmGenres de sqlc para obtener todos los registros y devolverlos como un array JSON.
func (h *FilmGenreHandler) getFilmGenres(w http.ResponseWriter, r *http.Request) {
	filmGenres, err := h.service.GetFilmGenres(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, filmGenres)
}

// POST /filmgenres: Debe recibir datos en formato JSON, validarlos y usar el método CreateMultipleFilmGenres.
func (h *FilmGenreHandler) createMultipleFilmGenres(w http.ResponseWriter, r *http.Request) {
	var filmGenres []db.Filmgenre
	if err := json.NewDecoder(r.Body).Decode(&filmGenres); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateMultipleFilmGenres(r.Context(), filmGenres)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created) // 201 (statusCreated)
}

// DELETE /filmgenres/{genreID}/{filmID}: Debe eliminar el registro y devolver un estado 204.
func (h *FilmGenreHandler) deleteFilmGenre(w http.ResponseWriter, r *http.Request, genreID int32, filmID int32) {
	err := h.service.Delete(r.Context(), genreID, filmID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FilmGenreHandler) createFilmGenre(w http.ResponseWriter, r *http.Request) {
	var receivedFilmGenre db.Filmgenre
	if err := json.NewDecoder(r.Body).Decode(&receivedFilmGenre); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	createdFilmGenre, err := h.service.Create(r.Context(), receivedFilmGenre)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createdFilmGenre) // 201 (statusCreated)
}

