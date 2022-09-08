shell=cmd.exe
STRIPE_SECRET=123abc
STRIPE_KEY=abc123
GOSTRIPE_PORT=4000
API_PORT=4001
STRIPE_BINARY=gostripe_api.exe
STRIPE_API_BINARY=gostripe_api.exe

## build: builds all binaries
build: clean build_front build_back
	@echo All binaries built!

## clean: cleans all binaries and runs go clean
clean:
	@echo Cleaning...
	@echo y | DEL /S dist
	@go clean
	@echo Cleaned and deleted binaries

## build_front: builds the front end
build_front:
	@echo Building front end...
	@set STRIPE_KEY=${STRIPE_KEY} && set STRIPE_SECRET=${STRIPE_SECRET} && go build -o dist/gostripe.exe ./cmd/web
	@echo Front end built!

## build_back: builds the back end
build_back:
	@echo Building back end...
	@set STRIPE_KEY=${STRIPE_KEY} && set STRIPE_SECRET=${STRIPE_SECRET} && go build -o dist/gostripe_api.exe ./cmd/api
	@echo Back end built!

## start: starts front and back end
start: start_front start_back
	
## start_front: starts the front end
start_front: build_front
	@echo "Starting the front end..."
	@chdir dist && start /min cmd /c ${STRIPE_BINARY} /v -port=${GOSTRIPE_PORT}
	@echo "Front end running!"

## start_back: starts the back end
start_back: build_back
	@echo Starting the back end...
	@chdir dist && start /min cmd /c ${STRIPE_API_BINARY} /v -port=${API_PORT} 
	@echo "Back end running!"

## stop: stops the front and back end
stop: stop_front stop_back
	@echo "All applications stopped"

## stop_front: stops the front end
stop_front:
	@echo "Stopping the front end..."
	@taskkill /IM gostripe.exe /F
	@echo "Stopped front end"

## stop_back: stops the back end
stop_back:
	@echo "Stopping the back end..."
	@taskkill /IM gostripe_api.exe /F
	@echo "Stopped back end"