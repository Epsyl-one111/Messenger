# О самом проекте
Данный проект представляет собой мессенджер с общим чатом, авторизацией и регистрацией через почту, написанный на языке программирования Go. В этом проекте использовались такие технологиии как базы данных SQL и NoSQL(Redis, PostgreSQL), WebSocket, работа с фреймвроком Echo, контейнеризация Docker, а также управление версиями с помощью Git. 

![Screenshot](https://github.com/BiFroZZy/Messenger/blob/main/web/static/photos/2025-11-29_23-07-54.png)
![Screenshot](https://github.com/BiFroZZy/Messenger/blob/main/web/static/photos/2025-11-29_23-11-31.png)

# Как установить?
    1) Клонируйте репозиторий:
        git clone https://github.com/BiFroZZy/Messenger
        
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
    Если захотите остановить, то напишите команду:
        docker-compose down 
![Screenshot](https://github.com/BiFroZZy/Messenger/blob/main/web/static/photos/2025-11-29_23-14-52.png)
# Структура проекта
    - cmd\app\main.go - главный запускаемый файл
    - cmd\handlers\handlers.go - хэндлеры рутов страниц, middleware и обьявление хоста

    - internal\database\db.go - работа с базой данных
    - internal\mail\mail.go - отправка кода на почту / работа с почтой
    - internal\websocket\ws.go - функции для поддерджки WebSocket соединения  

    - web\handlers\pghandlers.go - хэндлеры страниц проекта
    - web\static - CSS даанные (оформление страниц проекта)
    - web\templates - HTML-файлы (страницы)

    - docker-compose.yml - управление контейнерами docker
    - Dockerfile - инструкция сборки контейнера

    - go.mod - библиотеки для работы с мессенджером
    - go.sum

    - .env - секреты
