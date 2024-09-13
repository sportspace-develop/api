# Sport Space

## develop server
1. git clone https://github.com/sportspace-develop/api.git
2. cd api/deploy
3. запустить ```docker-compose up``` или пересобрать и запустить ```docker-compose up --build```

### file .env example
```text
HTTP_ADDRESS=localhost:8080
SECRET_KEY=secret_key_phrase
LOG_LEVEL=debug
DATABASE_URI=postgres://root:root@localhost:5432/sportspace?sslmode=disable
GIN_MODE=release

SOURCE=dev

MAIL_SMTP_HOST=localhost
MAIL_SMTP_PORT=1025
MAIL_SENDER=no-report@test.ru
MAIL_SENDER_PASSWORD=test
MAIL_SECURE=0
```

### Swagger docs
http://localhost:8080/swagger/index.html

### Email UI (dev)
http://localhost:1080/

# build
```bash
go build -ldflags "-X main.buildVersion=v0.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git show --oneline -s)'" ./cmd/sportspace/sportspace.go
```