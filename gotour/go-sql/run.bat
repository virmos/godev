go build -o cycir.exe ./cmd/web/.
cycir -dbuser='postgres' -dbpass='qwerqwer' -db="cycir"
@REM cycir -dbuser='postgres' -dbpass='qwerqwer' -pusherHost='localhost:4001' -pusherSecret='somesecret' -pusherKey='somekey' -pusherSecure=false pusherApp="1" -db="cycir"
