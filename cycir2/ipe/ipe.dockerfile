FROM alpine:latest

RUN mkdir /app

COPY ipe /app
COPY config.yml /app

CMD [ "/app/ipe"]