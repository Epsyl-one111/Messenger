package main

import (
	"log"
	"github.com/joho/godotenv"
	"Messanger/cmd/handlers"
	"Messanger/internal/database"
)

func main(){
	if err := godotenv.Load(); err != nil{log.Println("Can't connect to .env file!")}
	database.InitDB() // Проверка на базу данных
	handlers.HandleRequests() // Поддержка запросов
}
