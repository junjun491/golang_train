# backend-go

Go / Gin / PostgreSQL で作成したAPIです。  
既存のRailsポートフォリオで実装していた teacher 認証導線の一部を、Goで再実装しています。

## 実装している機能

- health check
- teachers 一覧取得
- teacher 1件取得
- teacher 登録
- teacher ログイン
- teacher 自分情報取得
- JWT 認証
- bcrypt によるパスワードハッシュ化

## 技術スタック

- Go
- Gin
- PostgreSQL
- pgx
- JWT
- bcrypt
- Docker Compose

## ディレクトリ構成

@@@
backend-go
├ cmd
│ └ server
│ └ main.go
├ internal
│ ├ auth
│ │ ├ jwt.go
│ │ └ password.go
│ ├ config
│ │ └ config.go
│ ├ db
│ │ └ db.go
│ ├ handler
│ │ ├ task_handler.go
│ │ ├ teacher_auth_handler.go
│ │ └ teacher_handler.go
│ ├ middleware
│ │ └ auth_middleware.go
│ ├ model
│ │ ├ task.go
│ │ └ teacher.go
│ └ repository
│ └ teacher_repository.go
├ docker-compose.yml
├ schema.sql
├ .env.example
└ go.mod
@@@

## セットアップ

### 1. DBを起動

@@@
docker compose up -d
@@@

### 2. 環境変数ファイルを作成

@@@
cp .env.example .env
@@@

.env の例:

@@@
JWT_SECRET=dev-secret
DATABASE_URL=postgres://appuser:password@localhost:5432/app_db
@@@

### 3. スキーマを投入

@@@
docker exec -i go-app-db psql -U appuser -d app_db < schema.sql
@@@

### 4. APIサーバを起動

@@@
go run ./cmd/server
@@@

## 動作確認

### health check

@@@
curl localhost:8080/healthz
@@@

期待レスポンス:

@@@
{"status":"ok"}
@@@

---

## API一覧

- GET /healthz
- GET /teachers
- GET /teachers/:id
- POST /teachers/register
- POST /teachers/login
- GET /teachers/me

---

## API利用例

### teachers 一覧取得

@@@
curl localhost:8080/teachers
@@@

---

### teacher 1件取得

@@@
curl localhost:8080/teachers/1
@@@

---

### teacher 登録

@@@
curl -i -X POST localhost:8080/teachers/register \
 -H "Content-Type: application/json" \
 -d '{
"teacher": {
"name": "Atsushi",
"email": "test@example.com",
"password": "password123",
"password_confirmation": "password123"
}
}'
@@@

成功時のレスポンス例:

@@@
HTTP/1.1 201 Created
Authorization: Bearer <JWT>
@@@

@@@
{
"data": {
"id": 1,
"name": "Atsushi",
"email": "test@example.com"
}
}
@@@

---

### teacher ログイン

@@@
curl -i -X POST localhost:8080/teachers/login \
 -H "Content-Type: application/json" \
 -d '{
"teacher": {
"email": "test@example.com",
"password": "password123"
}
}'
@@@

成功時のレスポンス例:

@@@
HTTP/1.1 200 OK
Authorization: Bearer <JWT>
@@@

@@@
{
"data": {
"id": 1,
"name": "Atsushi",
"email": "test@example.com"
}
}
@@@

---

### teacher 自分情報取得

ログインまたは登録で取得した JWT を Authorization ヘッダに付与して呼び出します。

@@@
curl -i localhost:8080/teachers/me \
 -H "Authorization: Bearer <JWT>"
@@@

成功時のレスポンス例:

@@@
{
"data": {
"id": 1,
"name": "Atsushi",
"email": "test@example.com"
}
}
@@@

---

## このAPIの位置づけ

この backend-go は、既存のRailsポートフォリオで実装していた teacher 認証導線をベースに、以下を Go で再実装したものです。

- teacher 登録
- teacher ログイン
- 認証済み teacher 情報取得
- JWT 認証 middleware

単純な Todo API ではなく、既存アプリのドメインと認証導線を別言語で再実装することを目的にしています。

## 今後の改善候補

- repository 層の拡張
- service 層の追加
- classroom / student / message 系APIの追加
- Docker で Go API も起動できる構成への拡張
- CI 追加
- 環境ごとの設定整理
