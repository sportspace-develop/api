# Sport Space

## develop server
1. git clone https://github.com/sportspace-develop/api.git
2. cd api/deploy
3. создать сеть для работы с фронтом ```docker network create sportspace-network```
4. запустить ```docker-compose up``` или пересобрать и запустить ```docker-compose up --build```

### file .env example
```text
TLS_ENABLE=1
TLS_CERT=./cert/server.crt
TLS_KEY=./cert/server.key
TLS_HOSTS=sportspace.ru;www.sportpace.ru
TLS_DIR_CACHE=./certs

HTTP_ADDRESS=localhost:8080
BASE_URL=http://localhost:8443
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

* TLS_ENABLE - включение шифрование, 0 - выкл, 1 - вкл, 2 - вкл. авто сер.
* TLS_CERT - путь до сертификата, обязательно при TLS_ENABLE=1
* TLS_KEY - путь до ключа, обязательно при TLS_ENABLE=1
* TLS_HOSTS - хосты, обязательно при TLS_ENABLE=2
* TLS_DIR_CACHE - путь до хранения сертификата, обязательно при TLS_ENABLE=2
* HTTP_ADDRESS - адрес запуска сервиса
* BASE_URL - базовый адрес сервера
* SECRET_KEY - ключ шифрования JWT



### Swagger docs
http://localhost:8080/swagger/index.html

### Email UI (dev)
http://localhost:1080/

# build
```bash
go build -ldflags "-X main.buildVersion=v0.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git show --oneline -s)'" ./cmd/sportspace/sportspace.go
```