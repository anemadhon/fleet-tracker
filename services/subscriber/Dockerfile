FROM golang:1.24-alpine AS builder

WORKDIR /app

# 1. copy go.mod + go.sum dari root
COPY go.mod go.sum ./
RUN go mod download

# 2. copy seluruh source project (termasuk tj/config, tj/pkg, services/...)
COPY . .

# RUN go test ./... -v

# 3. build main subscriber
RUN CGO_ENABLED=0 GOOS=linux go build -o subscriber ./services/subscriber/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/subscriber .

CMD ["./subscriber"]
