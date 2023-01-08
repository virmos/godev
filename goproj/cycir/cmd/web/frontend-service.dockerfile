FROM alpine:latest

WORKDIR /app

COPY static ./static
COPY views ./cmd/web/views
COPY frontApp ./

CMD [ "/app/frontApp"]
