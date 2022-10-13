FROM alpine:latest

RUN mkdir /app

COPY frontApp.exe /app

CMD [ "/app/frontApp.exe"]