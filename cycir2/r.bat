@REM cd ipe
@REM start /min cmd /c ipe.exe &
@REM start /min cmd /c MailHog_windows_386.exe &
@REM cd ..

go build -o cycir-api.exe ./cmd/api/.
@REM cycir-api -dbuser='postgres' -dbpass='qwerqwer' -esAddress="http://localhost:9200" -esUsername="elastic" -esPassword="EWAq+EaS8dyQV_82TSQd" -esIndex="my-index-000001" -pusherHost="localhost" -pusherPort="4001" -pusherSecret="123abc" -pusherKey="abc123" -pusherSecure=false -pusherApp="1" -db="temp"
cycir-api -dbuser='postgres' -dbpass='qwerqwer' -esAddress="http://localhost:9200" -esUsername="elastic" -esPassword="jHt1swMPzkAHttVEJ3si" -esIndex="cycir" -pusherHost="localhost" -pusherPort="4001" -pusherSecret="123abc" -pusherKey="abc123" -pusherSecure=false -pusherApp="1" -db="temp"