package Review

import (
	dto "WebGrupal/internal"
	views "WebGrupal/views"
	"errors"
	"github.com/a-h/templ"
	"log"
	"net/http"
	"strconv"
)

// Handler
type ReviewHandler struct {
	service *ReviewService
}

// Constructor
func NewReviewHandler(service *ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

// GET /reviews: pide los registros al service y renderiza solo la lista (views.ReviewList).
func (h *ReviewHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.service.GetReviews(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}

	err = views.ReviewList(reviews).Render(r.Context(), w)
	if err != nil {
		log.Printf("Error renderizando la lista de reseñas: %v", err)
	}
}

// POST /reviews: recibe el formulario, crea el registro y redirige a "/" (PRG).
func (h *ReviewHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	filmID, err := strconv.ParseInt(r.FormValue("reviewfilmid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.FormValue("reviewuserid"), 10, 32)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	stars, err := strconv.ParseInt(r.FormValue("reviewstars"), 10, 16)
	if err != nil {
		http.Error(w, "Estrellas inválidas", http.StatusBadRequest)
		return
	}

	review := dto.Review{
		Reviewfilmid:      int32(filmID),
		Reviewuserid:      int32(userID),
		Reviewstars:       int16(stars),
		Reviewdescription: r.FormValue("reviewdescription"),
	}

	_, err = h.service.Create(r.Context(), review)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Si es una petición HTMX, renderiza solo la lista actualizada
	if r.Header.Get("HX-Request") == "true" {
		h.ListReviews(w, r)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DELETE /reviews/{filmid}/{userid}: borra la reseña y responde 200 vacío,
// así HTMX reemplaza el elemento de la lista por nada.
func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	filmID, err := strconv.ParseInt(r.PathValue("filmid"), 10, 32)
	if err != nil || filmID <= 0 {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("userid"), 10, 32)
	if err != nil || userID <= 0 {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), int32(filmID), int32(userID))
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GET /reviews/{filmid}/{userid}: devuelve la fila en modo lectura (la usa el botón Cancelar).
func (h *ReviewHandler) Item(w http.ResponseWriter, r *http.Request) {
	filmID, okFilm := pathID(r, "filmid")
	userID, okUser := pathID(r, "userid")
	if !okFilm || !okUser {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	review, err := h.service.Get(r.Context(), filmID, userID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.ReviewItem(review))
}

// GET /reviews/{filmid}/{userid}/edit: devuelve la fila convertida en formulario de edición.
func (h *ReviewHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	filmID, okFilm := pathID(r, "filmid")
	userID, okUser := pathID(r, "userid")
	if !okFilm || !okUser {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	review, err := h.service.Get(r.Context(), filmID, userID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.ReviewEditItem(review))
}

// PUT /reviews/{filmid}/{userid}: recibe el formulario de edición, guarda los cambios y devuelve la fila actualizada.
func (h *ReviewHandler) Update(w http.ResponseWriter, r *http.Request) {
	filmID, okFilm := pathID(r, "filmid")
	userID, okUser := pathID(r, "userid")
	if !okFilm || !okUser {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	stars, err := strconv.ParseInt(r.FormValue("reviewstars"), 10, 16)
	if err != nil {
		http.Error(w, "Estrellas inválidas", http.StatusBadRequest)
		return
	}

	review := dto.Review{
		Reviewfilmid:      filmID,
		Reviewuserid:      userID,
		Reviewstars:       int16(stars),
		Reviewdescription: r.FormValue("reviewdescription"),
	}

	if err := h.service.Update(r.Context(), review); err != nil {
		h.writeError(w, err)
		return
	}

	// Se relee de la base: un UPDATE sobre un ID inexistente no da error, el Get sí (404).
	updated, err := h.service.Get(r.Context(), filmID, userID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	render(w, r, views.ReviewItem(updated))
}

// writeError traduce errores del negocio a codigos HTTP.
func (h *ReviewHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidReviewStars),
		errors.Is(err, ErrInvalidReviewDescriptionEmpty):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrReviewNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		log.Printf("error interno: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
	}
}

// pathID lee un ID positivo del segmento {name} de la ruta.
func pathID(r *http.Request, name string) (int32, bool) {
	id, err := strconv.ParseInt(r.PathValue(name), 10, 32)
	return int32(id), err == nil && id > 0
}

// render escribe el componente; si falla solo se puede loguear, porque el código HTTP ya se envió.
func render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	if err := c.Render(r.Context(), w); err != nil {
		log.Printf("Error renderizando: %v", err)
	}
}
