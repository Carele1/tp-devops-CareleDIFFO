# ---------- Étape 1 : compilation ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /backend .

# ---------- Étape 2 : image finale ----------
FROM alpine:3.20

COPY --from=builder /backend /backend

EXPOSE 8080

CMD ["/backend"]