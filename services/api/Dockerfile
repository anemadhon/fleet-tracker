FROM golang:1.24-alpine AS builder

WORKDIR /app

# 1. copy go.mod + go.sum dari root
COPY go.mod go.sum ./
RUN go mod download

# 2. copy seluruh source project (termasuk tj/config, tj/pkg, services/...)
COPY . .

# 3. build main api
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./services/api/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/api .

CMD ["./api"]
