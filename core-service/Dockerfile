FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git gcc musl-dev protobuf-dev protobuf ca-certificates && \
    go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN /go/bin/swag init -g cmd/app/main.go


RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

RUN protoc --proto_path=. \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/*.proto

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/app/main.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
# Используем /app везде для единообразия
WORKDIR /app

# Копируем из builder в текущую папку (/app)
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

EXPOSE 8080 50051
# Запускаем из /app
CMD ["./main"]