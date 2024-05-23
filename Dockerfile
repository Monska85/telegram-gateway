FROM golang:1.24-alpine3.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY lib ./lib

RUN CGO_ENABLED=0 GOOS=linux go build -o /telegram-gateway

FROM scratch

# Copy certificates from alpine
COPY --from=alpine:3.22 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the statically compiled Go binary
COPY --from=build /telegram-gateway /telegram-gateway

EXPOSE 8080

ENTRYPOINT ["/telegram-gateway"]
