#!/usr/bin/env bash
# ==============================================================================
# Script de prueba completo para todos los endpoints de la API
# Prueba: Crear (POST), Listar (GET), Ver (GET por ID/clave), Actualizar (PUT) y Eliminar (DELETE)
# ==============================================================================

BASE_URL="http://localhost:8080"

# INICIO CAMBIOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
ERROR_LOG_FILE=$(mktemp)
trap 'rm -f "$ERROR_LOG_FILE"' EXIT

# Wrapper de curl para contar y registrar cualquier error (fallo de conexión o HTTP >= 400)
curl() {
  local tmp_headers
  tmp_headers=$(mktemp)
  local exit_code=0
  command curl -D "$tmp_headers" "$@" || exit_code=$?

  local http_code
  http_code=$(awk '/^HTTP\// {code=$2} END {print code}' "$tmp_headers" | tr -d '\r')
  rm -f "$tmp_headers"

  if [ "$exit_code" -ne 0 ] || [ -z "$http_code" ] || [ "$http_code" -ge 400 ] 2>/dev/null; then
    local method="GET"
    local url=""
    local prev=""
    for arg in "$@"; do
      if [ "$prev" = "-X" ]; then
        method="$arg"
      fi
      case "$arg" in
        http://*|https://*) url="$arg" ;;
      esac
      prev="$arg"
    done
    echo "HTTP ${http_code:-ERR_CONN} [$method $url]" >> "$ERROR_LOG_FILE"
  fi
  return 0
}

# FIN CAMBIOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO


# Helper para extraer IDs de JSON sin requerir dependencias externas complejas
extract_id() {
  local json="$1"
  local key="$2"
  if command -v jq >/dev/null 2>&1; then
    echo "$json" | jq -r ".$key // empty"
  else
    echo "$json" | grep -o "\"$key\":[0-9]*" | head -1 | cut -d':' -f2
  fi
}

echo "===================================================================="
echo "    VERIFICANDO CONEXION A LA API ($BASE_URL)"
echo "===================================================================="
curl -s -o /dev/null -w "Estado del servidor: HTTP %{http_code}\n" "$BASE_URL/films"

# ------------------------------------------------------------------------------
# 1. STATUS
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 1. STATUS ] ===================="
echo "1.1. POST /statuses/ (Crear unitario)"
STATUS_RESP=$(curl -s -X POST "$BASE_URL/statuses/" \
  -H "Content-Type: application/json" \
  -d '{"statusname":"Pendiente"}')
echo "Respuesta: $STATUS_RESP"
STATUS_ID=$(extract_id "$STATUS_RESP" "statusid")
echo "-> statusid obtenido: $STATUS_ID"

echo -e "\n1.2. POST /statuses (Crear masivo)"
curl -s -X POST "$BASE_URL/statuses" \
  -H "Content-Type: application/json" \
  -d '[{"statusname":"Viendo"},{"statusname":"Visto"}]'
echo ""

echo -e "\n1.3. GET /statuses (Listar colección)"
curl -s "$BASE_URL/statuses"
echo ""

echo -e "\n1.4. GET /statuses/$STATUS_ID (Ver por ID)"
curl -s "$BASE_URL/statuses/$STATUS_ID"
echo ""

echo -e "\n1.5. PUT /statuses/$STATUS_ID (Actualizar)"
curl -s -X PUT "$BASE_URL/statuses/$STATUS_ID" \
  -H "Content-Type: application/json" \
  -d "{\"statusid\":$STATUS_ID,\"statusname\":\"Pendiente de ver\"}"
echo ""

# ------------------------------------------------------------------------------
# 2. GENRE
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 2. GENRE ] ===================="
echo "2.1. POST /genres/ (Crear unitario)"
GENRE_RESP=$(curl -s -X POST "$BASE_URL/genres/" \
  -H "Content-Type: application/json" \
  -d '{"genrename":"Accion"}')
echo "Respuesta: $GENRE_RESP"
GENRE_ID=$(extract_id "$GENRE_RESP" "genreid")
echo "-> genreid obtenido: $GENRE_ID"

echo -e "\n2.2. POST /genres (Crear masivo)"
curl -s -X POST "$BASE_URL/genres" \
  -H "Content-Type: application/json" \
  -d '[{"genrename":"Ciencia Ficcion"},{"genrename":"Aventura"}]'
echo ""

echo -e "\n2.3. GET /genres (Listar colección)"
curl -s "$BASE_URL/genres"
echo ""

echo -e "\n2.4. GET /genres/$GENRE_ID (Ver por ID)"
curl -s "$BASE_URL/genres/$GENRE_ID"
echo ""

echo -e "\n2.5. PUT /genres/$GENRE_ID (Actualizar)"
curl -s -X PUT "$BASE_URL/genres/$GENRE_ID" \
  -H "Content-Type: application/json" \
  -d "{\"genreid\":$GENRE_ID,\"genrename\":\"Accion y Suspenso\"}"
echo ""

# ------------------------------------------------------------------------------
# 3. CAST MEMBER
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 3. CAST MEMBER ] ===================="
echo "3.1. POST /castmembers/ (Crear unitario)"
CAST_RESP=$(curl -s -X POST "$BASE_URL/castmembers/" \
  -H "Content-Type: application/json" \
  -d '{"castmembername":"Keanu Reeves","castmemberrole":"A","castmembergender":"M"}')
echo "Respuesta: $CAST_RESP"
CAST_ID=$(extract_id "$CAST_RESP" "castmemberid")
echo "-> castmemberid obtenido: $CAST_ID"

echo -e "\n3.2. POST /castmembers (Crear masivo)"
curl -s -X POST "$BASE_URL/castmembers" \
  -H "Content-Type: application/json" \
  -d '[{"castmembername":"Laurence Fishburne","castmemberrole":"A","castmembergender":"M"}]'
echo ""

echo -e "\n3.3. GET /castmembers (Listar colección)"
curl -s "$BASE_URL/castmembers"
echo ""

echo -e "\n3.4. GET /castmembers/$CAST_ID (Ver por ID)"
curl -s "$BASE_URL/castmembers/$CAST_ID"
echo ""

echo -e "\n3.5. PUT /castmembers/$CAST_ID (Actualizar)"
curl -s -X PUT "$BASE_URL/castmembers/$CAST_ID" \
  -H "Content-Type: application/json" \
  -d "{\"castmemberid\":$CAST_ID,\"castmembername\":\"Keanu Charles Reeves\",\"castmemberrole\":\"A\",\"castmembergender\":\"M\"}"
echo ""

# ------------------------------------------------------------------------------
# 4. USER
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 4. USER ] ===================="
echo "4.1. POST /users/ (Crear unitario)"
USER_RESP=$(curl -s -X POST "$BASE_URL/users/" \
  -H "Content-Type: application/json" \
  -d '{"username":"neo","userpasswordhashed":"hash123","userbirthday":"1990-01-01T00:00:00Z","usergenre":"M"}')
echo "Respuesta: $USER_RESP"
USER_ID=$(extract_id "$USER_RESP" "userid")
echo "-> userid obtenido: $USER_ID"

echo -e "\n4.2. POST /users (Crear masivo)"
curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '[{"username":"trinity","userpasswordhashed":"hash456","userbirthday":"1992-05-12T00:00:00Z","usergenre":"F"}]'
echo ""

echo -e "\n4.3. GET /users (Listar colección)"
curl -s "$BASE_URL/users"
echo ""

echo -e "\n4.4. GET /users/$USER_ID (Ver por ID)"
curl -s "$BASE_URL/users/$USER_ID"
echo ""

echo -e "\n4.5. PUT /users/$USER_ID (Actualizar)"
curl -s -X PUT "$BASE_URL/users/$USER_ID" \
  -H "Content-Type: application/json" \
  -d "{\"userid\":$USER_ID,\"username\":\"thomas_anderson\",\"userpasswordhashed\":\"newhash123\",\"userbirthday\":\"1990-01-01T00:00:00Z\",\"usergenre\":\"M\"}"
echo ""

# ------------------------------------------------------------------------------
# 5. FILM
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 5. FILM ] ===================="
echo "5.1. POST /films/ (Crear unitario)"
FILM_RESP=$(curl -s -X POST "$BASE_URL/films/" \
  -H "Content-Type: application/json" \
  -d '{"filmtitle":"The Matrix","filmreleasedate":"1999-03-31T00:00:00Z","filmduration":136}')
echo "Respuesta: $FILM_RESP"
FILM_ID=$(extract_id "$FILM_RESP" "filmid")
echo "-> filmid obtenido: $FILM_ID"

echo -e "\n5.2. POST /films (Crear masivo)"
curl -s -X POST "$BASE_URL/films" \
  -H "Content-Type: application/json" \
  -d '[{"filmtitle":"The Matrix Reloaded","filmreleasedate":"2003-05-15T00:00:00Z","filmduration":138}]'
echo ""

echo -e "\n5.3. GET /films (Listar colección)"
curl -s "$BASE_URL/films"
echo ""

echo -e "\n5.4. GET /films/$FILM_ID (Ver por ID)"
curl -s "$BASE_URL/films/$FILM_ID"
echo ""

echo -e "\n5.5. PUT /films/$FILM_ID (Actualizar)"
curl -s -X PUT "$BASE_URL/films/$FILM_ID" \
  -H "Content-Type: application/json" \
  -d "{\"filmid\":$FILM_ID,\"filmtitle\":\"The Matrix (Remastered)\",\"filmreleasedate\":\"1999-03-31T00:00:00Z\",\"filmduration\":136}"
echo ""

# ------------------------------------------------------------------------------
# 6. FILMCAST (Clave compuesta: filmid, castmemberid - N:M sin PUT)
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 6. FILMCAST ] ===================="
echo "6.1. POST /filmcasts/ (Crear relación)"
curl -s -X POST "$BASE_URL/filmcasts/" \
  -H "Content-Type: application/json" \
  -d "{\"filmid\":$FILM_ID,\"castmemberid\":$CAST_ID}"
echo ""

echo -e "\n6.2. GET /filmcasts (Listar relaciones)"
curl -s "$BASE_URL/filmcasts"
echo ""

echo -e "\n6.3. GET /filmcasts/$FILM_ID/$CAST_ID (Ver relación por IDs compuestos)"
curl -s "$BASE_URL/filmcasts/$FILM_ID/$CAST_ID"
echo ""

echo -e "\n6.4. DELETE /filmcasts/$FILM_ID/$CAST_ID (Eliminar relación)"
curl -s -i -X DELETE "$BASE_URL/filmcasts/$FILM_ID/$CAST_ID" | head -n 1

# Volvemos a crear para dejar la relación disponible para pruebas posteriores
curl -s -X POST "$BASE_URL/filmcasts/" \
  -H "Content-Type: application/json" \
  -d "{\"filmid\":$FILM_ID,\"castmemberid\":$CAST_ID}" >/dev/null

# ------------------------------------------------------------------------------
# 7. FILMGENRE (Clave compuesta: genre_genreid, film_filmid - N:M sin PUT)
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 7. FILMGENRE ] ===================="
echo "7.1. POST /filmgenres/ (Crear relación)"
curl -s -X POST "$BASE_URL/filmgenres/" \
  -H "Content-Type: application/json" \
  -d "{\"genre_genreid\":$GENRE_ID,\"film_filmid\":$FILM_ID}"
echo ""

echo -e "\n7.2. GET /filmgenres (Listar relaciones)"
curl -s "$BASE_URL/filmgenres"
echo ""

echo -e "\n7.3. GET /filmgenres/$GENRE_ID/$FILM_ID (Ver relación por IDs compuestos)"
curl -s "$BASE_URL/filmgenres/$GENRE_ID/$FILM_ID"
echo ""

echo -e "\n7.4. DELETE /filmgenres/$GENRE_ID/$FILM_ID (Eliminar relación)"
curl -s -i -X DELETE "$BASE_URL/filmgenres/$GENRE_ID/$FILM_ID" | head -n 1

# Volvemos a crear para dejar la relación disponible
curl -s -X POST "$BASE_URL/filmgenres/" \
  -H "Content-Type: application/json" \
  -d "{\"genre_genreid\":$GENRE_ID,\"film_filmid\":$FILM_ID}" >/dev/null

# ------------------------------------------------------------------------------
# 8. FILMUSERSTATUS (Clave compuesta: film_filmid, user_userid + status_statusid)
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 8. FILMUSERSTATUS ] ===================="
echo "8.1. POST /filmuserstatuses/ (Crear relación)"
curl -s -X POST "$BASE_URL/filmuserstatuses/" \
  -H "Content-Type: application/json" \
  -d "{\"status_statusid\":$STATUS_ID,\"film_filmid\":$FILM_ID,\"user_userid\":$USER_ID}"
echo ""

echo -e "\n8.2. GET /filmuserstatuses (Listar relaciones)"
curl -s "$BASE_URL/filmuserstatuses"
echo ""

echo -e "\n8.3. GET /filmuserstatuses/$FILM_ID/$USER_ID (Ver por IDs compuestos)"
curl -s "$BASE_URL/filmuserstatuses/$FILM_ID/$USER_ID"
echo ""

echo -e "\n8.4. PUT /filmuserstatuses/$FILM_ID/$USER_ID (Actualizar estado)"
curl -s -X PUT "$BASE_URL/filmuserstatuses/$FILM_ID/$USER_ID" \
  -H "Content-Type: application/json" \
  -d "{\"status_statusid\":$STATUS_ID,\"film_filmid\":$FILM_ID,\"user_userid\":$USER_ID}"
echo ""

echo -e "\n8.5. DELETE /filmuserstatuses/$FILM_ID/$USER_ID (Eliminar relación)"
curl -s -i -X DELETE "$BASE_URL/filmuserstatuses/$FILM_ID/$USER_ID" | head -n 1

# ------------------------------------------------------------------------------
# 9. REVIEW (Clave compuesta: reviewfilmid, reviewuserid + estrellas y descripción)
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 9. REVIEW ] ===================="
echo "9.1. POST /reviews/ (Crear reseña)"
curl -s -X POST "$BASE_URL/reviews/" \
  -H "Content-Type: application/json" \
  -d "{\"reviewstars\":5,\"reviewdescription\":\"Excelente pelicula de ciencia ficcion\",\"reviewuserid\":$USER_ID,\"reviewfilmid\":$FILM_ID}"
echo ""

echo -e "\n9.2. GET /reviews (Listar reseñas)"
curl -s "$BASE_URL/reviews"
echo ""

echo -e "\n9.3. GET /reviews/$FILM_ID/$USER_ID (Ver reseña por IDs compuestos)"
curl -s "$BASE_URL/reviews/$FILM_ID/$USER_ID"
echo ""

echo -e "\n9.4. PUT /reviews/$FILM_ID/$USER_ID (Actualizar reseña)"
curl -s -X PUT "$BASE_URL/reviews/$FILM_ID/$USER_ID" \
  -H "Content-Type: application/json" \
  -d "{\"reviewstars\":4,\"reviewdescription\":\"Muy buena pelicula y efectos visuales\",\"reviewuserid\":$USER_ID,\"reviewfilmid\":$FILM_ID}"
echo ""

echo -e "\n9.5. DELETE /reviews/$FILM_ID/$USER_ID (Eliminar reseña)"
curl -s -i -X DELETE "$BASE_URL/reviews/$FILM_ID/$USER_ID" | head -n 1

# ------------------------------------------------------------------------------
# 10. LIMPIEZA DE ENTIDADES CREADAS (DELETE)
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 10. LIMPIEZA (DELETE) ] ===================="
echo "DELETE /filmcasts/$FILM_ID/$CAST_ID"
curl -s -i -X DELETE "$BASE_URL/filmcasts/$FILM_ID/$CAST_ID" | head -n 1

echo "DELETE /filmgenres/$GENRE_ID/$FILM_ID"
curl -s -i -X DELETE "$BASE_URL/filmgenres/$GENRE_ID/$FILM_ID" | head -n 1

echo "DELETE /films/$FILM_ID"
curl -s -i -X DELETE "$BASE_URL/films/$FILM_ID" | head -n 1

echo "DELETE /users/$USER_ID"
curl -s -i -X DELETE "$BASE_URL/users/$USER_ID" | head -n 1

echo "DELETE /castmembers/$CAST_ID"
curl -s -i -X DELETE "$BASE_URL/castmembers/$CAST_ID" | head -n 1

echo "DELETE /genres/$GENRE_ID"
curl -s -i -X DELETE "$BASE_URL/genres/$GENRE_ID" | head -n 1

echo "DELETE /statuses/$STATUS_ID"
curl -s -i -X DELETE "$BASE_URL/statuses/$STATUS_ID" | head -n 1


# INICIO CAMBIOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO

echo -e "\n===================================================================="
ERROR_COUNT=$(wc -l < "$ERROR_LOG_FILE" | tr -d ' ')
echo "    CANTIDAD DE ERRORES DETECTADOS: $ERROR_COUNT"
if [ "$ERROR_COUNT" -gt 0 ]; then
  echo "--------------------------------------------------------------------"
  echo "    DETALLE DE ERRORES:"
  while IFS= read -r line; do
    echo "      - $line"
  done < "$ERROR_LOG_FILE"
  echo "===================================================================="
  exit 1
else
  echo "    TODAS LAS PRUEBAS DE LA API SE COMPLETARON EXITOSAMENTE"
  echo "===================================================================="
fi


# FIN CAMBIOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO