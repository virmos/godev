cd ipe
start /min cmd /c ipe.exe &
cd ..

go build -o cycir.exe ./cmd/api/.
cycir -dbuser='postgres' -dbpass='qwerqwer' -pusherHost="localhost" -pusherPort="4001" -pusherSecret="123abc" -pusherKey="abc123" -pusherSecure=false -pusherApp="1" -db="temp"