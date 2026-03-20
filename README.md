# backend-go

Go / Gin / PostgreSQL で作成した REST API です。

[既存の Rails ポートフォリオ](https://github.com/junjun491/otayori_app)（Next.js + Rails API + Terraform + ECS）で実装していた  
**teacher 認証導線の一部を Go で再実装したバックエンド API**です。

# このリポジトリの目的

このリポジトリでは以下を目的に実装しています。

- Go による Web API 実装
- JWT を用いた認証フロー
- PostgreSQL 接続
- Repository 層による責務分離

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

# セットアップ

## 環境変数

- DATABASE_URL: PostgreSQL の接続文字列
- JWT_SECRET: JWT の署名に使用するシークレットキー
- PORT: API サーバの起動ポート（未指定の場合は 3001）

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
一部実装中です

```
Authorization: Bearer <JWT>
```

# 今後の改善候補

- service 層の追加
- classroom / student API 追加
- CI 追加
- テスト追加
