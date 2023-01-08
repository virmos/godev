FROM alpine:latest

WORKDIR /app

COPY database.yml ./
COPY migrationApp ./

CMD [ "/app/migrationApp"]
