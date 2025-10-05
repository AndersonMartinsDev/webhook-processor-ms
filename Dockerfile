

FROM golang:1.24.5-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

# COPIA A PASTA MIGRATIONS DO ESTÃGIO DE BUILD PARA O ESTÃGIO FINAL
# A pasta 'migrations' estÃ¡ em /app/migrations no estÃ¡gio 'builder'
COPY --from=builder /app/cmd/api/migrations ./cmd/api/migrations

EXPOSE 50051

CMD ["./main"]