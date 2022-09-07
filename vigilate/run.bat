cd ipe
start /min cmd /c ipe.exe &
cd ..
go build -o vigilate.exe ./cmd/web/.
<<<<<<< HEAD
docker compose up -d
vigilate -dbuser='postgres' -dbpass='qwerqwer' -pusherHost="localhost" -pusherPort="4001" -pusherSecret="123abc" -pusherKey="abc123" -pusherSecure=false -pusherApp="1" -db="vigilate" -redisHost="localhost" -redisPort="6379"
=======
@REM docker compose up -d
soda migrate
vigilate -dbuser='postgres' -dbpass='qwerqwer' -pusherHost="localhost" -pusherPort="4001" -pusherSecret="123abc" -pusherKey="abc123" -pusherSecure=false -pusherApp="1" -db="temp" -redisHost="localhost:6379" -redisPrefix="vigilate"
>>>>>>> 71f64a59cec07db2166a805947a4839fd2e9d44c
