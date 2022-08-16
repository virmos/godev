/*
docker run --name pgacid -d POSTGRES_PASSWORD=postgres postgres:13
docker exec -it pgacid psql -U postgres
*/

begin transaction isolation level repeatable read;
