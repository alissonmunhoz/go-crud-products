# ---------- STAGE 1: Builder ----------
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Informe o diretório/arquivo onde está o package main
# Exemplos de valores válidos:
#   .                  (se o main.go está na raiz)
#   ./cmd/api          (se o main está em cmd/api)
#   ./server           (se o main está em server/)
#   ./cmd/api/main.go  (apontando direto pro arquivo)
ARG APP_PATH=.
RUN test -e $APP_PATH || (echo "APP_PATH não encontrado: $APP_PATH" && ls -la && exit 1)

# Compila
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server $APP_PATH

# ---------- STAGE 2: Runtime ----------
FROM alpine:3.20
RUN addgroup -S app && adduser -S app -G app
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/server /app/server

ENV DB_HOST=mysql \
    DB_PORT=3306 \
    DB_USER=root \
    DB_PASSWORD=root \
    DB_NAME=products

EXPOSE 8080
USER app
CMD ["./server"]
