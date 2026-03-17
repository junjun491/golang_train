# ==========
# build stage
# ==========
FROM golang:1.26.1 AS builder

WORKDIR /app

# 依存解決を先にやってキャッシュを効かせる
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# Linux向けに静的バイナリを作成
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# ==========
# runtime stage
# ==========
FROM alpine:3.20

WORKDIR /app

# HTTPS通信用の証明書だけ入れる
RUN apk add --no-cache ca-certificates

# builder で作ったバイナリをコピー
COPY --from=builder /app/server .

# ECS/Fargateで使う想定ポート
EXPOSE 8080

# 実行
CMD ["./server"]