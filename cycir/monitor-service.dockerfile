FROM alpine:latest

RUN mkdir /app

COPY monitorApp /app

CMD [ "/app/monitorApp"]