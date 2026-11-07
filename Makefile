# Si el usuario no envia parametros, corre esto
.DEFAULT_GOAL := help

.PHONY: help start start-server stop-server clean-db test up-docker down-docker build-docker rebuild \
	check-sqlc generate-sqlc delete-sqlc generate-templ delete-templ copy-env clean

# Variables
APP_NAME := app
SQLC_GENERATED := $(shell ls ./db/sqlc 2>/dev/null || echo "")
# Usa templ del PATH, o el instalado con go install si no está en el PATH
TEMPL := $(shell command -v templ 2>/dev/null || echo $(HOME)/go/bin/templ)

# Chequea si está generado sqlc, sino lo genera
check-sqlc:
	@if [ -z "$(SQLC_GENERATED)" ]; then \
		echo "Generando SQLC"; \
		sqlc generate;  \
	else \
		echo "sqlc ya fue ejecutado"; \
	fi

# Elimina el código generado por sqlc
delete-sqlc:
	@echo "DEBUG delete-sqlc"
	rm -rf db/sqlc

# Elimina sqlc anterior y lo regenera desde cero
generate-sqlc: delete-sqlc
	@echo "Sqlc anterior eliminado"
	@$(MAKE) check-sqlc
	@echo "Nuevo sqlc generado"

# Elimina los archivos .go generados por templ en views/
delete-templ:
	@echo "DEBUG delete-templ"
	rm -f views/*_templ.go

# Elimina templ anterior y lo regenera desde cero
generate-templ: delete-templ
	@echo "DEBUG generate-templ"
	$(TEMPL) generate -path ./views
	@echo "Nuevo templ generado"

# Ejecuta los tests de Go en el contenedor (se borra al terminar)
test:
	@echo "DEBUG Testeo"
	docker compose run --rm api go test -v

# Levanta los contenedores en segundo plano
up-docker:
	@echo "DEBUG up-docker"
	docker compose up -d

# Construye las imágenes de Docker
build-docker:
	@echo "DEBUG build-docker"
	docker compose build

# Frena y elimina los contenedores
down-docker:
	@echo "DEBUG down-docker"
	docker compose down

# Elimina el binario y archivos temporales
clean:
	@echo "DEBUG clean"
	rm -f $(APP_NAME)
	rm -rf tmp/

# Baja los contenedores y los reconstruye
rebuild: down-docker
	@echo "DEBUG rebuild"
	docker compose up -d --build

# Copia las variables de entorno de ejemplo
copy-env:
	cp .env.example .env

# Limpia la BD, regenera sqlc y templ, levanta Docker, corre requests.hurl y apaga
start: copy-env generate-sqlc generate-templ clean-db
	@echo "Limpiando e iniciando Docker"
	@docker compose up --build -d
	@echo "Ejecutando tests de Hurl"
	@hurl --test requests.hurl
	@echo "Testeos finalizados"
	@echo "Limpiando y apagando"
	@docker compose down >/dev/null 2>&1

# Regenera templ y levanta el servidor en primer plano para pruebas manuales
start-server: copy-env check-sqlc generate-templ
	@echo "=== Iniciando Servidor de Pruebas ==="
	@docker compose up --build

# Frena el servidor sin borrar los datos de la BD
stop-server:
	@echo "=== Apagando Servidor ==="
	@docker compose down >/dev/null 2>&1
	@echo "=== Finalizado ==="

# Frena el servidor y borra los volúmenes de la BD para empezar de cero
clean-db:
	@echo "=== Limpiando Base de Datos (volúmenes) ==="
	@docker compose down -v >/dev/null 2>&1
	@echo "=== Base de Datos Limpia ==="

help:
	@echo "Posibles Comandos:"
	@echo ""
	@echo "  Uso general:"
	@echo "    help             Muestra esta ayuda"
	@echo "    start            Limpia BD, regenera sqlc y templ, levanta Docker, corre requests.hurl y apaga"
	@echo "    start-server     Regenera templ y levanta el servidor en http://localhost:8080"
	@echo "    stop-server      Frena el servidor (conserva datos de la BD)"
	@echo "    clean-db         Frena el servidor y borra los datos de la BD"
	@echo "    test             Ejecuta 'go test' dentro del contenedor"
	@echo ""
	@echo "  Docker:"
	@echo "    up-docker        Inicia los contenedores en segundo plano"
	@echo "    down-docker      Frena los contenedores"
	@echo "    build-docker     Construye las imágenes de Docker"
	@echo "    rebuild          Reconstruye y levanta los contenedores"
	@echo ""
	@echo "  Generación de código:"
	@echo "    check-sqlc       Genera el código de sqlc si no existe"
	@echo "    generate-sqlc    Borra y regenera el código de sqlc"
	@echo "    delete-sqlc      Elimina la carpeta db/sqlc"
	@echo "    generate-templ   Borra y regenera los *_templ.go de views/"
	@echo "    delete-templ     Elimina los *_templ.go de views/"
	@echo ""
	@echo "  Otros:"
	@echo "    copy-env         Copia .env.example a .env"
	@echo "    clean            Elimina binarios y temporales locales"
	@echo ""
