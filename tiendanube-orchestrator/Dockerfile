FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api/main.go


FROM alpine:latest

COPY --from=builder /app/server /server

EXPOSE 8000

ENTRYPOINT ["/server"]