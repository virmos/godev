[![Go Report Card](https://goreportcard.com/badge/cycir)](https://goreportcard.com/report/cycir)

# cycir

## Build
 
On Windows:

~~~
go build -o cycir.exe cmd/web/.
~~~

Or for a particular platform:

~~~
env GOOS=linux GOARCH=amd64 go build -o cycir cmd/web/*.go
~~~

## Test

go test -coverprofile=coverage.out && go tool cover -html=coverage.out

go test -cover

## Requirements

cycir requires:
- Postgres 11 or later (db is set up as a repository, so other databases are possible)
- An account with [Pusher](https://pusher.com/), or a Pusher alternative 
(like [ipê](https://github.com/dimiro1/ipe))
