BINARY_NAME=app.exe

## build: builds all binaries
build:
	@go mod vendor
	@go build -o tmp/${BINARY_NAME} .
	@echo Django built!

run:
	@echo Staring Django...
	@start /min cmd /c tmp/${BINARY_NAME} &
	@echo Django started!

clean:
	@echo Cleaning...
	@DEL ${BINARY_NAME}
	@go clean
	@echo Cleaned!

test:
	@echo Testing...
	@go test ./...
	@echo Done!

start: run
	
stop:
	@echo "Starting the front end..."
	@taskkill /IM ${BINARY_NAME} /F
	@echo Stopped Django

restart: stop start
