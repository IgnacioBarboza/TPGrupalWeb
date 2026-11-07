# ==============================================================================
# Prueba de los formularios (patron Post/Redirect/Get) de las 9 entidades:
#   POST /<entidades> valido   -> 303 con Location: /
#   POST /<entidades> invalido -> 400
#   GET  /                     -> 200 y el registro nuevo aparece en el HTML
# Los IDs se capturan del HTML de GET / para usarlos en las relaciones.
# Cada corrida usa un sufijo unico ({{run}}) en los nombres, asi el test
# funciona aunque la BD ya tenga datos de corridas anteriores.
# Ejecutar con: hurl --test requests.hurl
# ==============================================================================

# ------------------------------------------------------------------------------
# Espera a que la API este arriba (30 intentos, 1 por segundo)
# ------------------------------------------------------------------------------
GET http://localhost:8080/
[Options]
retry: 30
retry-interval: 1000
variable: uuid={{newUuid}}
HTTP 200
[Captures]
run: variable "uuid" regex /^([0-9a-f]{8})/


# ------------------------------------------------------------------------------
# STATUS
# ------------------------------------------------------------------------------
POST http://localhost:8080/statuses
[FormParams]
statusname: Pendiente {{run}}
HTTP 303
[Asserts]
header "Location" == "/"

# Nombre vacio (validacion)
POST http://localhost:8080/statuses
[FormParams]
statusname:
HTTP 400

GET http://localhost:8080/
HTTP 200
[Captures]
status_id: regex "\\[ID: (\\d+)\\] Pendiente {{run}}"
[Asserts]
body contains "Pendiente {{run}}"


# ------------------------------------------------------------------------------
# GENRE
# ------------------------------------------------------------------------------
POST http://localhost:8080/genres
[FormParams]
genrename: Accion {{run}}
HTTP 303
[Asserts]
header "Location" == "/"

# Nombre vacio (validacion)
POST http://localhost:8080/genres
[FormParams]
genrename:
HTTP 400

GET http://localhost:8080/
HTTP 200
[Captures]
genre_id: regex "\\[ID: (\\d+)\\] Accion {{run}}"
[Asserts]
body contains "Accion {{run}}"


# ------------------------------------------------------------------------------
# CAST MEMBER
# ------------------------------------------------------------------------------
POST http://localhost:8080/castmembers
[FormParams]
castmembername: Keanu Reeves {{run}}
castmemberrole: A
castmembergender: M
HTTP 303
[Asserts]
header "Location" == "/"

# Nombre vacio (validacion)
POST http://localhost:8080/castmembers
[FormParams]
castmembername:
castmemberrole: A
castmembergender: M
HTTP 400

# Rol de mas de 1 caracter (validacion)
POST http://localhost:8080/castmembers
[FormParams]
castmembername: Laurence Fishburne
castmemberrole: AB
castmembergender: M
HTTP 400

GET http://localhost:8080/
HTTP 200
[Captures]
cast_id: regex "\\[ID: (\\d+)\\] Keanu Reeves {{run}}"
[Asserts]
body contains "Keanu Reeves {{run}} - Rol: A - Género: M"


# ------------------------------------------------------------------------------
# USER
# ------------------------------------------------------------------------------
POST http://localhost:8080/users
[FormParams]
username: neo-{{run}}
userpasswordhashed: hash123
userbirthday: 1990-01-01
usergenre: M
HTTP 303
[Asserts]
header "Location" == "/"

# Username vacio (validacion)
POST http://localhost:8080/users
[FormParams]
username:
userpasswordhashed: hash123
userbirthday: 1990-01-01
usergenre: M
HTTP 400

# Fecha con formato invalido
POST http://localhost:8080/users
[FormParams]
username: trinity
userpasswordhashed: hash456
userbirthday: 12/05/1992
usergenre: F
HTTP 400

GET http://localhost:8080/
HTTP 200
[Captures]
user_id: regex "\\[ID: (\\d+)\\] neo-{{run}}"
[Asserts]
body contains "neo-{{run}} - Nacimiento: 1990-01-01 - Género: M"


# ------------------------------------------------------------------------------
# FILM
# ------------------------------------------------------------------------------
POST http://localhost:8080/films
[FormParams]
filmtitle: The Matrix {{run}}
filmreleasedate: 1999-03-31
filmduration: 136
HTTP 303
[Asserts]
header "Location" == "/"

# Titulo vacio (validacion)
POST http://localhost:8080/films
[FormParams]
filmtitle:
filmreleasedate: 1999-03-31
filmduration: 136
HTTP 400

# Duracion 0 (validacion)
POST http://localhost:8080/films
[FormParams]
filmtitle: X
filmreleasedate: 1999-03-31
filmduration: 0
HTTP 400

# Duracion no numerica
POST http://localhost:8080/films
[FormParams]
filmtitle: X
filmreleasedate: 1999-03-31
filmduration: abc
HTTP 400

GET http://localhost:8080/
HTTP 200
[Captures]
film_id: regex "\\[ID: (\\d+)\\] The Matrix {{run}}"
[Asserts]
body contains "The Matrix {{run}} - Estreno: 1999-03-31 - Duración: 136 min"


# ------------------------------------------------------------------------------
# FILMCAST (filmid, castmemberid)
# ------------------------------------------------------------------------------
POST http://localhost:8080/filmcasts
[FormParams]
filmid: {{film_id}}
castmemberid: {{cast_id}}
HTTP 303
[Asserts]
header "Location" == "/"

# ID 0 (validacion)
POST http://localhost:8080/filmcasts
[FormParams]
filmid: 0
castmemberid: {{cast_id}}
HTTP 400

GET http://localhost:8080/
HTTP 200
[Asserts]
body contains "Película ID: {{film_id}} ↔ Miembro Elenco ID: {{cast_id}}"


# ------------------------------------------------------------------------------
# FILMGENRE (genre_genreid, film_filmid)
# ------------------------------------------------------------------------------
POST http://localhost:8080/filmgenres
[FormParams]
genre_genreid: {{genre_id}}
film_filmid: {{film_id}}
HTTP 303
[Asserts]
header "Location" == "/"

# ID no numerico
POST http://localhost:8080/filmgenres
[FormParams]
genre_genreid: abc
film_filmid: {{film_id}}
HTTP 400

GET http://localhost:8080/
HTTP 200
[Asserts]
body contains "Género ID: {{genre_id}} ↔ Película ID: {{film_id}}"


# ------------------------------------------------------------------------------
# FILMUSERSTATUS (film_filmid, user_userid, status_statusid)
# ------------------------------------------------------------------------------
POST http://localhost:8080/filmuserstatuses
[FormParams]
film_filmid: {{film_id}}
user_userid: {{user_id}}
status_statusid: {{status_id}}
HTTP 303
[Asserts]
header "Location" == "/"

# ID 0 (validacion)
POST http://localhost:8080/filmuserstatuses
[FormParams]
film_filmid: {{film_id}}
user_userid: {{user_id}}
status_statusid: 0
HTTP 400

GET http://localhost:8080/
HTTP 200
[Asserts]
body contains "Película ID: {{film_id}} | Usuario ID: {{user_id}} → Estado ID: {{status_id}}"


# ------------------------------------------------------------------------------
# REVIEW (reviewfilmid, reviewuserid, estrellas y descripcion)
# ------------------------------------------------------------------------------
POST http://localhost:8080/reviews
[FormParams]
reviewfilmid: {{film_id}}
reviewuserid: {{user_id}}
reviewstars: 5
reviewdescription: Excelente pelicula de ciencia ficcion
HTTP 303
[Asserts]
header "Location" == "/"

# Estrellas fuera de rango (validacion)
POST http://localhost:8080/reviews
[FormParams]
reviewfilmid: {{film_id}}
reviewuserid: {{user_id}}
reviewstars: 99
reviewdescription: X
HTTP 400

# Descripcion vacia (validacion)
POST http://localhost:8080/reviews
[FormParams]
reviewfilmid: {{film_id}}
reviewuserid: {{user_id}}
reviewstars: 3
reviewdescription:
HTTP 400

GET http://localhost:8080/
HTTP 200
[Asserts]
body contains "[Película ID: {{film_id}} | Usuario ID: {{user_id}}] 5★"
body contains "Excelente pelicula de ciencia ficcion"
