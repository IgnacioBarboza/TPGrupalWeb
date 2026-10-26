# Si el usuario no envia parametros, corre esto
.DEFAULT_GOAL := help

# Variables
APP_NAME := app
SQLC_GENERATED := $(shell ls ./db/sqlc 2>/dev/null || echo "")

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
generatesqlc: delete-sqlc 
	@echo "Sqlc anterior eliminado"
	@make check-sqlc
	@echo "Nuevo sqlc generado"

# Ejecuta los tests de Go en el contenedor (se borra al terminar)
test:
	@echo "DEBUG Testeo"
	docker compose run --rm api go test -v

# Levanta los contenedores en segundo plano
upDocker:
	@echo "DEBUG upDocker"
	docker compose up -d

# Construye las imágenes de Docker
buildDocker:
	@echo "DEBUG buildDocker"
	docker compose build

# Frena y elimina los contenedores
downDocker:
	@echo "DEBUG downDocker"
	docker compose down

# Frena los contenedores y elimina los volúmenes
cleanVolumenes:
	@echo "DEBUG cleanVolumenes"
	docker compose down -v

# Elimina el binario y archivos temporales
clean:
	@echo "DEBUG clean"
	rm -f $(APP_NAME)
	rm -rf tmp/

# Baja los contenedores y los reconstruye
rebuild: downDocker
	@echo "DEBUG rebuild"
	docker compose up -d --build    

# Copia las variables de entorno de ejemplo
copyEnv:
	cp .env.example .env

# Limpia la BD, regenera sqlc, levanta Docker, corre requests.bash y apaga
start: copyEnv generatesqlc cleanDB
	@echo "Limpiando e iniciando Docker"
	@docker compose up --build -d 
	@echo "Ejecutando batch de testing de Curls"
	@./requests.bash
	@echo "Testeos finalizados"
	@echo "Limpiando y apagando"
	@docker compose down >/dev/null 2>&1

# Levanta el servidor en primer plano para pruebas manuales
StartServer: copyEnv check-sqlc
	@echo "\n=== Iniciando Servidor de Pruebas ==="
	@docker compose up --build

# Frena el servidor sin borrar los datos de la BD
StopServer:
	@echo "\n=== Apagando Servidor ==="
	@docker compose down >/dev/null 2>&1
	@echo "=== Finalizado ==="

# Frena el servidor y borra los volúmenes de la BD para empezar de cero
cleanDB:
	@echo "\n=== Limpiando Base de Datos (volúmenes) ==="
	@docker compose down -v >/dev/null 2>&1
	@echo "=== Base de Datos Limpia ==="


help:
	@echo "Posibles Comandos:"
	@echo ""
	@echo "  help             Muestra esta ayuda"
	@echo "  start            Limpia BD, levanta Docker, corre requests.bash y apaga"
	@echo "  StartServer      Levanta el servidor para pruebas manuales"
	@echo "  StopServer       Frena el servidor (conserva datos de la BD)"
	@echo "  cleanDB          Frena el servidor y borra los datos de la BD"
	@echo "  test             Ejecuta 'go test' dentro del contenedor"
	@echo "  upDocker         Inicia los contenedores en segundo plano"
	@echo "  downDocker       Frena los contenedores"
	@echo "  buildDocker      Construye las imágenes de Docker"
	@echo "  rebuild          Reconstruye y levanta los contenedores"
	@echo "  check-sqlc       Genera el código de sqlc si no existe"
	@echo "  generatesqlc     Borra y regenera el código de sqlc"
	@echo "  delete-sqlc      Elimina la carpeta db/sqlc"
	@echo "  copyEnv          Copia .env.example a .env"
	@echo "  clean            Elimina binarios y temporales locales"
	@echo ""