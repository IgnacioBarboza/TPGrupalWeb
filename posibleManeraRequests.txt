#!/usr/bin/env bash
# ==============================================================================
# Prueba de todos los endpoints de la API, verificando el codigo HTTP esperado
# segun el enunciado:
#   POST   -> 201 (y 400 si la validacion falla)
#   GET    -> 200 (lista y por id), 404 si no existe
#   PUT    -> 200 (y 404 si no existe)
#   DELETE -> 204
# Cada request imprime OK o FALLO. Al final se listan los fallos y se sale con
# exit 1 si hubo alguno.
# ==============================================================================

BASE_URL="http://localhost:8080"
FAIL_FILE=$(mktemp)
trap 'rm -f "$FAIL_FILE"' EXIT
TOTAL=0
BODY=""

# ------------------------------------------------------------------------------
# Helpers
# ------------------------------------------------------------------------------

# req CODIGO_ESPERADO "descripcion" METODO RUTA [JSON]
# Deja la respuesta en $BODY. Llamarla siempre directo (no dentro de $(...)).
req() {
  local expected="$1"
  local desc="$2"
  local method="$3"
  local path="$4"
  local data="${5:-}"
  local out
  local code

  if [ -n "$data" ]; then
    out=$(command curl -s -w $'\n%{http_code}' -X "$method" "$BASE_URL$path" \
      -H "Content-Type: application/json" -d "$data")
  else
    out=$(command curl -s -w $'\n%{http_code}' -X "$method" "$BASE_URL$path")
  fi

  code="${out##*$'\n'}"
  BODY="${out%$'\n'*}"
  TOTAL=$((TOTAL + 1))

  if [ "$code" = "$expected" ]; then
    echo "   OK    [$code] $desc  ($method $path)"
  else
    echo "   FALLO [$code, esperaba $expected] $desc  ($method $path)"
    echo "[$code, esperaba $expected] $desc ($method $path)" >> "$FAIL_FILE"
  fi
  if [ -n "$BODY" ]; then
    echo "         -> $BODY"
  fi
}

# Extrae un valor numerico de un JSON
extract_id() {
  local json="$1"
  local key="$2"
  if command -v jq >/dev/null 2>&1; then
    echo "$json" | jq -r ".$key // empty"
  else
    echo "$json" | grep -o "\"$key\":[0-9]*" | head -1 | cut -d':' -f2
  fi
}

# Corta el script si no se pudo obtener un ID (sino las requests siguientes
# se arman con URLs rotas)
require_id() {
  local name="$1"
  local value="$2"
  if [ -z "$value" ]; then
    echo ""
    echo "No se obtuvo $name. Abortando: las pruebas siguientes dependen de este ID."
    exit 1
  fi
}

# Espera a que la API responda; corta si no lo hace
wait_for_api() {
  local attempts=30
  local code
  local i
  for ((i = 1; i <= attempts; i++)); do
    code=$(command curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/films")
    if [ "$code" = "200" ]; then
      echo "API lista (intento $i)"
      return 0
    fi
    sleep 1
  done
  echo "La API no respondio tras $attempts intentos. Abortando."
  exit 1
}

echo "===================================================================="
echo "    VERIFICANDO CONEXION A LA API ($BASE_URL)"
echo "===================================================================="
wait_for_api

# ------------------------------------------------------------------------------
# 1. STATUS
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 1. STATUS ] ===================="
req 201 "POST unitario" POST "/statuses/" '{"statusname":"Pendiente"}'
STATUS_ID=$(extract_id "$BODY" "statusid")
require_id "statusid" "$STATUS_ID"

req 400 "POST con nombre vacio (validacion)" POST "/statuses/" '{"statusname":""}'
req 400 "POST con JSON invalido" POST "/statuses/" '{"statusname":'
req 201 "POST masivo" POST "/statuses" '[{"statusname":"Viendo"},{"statusname":"Visto"}]'
req 200 "GET lista" GET "/statuses"
req 200 "GET por id" GET "/statuses/$STATUS_ID"
req 404 "GET id inexistente" GET "/statuses/999999"
req 400 "GET id no numerico" GET "/statuses/abc"
req 200 "PUT actualizar" PUT "/statuses/$STATUS_ID" \
  "{\"statusid\":$STATUS_ID,\"statusname\":\"Pendiente de ver\"}"
req 404 "PUT id inexistente" PUT "/statuses/999999" \
  '{"statusid":999999,"statusname":"X"}'
req 400 "PUT con nombre vacio (validacion)" PUT "/statuses/$STATUS_ID" \
  "{\"statusid\":$STATUS_ID,\"statusname\":\"\"}"

# ------------------------------------------------------------------------------
# 2. GENRE
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 2. GENRE ] ===================="
req 201 "POST unitario" POST "/genres/" '{"genrename":"Accion"}'
GENRE_ID=$(extract_id "$BODY" "genreid")
require_id "genreid" "$GENRE_ID"

req 400 "POST con nombre vacio (validacion)" POST "/genres/" '{"genrename":""}'
req 400 "POST con JSON invalido" POST "/genres/" '{"genrename":'
req 201 "POST masivo" POST "/genres" '[{"genrename":"Ciencia Ficcion"},{"genrename":"Aventura"}]'
req 200 "GET lista" GET "/genres"
req 200 "GET por id" GET "/genres/$GENRE_ID"
req 404 "GET id inexistente" GET "/genres/999999"
req 400 "GET id no numerico" GET "/genres/abc"
req 200 "PUT actualizar" PUT "/genres/$GENRE_ID" \
  "{\"genreid\":$GENRE_ID,\"genrename\":\"Accion y Suspenso\"}"
req 404 "PUT id inexistente" PUT "/genres/999999" \
  '{"genreid":999999,"genrename":"X"}'
req 400 "PUT con nombre vacio (validacion)" PUT "/genres/$GENRE_ID" \
  "{\"genreid\":$GENRE_ID,\"genrename\":\"\"}"

# ------------------------------------------------------------------------------
# 3. CAST MEMBER
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 3. CAST MEMBER ] ===================="
req 201 "POST unitario" POST "/castmembers/" \
  '{"castmembername":"Keanu Reeves","castmemberrole":"A","castmembergender":"M"}'
CAST_ID=$(extract_id "$BODY" "castmemberid")
require_id "castmemberid" "$CAST_ID"

req 400 "POST con nombre vacio (validacion)" POST "/castmembers/" \
  '{"castmembername":"","castmemberrole":"A","castmembergender":"M"}'
req 400 "POST con JSON invalido" POST "/castmembers/" '{"castmembername":'
req 201 "POST masivo" POST "/castmembers" \
  '[{"castmembername":"Laurence Fishburne","castmemberrole":"A","castmembergender":"M"}]'
req 200 "GET lista" GET "/castmembers"
req 200 "GET por id" GET "/castmembers/$CAST_ID"
req 404 "GET id inexistente" GET "/castmembers/999999"
req 400 "GET id no numerico" GET "/castmembers/abc"
req 200 "PUT actualizar" PUT "/castmembers/$CAST_ID" \
  "{\"castmemberid\":$CAST_ID,\"castmembername\":\"Keanu Charles Reeves\",\"castmemberrole\":\"A\",\"castmembergender\":\"M\"}"
req 404 "PUT id inexistente" PUT "/castmembers/999999" \
  '{"castmemberid":999999,"castmembername":"X","castmemberrole":"A","castmembergender":"M"}'
req 400 "PUT con nombre vacio (validacion)" PUT "/castmembers/$CAST_ID" \
  "{\"castmemberid\":$CAST_ID,\"castmembername\":\"\",\"castmemberrole\":\"A\",\"castmembergender\":\"M\"}"

# ------------------------------------------------------------------------------
# 4. USER
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 4. USER ] ===================="
req 201 "POST unitario" POST "/users/" \
  '{"username":"neo","userpasswordhashed":"hash123","userbirthday":"1990-01-01T00:00:00Z","usergenre":"M"}'
USER_ID=$(extract_id "$BODY" "userid")
require_id "userid" "$USER_ID"

req 400 "POST con username vacio (validacion)" POST "/users/" \
  '{"username":"","userpasswordhashed":"hash123","userbirthday":"1990-01-01T00:00:00Z","usergenre":"M"}'
req 400 "POST con JSON invalido" POST "/users/" '{"username":'
req 201 "POST masivo" POST "/users" \
  '[{"username":"trinity","userpasswordhashed":"hash456","userbirthday":"1992-05-12T00:00:00Z","usergenre":"F"}]'
req 200 "GET lista" GET "/users"
req 200 "GET por id" GET "/users/$USER_ID"
req 404 "GET id inexistente" GET "/users/999999"
req 400 "GET id no numerico" GET "/users/abc"
req 200 "PUT actualizar" PUT "/users/$USER_ID" \
  "{\"userid\":$USER_ID,\"username\":\"thomas_anderson\",\"userpasswordhashed\":\"newhash123\",\"userbirthday\":\"1990-01-01T00:00:00Z\",\"usergenre\":\"M\"}"
req 404 "PUT id inexistente" PUT "/users/999999" \
  '{"userid":999999,"username":"x","userpasswordhashed":"h","userbirthday":"1990-01-01T00:00:00Z","usergenre":"M"}'
req 400 "PUT con username vacio (validacion)" PUT "/users/$USER_ID" \
  "{\"userid\":$USER_ID,\"username\":\"\",\"userpasswordhashed\":\"newhash123\",\"userbirthday\":\"1990-01-01T00:00:00Z\",\"usergenre\":\"M\"}"

# ------------------------------------------------------------------------------
# 5. FILM
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 5. FILM ] ===================="
req 201 "POST unitario" POST "/films/" \
  '{"filmtitle":"The Matrix","filmreleasedate":"1999-03-31T00:00:00Z","filmduration":136}'
FILM_ID=$(extract_id "$BODY" "filmid")
require_id "filmid" "$FILM_ID"

req 400 "POST con titulo vacio (validacion)" POST "/films/" \
  '{"filmtitle":"","filmreleasedate":"1999-03-31T00:00:00Z","filmduration":136}'
req 400 "POST con duracion invalida (validacion)" POST "/films/" \
  '{"filmtitle":"X","filmreleasedate":"1999-03-31T00:00:00Z","filmduration":0}'
req 400 "POST con JSON invalido" POST "/films/" '{"filmtitle":'
req 201 "POST masivo" POST "/films" \
  '[{"filmtitle":"The Matrix Reloaded","filmreleasedate":"2003-05-15T00:00:00Z","filmduration":138}]'
req 400 "POST masivo con un elemento invalido" POST "/films" \
  '[{"filmtitle":"Ok","filmreleasedate":"2003-05-15T00:00:00Z","filmduration":100},{"filmtitle":"","filmreleasedate":"2003-05-15T00:00:00Z","filmduration":100}]'
req 200 "GET lista" GET "/films"
req 200 "GET por id" GET "/films/$FILM_ID"
req 404 "GET id inexistente" GET "/films/999999"
req 400 "GET id no numerico" GET "/films/abc"
req 200 "PUT actualizar" PUT "/films/$FILM_ID" \
  "{\"filmid\":$FILM_ID,\"filmtitle\":\"The Matrix (Remastered)\",\"filmreleasedate\":\"1999-03-31T00:00:00Z\",\"filmduration\":136}"
req 404 "PUT id inexistente" PUT "/films/999999" \
  '{"filmid":999999,"filmtitle":"X","filmreleasedate":"1999-03-31T00:00:00Z","filmduration":1}'
req 400 "PUT con titulo vacio (validacion)" PUT "/films/$FILM_ID" \
  "{\"filmid\":$FILM_ID,\"filmtitle\":\"\",\"filmreleasedate\":\"1999-03-31T00:00:00Z\",\"filmduration\":136}"
req 400 "PUT con id distinto en URL y body" PUT "/films/$FILM_ID" \
  '{"filmid":999999,"filmtitle":"X","filmreleasedate":"1999-03-31T00:00:00Z","filmduration":1}'

# ------------------------------------------------------------------------------
# 6. FILMCAST (clave compuesta: filmid, castmemberid)
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 6. FILMCAST ] ===================="
req 201 "POST relacion" POST "/filmcasts/" \
  "{\"filmid\":$FILM_ID,\"castmemberid\":$CAST_ID}"
req 200 "GET lista" GET "/filmcasts"
req 200 "GET por ids compuestos" GET "/filmcasts/$FILM_ID/$CAST_ID"
req 404 "GET relacion inexistente" GET "/filmcasts/999999/999999"
req 204 "DELETE relacion" DELETE "/filmcasts/$FILM_ID/$CAST_ID"
req 404 "GET tras el DELETE" GET "/filmcasts/$FILM_ID/$CAST_ID"
# Se vuelve a crear para que la limpieza final tenga algo que borrar
req 201 "POST relacion (recrear)" POST "/filmcasts/" \
  "{\"filmid\":$FILM_ID,\"castmemberid\":$CAST_ID}"

# ------------------------------------------------------------------------------
# 7. FILMGENRE (clave compuesta: genre_genreid, film_filmid)
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 7. FILMGENRE ] ===================="
req 201 "POST relacion" POST "/filmgenres/" \
  "{\"genre_genreid\":$GENRE_ID,\"film_filmid\":$FILM_ID}"
req 200 "GET lista" GET "/filmgenres"
req 200 "GET por ids compuestos" GET "/filmgenres/$GENRE_ID/$FILM_ID"
req 404 "GET relacion inexistente" GET "/filmgenres/999999/999999"
req 204 "DELETE relacion" DELETE "/filmgenres/$GENRE_ID/$FILM_ID"
req 404 "GET tras el DELETE" GET "/filmgenres/$GENRE_ID/$FILM_ID"
req 201 "POST relacion (recrear)" POST "/filmgenres/" \
  "{\"genre_genreid\":$GENRE_ID,\"film_filmid\":$FILM_ID}"

# ------------------------------------------------------------------------------
# 8. FILMUSERSTATUS (film_filmid, user_userid + status_statusid)
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 8. FILMUSERSTATUS ] ===================="
req 201 "POST relacion" POST "/filmuserstatuses/" \
  "{\"status_statusid\":$STATUS_ID,\"film_filmid\":$FILM_ID,\"user_userid\":$USER_ID}"
req 200 "GET lista" GET "/filmuserstatuses"
req 200 "GET por ids compuestos" GET "/filmuserstatuses/$FILM_ID/$USER_ID"
req 404 "GET relacion inexistente" GET "/filmuserstatuses/999999/999999"
req 200 "PUT actualizar estado" PUT "/filmuserstatuses/$FILM_ID/$USER_ID" \
  "{\"status_statusid\":$STATUS_ID,\"film_filmid\":$FILM_ID,\"user_userid\":$USER_ID}"
req 204 "DELETE relacion" DELETE "/filmuserstatuses/$FILM_ID/$USER_ID"
req 404 "GET tras el DELETE" GET "/filmuserstatuses/$FILM_ID/$USER_ID"

# ------------------------------------------------------------------------------
# 9. REVIEW (reviewfilmid, reviewuserid + estrellas y descripcion)
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 9. REVIEW ] ===================="
req 201 "POST resena" POST "/reviews/" \
  "{\"reviewstars\":5,\"reviewdescription\":\"Excelente pelicula de ciencia ficcion\",\"reviewuserid\":$USER_ID,\"reviewfilmid\":$FILM_ID}"
req 400 "POST con estrellas fuera de rango (validacion)" POST "/reviews/" \
  "{\"reviewstars\":99,\"reviewdescription\":\"X\",\"reviewuserid\":$USER_ID,\"reviewfilmid\":$FILM_ID}"
req 200 "GET lista" GET "/reviews"
req 200 "GET por ids compuestos" GET "/reviews/$FILM_ID/$USER_ID"
req 404 "GET resena inexistente" GET "/reviews/999999/999999"
req 200 "PUT actualizar" PUT "/reviews/$FILM_ID/$USER_ID" \
  "{\"reviewstars\":4,\"reviewdescription\":\"Muy buena pelicula y efectos visuales\",\"reviewuserid\":$USER_ID,\"reviewfilmid\":$FILM_ID}"
req 204 "DELETE resena" DELETE "/reviews/$FILM_ID/$USER_ID"
req 404 "GET tras el DELETE" GET "/reviews/$FILM_ID/$USER_ID"

# ------------------------------------------------------------------------------
# 10. LIMPIEZA (DELETE) en orden inverso a las dependencias
# ------------------------------------------------------------------------------
echo -e "\n==================== [ 10. LIMPIEZA (DELETE) ] ===================="
req 204 "DELETE filmcast" DELETE "/filmcasts/$FILM_ID/$CAST_ID"
req 204 "DELETE filmgenre" DELETE "/filmgenres/$GENRE_ID/$FILM_ID"
req 204 "DELETE film" DELETE "/films/$FILM_ID"
req 404 "GET film tras el DELETE" GET "/films/$FILM_ID"
req 204 "DELETE user" DELETE "/users/$USER_ID"
req 404 "GET user tras el DELETE" GET "/users/$USER_ID"
req 204 "DELETE castmember" DELETE "/castmembers/$CAST_ID"
req 404 "GET castmember tras el DELETE" GET "/castmembers/$CAST_ID"
req 204 "DELETE genre" DELETE "/genres/$GENRE_ID"
req 404 "GET genre tras el DELETE" GET "/genres/$GENRE_ID"
req 204 "DELETE status" DELETE "/statuses/$STATUS_ID"
req 404 "GET status tras el DELETE" GET "/statuses/$STATUS_ID"

# ------------------------------------------------------------------------------
# Resumen
# ------------------------------------------------------------------------------
echo -e "\n===================================================================="
FAIL_COUNT=$(wc -l < "$FAIL_FILE" | tr -d ' ')
echo "    REQUESTS EJECUTADAS: $TOTAL"
echo "    FALLOS: $FAIL_COUNT"
if [ "$FAIL_COUNT" -gt 0 ]; then
  echo "--------------------------------------------------------------------"
  echo "    DETALLE:"
  while IFS= read -r line; do
    echo "      - $line"
  done < "$FAIL_FILE"
  echo "===================================================================="
  exit 1
else
  echo "    TODAS LAS PRUEBAS COINCIDEN CON LO ESPERADO"
  echo "===================================================================="
fi