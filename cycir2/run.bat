@REM cd ipe
@REM start /min cmd /c ipe.exe &
@REM cd ..
go build -o cycir.exe ./cmd/web/.
@REM docker compose up -d
@REM soda migrate
cycir -dbuser='postgres' -dbpass='qwerqwer' -pusherHost="localhost" -pusherPort="4001" -pusherSecret="123abc" -pusherKey="abc123" -pusherSecure=false -pusherApp="1" -db="temp"

@REM cycir -dbuser='postgres' -dbpass='qwerqwer' -pusherHost="localhost" -pusherPort="4001" -pusherSecret="123abc" -pusherKey="abc123" -pusherSecure=false -pusherApp="1" -db="temp" -redisHost="localhost:6379" -redisPrefix="cycir"
