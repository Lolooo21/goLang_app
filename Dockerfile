FROM golang:1.25-alpine AS builder

# activer compilation statique
ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /app

# copier tout le projet
COPY . .

# build du binaire nommé "backend"
RUN go build -o backend .

# image finale ultra-légère
FROM scratch

WORKDIR /

# copier uniquement le binaire
COPY --from=builder /app/backend .

EXPOSE 8080

CMD ["./backend"]
