cd ipe
start /min cmd /c ipe.exe &
start /min cmd /c MailHog_windows_386.exe &
cd ..

go build -o cycir-api.exe ./cmd/api/.
cycir-api -dbuser='postgres' -dbpass='qwerqwer' -esAddress="http://localhost:9200" -esUsername="elastic" -esPassword="EWAq+EaS8dyQV_82TSQd" -esIndex="my-index-000001" -pusherHost="localhost" -pusherPort="4001" -pusherSecret="123abc" -pusherKey="abc123" -pusherSecure=false -pusherApp="1" -db="temp"
@REM cycir-api -dbuser='postgres' -dbpass='qwerqwer' -esAddress="http://localhost:9200" -esUsername="elastic" -esPassword="E08_6Tmn3QyzhI4XLxyR" -esIndex="cycir" -pusherHost="localhost" -pusherPort="4001" -pusherSecret="123abc" -pusherKey="abc123" -pusherSecure=false -pusherApp="1" -db="temp"
