# Cargar variables desde .env
include .env
export $(shell sed 's/=.*//' .env)

# Si el usuario no envia parametros, corre esto
.DEFAULT_GOAL := help

# Variables
APP_NAME := app
SQLC_BIN := $(shell which sqlc 2>/dev/null || echo "")
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

SQLC_GENERATED := $(shell ls ./db/sqlc 2>/dev/null || echo "")

# Instalar sqlc si no existe
install-sqlc:
	@if [ -z "$(SQLC_BIN)" ]; then \
		echo "Instalando sqlc..."; \
		go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest; \
	else \
		echo "sqlc ya está instalado en $(SQLC_BIN)"; \
	fi

# Chequea si esta esta generado sqlc, sino lo genera
check-sqlc:
	@if [ -z "$(SQLC_GENERATED)" ]; then \
		echo "Generando SQLC"; \
		sqlc generate;  \
	else \
		echo "sqlc ya fue ejecutado"; \
	fi

# Ejecuta install-sqlc y luego check-sqlc
generatesqlc: install-sqlc check-sqlc
	@echo "DEBUG generatesqlc terminado"

# Hacer testeos (rm = borrar luego de ejecutar)
test:
	@echo "DEBUG Testeo"
	docker compose run --rm api go test -v

# Empieza el container
upDocker:
	@echo "DEBUG upDocker"
	docker compose up -d

buildDocker:
	@echo "DEBUG buildDocker"
	docker compose build

# Frena el container
downDocker:
	@echo "DEBUG downDocker"
	docker compose down

# Limpiar volumenes
cleanVolumenes:
	@echo "DEBUG cleanVolumenes"
	docker compose down -v

# Limpiar archivos
clean:
	@echo "DEBUG clean"
	rm -f $(APP_NAME)
	rm -rf tmp/

# Baja el container y lo re-construye
rebuild: downDocker
	@echo "DEBUG rebuild"
	docker compose up -d --build    

# Limpia lo viejo, construye, levanta, testea y limpia
start: cleanVolumenes generatesqlc buildDocker upDocker
	@echo "Ejecutando tests..."
	make test
	@echo "Tests exitosos, procediendo a limpiar..."
	make cleanVolumenes
	make clean
	@echo "DEBUG start completado"

help:
	@echo "Posibles Comandos:"
	@echo ""
	@echo "  help             Muestra este mensaje de ayuda"
	@echo "  start            Construye, levanta el contenedor y ejecuta los tests"
	@echo "  test             Ejecuta los tests dentro del contenedor (se borra al terminar)"
	@echo "  upDocker         Inicia los contenedores"
	@echo "  buildDocker      Construye las imágenes de Docker"
	@echo "  downDocker       Frena los contenedores"
	@echo "  cleanVolumenes   Frena los contenedores y limpia los volúmenes"
	@echo "  clean            Elimina el binario y los archivos temporales"
	@echo "  rebuild          Detiene y reconstruye los contenedores desde cero"
	@echo "  check-sqlc       Verifica si el código está generado y ejecuta sqlc generate si no lo está"
	@echo "  install-sqlc     Instala sqlc en tu sistema si no lo tenés"
	@echo "  generatesqlc     Instala sqlc (si hace falta) y genera el código"
	@echo ""