FROM alpine:latest

RUN mkdir /app
RUN mkdir /excel

COPY backApp /app

CMD [ "/app/backApp"]
