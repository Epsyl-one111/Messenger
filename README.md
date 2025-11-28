# Messenger

Данный проект представляет собой мессенджер, написанный на языке программирования Go. В этом проекте использовались такие технологиии как базы данных SQL и NoSQL, а именно Redis, PostgreSQL, WebSocket, работа с фреймвроком Echo, контейнеризация Docker, а также управление версиями с помощью Git. 

# Как установить

    1) Клонируйте репозиторий:
        git clone https://github.com/Epsyl-one111/Messenger.git
        
    2) Установить библиотеки: 
        go get "github.com/go-gomail/gomail"
        go get "github.com/redis/go-redis/v9"
        go get "github.com/joho/godotenv"
        go get "github.com/jackc/pgx/v4"
	    go get "github.com/labstack/echo/v4"
        go get "github.com/gorilla/sessions"
        go get "github.com/labstack/echo/v4/middleware"
        go get "github.com/gorilla/websocket"
        ну или же:
        go mod tidy 
    
# Как запустить?
    Для запуска мессенджера без Docker-а, сначала нужно зайти в файл .env и поменять знаечние POST_HOST=localhost, затем нужно перейти в main.go, который находится в папке cmd/app и запустить его командой:
        go run cmd/app/main.go
    Перейти на localhost:8080 и только теперь Вы можете пользоваться эти мессенджером

    Если у Вас установлен Docker на вашем ПК, то Вы можете запустить программу с помощью команды:
        docker-compose build
        docker-compose up -d 

# Структура проекта
    - cmd\app\main.go - 
    - cmd\handlers\handlers.go

    - internal\database\db.go
    - internal\mail\mail.go
    - internal\websocket\ws.go

    - web\handlers\pghandlers.go
    - web\static
    - web\templates

    - docker-compose.yml - 
    - Dockerfile

    - go.mod - библиотеки для работы с мессенджером
    - go.sum

    - .env