FROM golang:1.25-alpine AS builder

WORKDIR /app


RUN apk add --no-cache git gcc musl-dev protobuf-dev


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

RUN ./scripts/generate.sh


RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/app


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

EXPOSE 8080 50051

CMD ["./main"]