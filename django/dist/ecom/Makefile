BINARY_NAME=ecom.exe

## build: builds all binaries
build:
	@go mod vendor
	@go build -o ${BINARY_NAME} .
	@echo App built!

run:
	@echo Staring App...
	@start /min cmd /c ${BINARY_NAME} &
	@echo App started!

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
	@echo Stopped App

restart: stop start