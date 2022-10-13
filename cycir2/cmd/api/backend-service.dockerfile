FROM alpine:latest

RUN mkdir /app

COPY backApp /app

CMD [ "/app/backApp"]