FROM golang:1.20.7-alpine3.18

WORKDIR /app

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE 6969

# Run
CMD ["/docker-gs-ping"]