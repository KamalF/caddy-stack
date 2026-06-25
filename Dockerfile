FROM golang:1.26 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -trimpath -o /caddy .

FROM caddy:2.11.4
COPY --from=builder /caddy /usr/bin/caddy
CMD ["caddy", "docker-proxy"]
