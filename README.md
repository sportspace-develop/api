# Sport Space

## develop server
1. git clone https://github.com/sportspace-develop/api.git
2. cd api/deploy
3. docker-compose up

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

### Docker
```cmd
docker-compose up --build
```

### API Docs
http://localhost:8080/swagger/index.html

### Email UI (dev)
http://localhost:1080/