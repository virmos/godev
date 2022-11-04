SHELL=cmd.exe
FRONT_END_BINARY=frontApp
BACK_END_BINARY=backApp
DSN="host=localhost port=5432 user=postgres password=qwerqwer dbname=temp sslmode=disable timezone=UTC connect_timeout=5"

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_front build_back
	@echo "Starting..."
	@echo Stopping docker images (if running...)
	docker-compose down
	@echo Building (when required) and starting docker images...
	docker-compose up --build -d
	@echo Docker images built and started!

## down: stop docker compose
down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

## build_back: builds the backend binary as a linux executable
build_back:
	@echo Building back end binary...
	chdir ..\cmd\api && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${BACK_END_BINARY} .
	@echo Done!

## build_front: builds the frontend binary
build_front:
	@echo Building front end binary...
	chdir ..\cmd\web && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${FRONT_END_BINARY} .
	@echo Done!

## start: starts the front end
start: build_front
	@echo Starting front end
	chdir ..\front-end && start /B ${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo Stopping front end...
	@taskkill /IM "${FRONT_END_BINARY}" /F
	@echo "Stopped front end!"