package main

import(
	"os" // leer env
	"fmt"
  "log"
	"database/sql"
  "errors"
  "encoding/json"
  "strconv"
  "net/http"
  db "WebGrupal/db/sqlc"
	_ "github.com/lib/pq" //Driver postgresql
)

func OpenDB() *sql.DB {

	// Seteo variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := "disable"

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbname,
		sslmode,
	)

	conn, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("Error al abrir la BD: %v", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	return conn
}
//////
type Server struct {
	queries *db.Queries
}

// GetFilmHandler responde GET /film?id=1
func (s *Server) GetFilmHandler(w http.ResponseWriter, r *http.Request) {
	// Si NO es tipo GET da error
	if r.Method != http.MethodGet {
		http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leemos y validamos el parametro id
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id invalido", http.StatusBadRequest)
		return
	}

	// Obtenemos la pelicula cuyo id sea el entregado
	film, err := s.queries.GetFilm(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "film no encontrado", http.StatusNotFound)
			return
		}
		log.Printf("error en GetFilm: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	// Cambiamos el header para notificar que sea de tipo JSON y usamos Encode para crear el json apartir del struct
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(film); err != nil {
		log.Printf("error al codificar respuesta: %v", err)
	}
}
//////

func main() {
  // Abrimos conexion
  conn := OpenDB();
  defer conn.Close();
  
  queries := db.New(conn);
  
  server := &Server{
      queries: queries,
  }

  // 2. Asociamos la ruta "/film" a nuestro handler
  http.HandleFunc("/film", server.GetFilmHandler)

  // 3. Levantamos el servidor en el puerto 8080
  port := ":8080"
  fmt.Printf("Servidor corriendo en el puerto %s\n", port)
  if err := http.ListenAndServe(port, nil); err != nil {
      log.Fatalf("Error al iniciar el servidor: %v", err)
  }
}