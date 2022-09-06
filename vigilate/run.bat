cd ipe
start /min cmd /c ipe.exe &
cd ..
go build -o vigilate.exe ./cmd/web/.
docker compose up -d
vigilate -dbuser='postgres' -dbpass='qwerqwer' -pusherHost="localhost" -pusherPort="4001" -pusherSecret="123abc" -pusherKey="abc123" -pusherSecure=false -pusherApp="1" -db="vigilate" -redisHost="localhost:6379" -redisPrefix="vigilate"