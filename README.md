# backend-go

Go / Gin / PostgreSQL で作成した REST API です。

既存の Rails ポートフォリオ（Next.js + Rails API + Terraform + ECS）で実装していた  
**teacher 認証導線の一部を Go で再実装したバックエンド API**です。

# このリポジトリの目的

このリポジトリでは以下を目的に実装しています。

- Go による Web API 実装
- JWT を用いた認証フロー
- PostgreSQL 接続
- Repository 層による責務分離

既存ポートフォリオでは以下の構成で動作させています。

```
Browser
↓
Next.js
↓
ALB
↓
Go API
↓
PostgreSQL
```

# 技術スタック

- Go
- Gin
- PostgreSQL
- JWT
- Docker Compose

# Architecture

この API は以下のような構成での利用を想定しています。

```
Browser
↓
Next.js (Frontend)
↓
ALB
↓
Go API
↓
PostgreSQL
```

既存ポートフォリオでは以下のインフラ構成で動作しています。

- Terraform
- ECS Fargate
- ALB
- RDS (PostgreSQL)
- Secrets Manager
- GitHub Actions

# セットアップ

### DB起動

```
docker compose up -d
```

### APIサーバ起動

```
go run ./cmd/server
```

# 動作確認

health check

```
curl localhost:3001/healthz
```

レスポンス

```
{"status":"ok"}
```

# 実装している主なエンドポイント

- GET /healthz
- POST /teachers/register
- POST /teachers/login
- GET /teachers/me

JWT 認証は Authorization ヘッダで行います。

```
Authorization: Bearer <JWT>
```

# 今後の改善候補

- service 層の追加
- classroom / student API 追加
- Docker で Go API 起動
- CI 追加
- テスト追加
