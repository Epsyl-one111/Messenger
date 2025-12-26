package database
import (
	"os"
	"time"
	"log"
	"net/http"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)
var (
	connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", 
		os.Getenv("POST_USER"),
		os.Getenv("POST_PASSWORD"),
		os.Getenv("POST_HOST"),
		os.Getenv("POST_PORT"),
		os.Getenv("POST_DB"),
	)
	person Person
)

type Person struct{
	Username string
	Password string
	Email string
}

// Проверка на наличие базы данных, если ее нет, он ее создает
func InitDB(){
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil{
		log.Fatalf("%v",err)
	}
	_, err = conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS data_user (
			ID SERIAL PRIMARY KEY,
            username VARCHAR(50),
            password VARCHAR(50)
        )
    `)
	if err != nil{
		time.Sleep(2 * time.Second)
		InitDB() // Рекурсия на проверку 
		return
	}
}
func RegPage(c echo.Context) error { 	
	if c.Request().Method != http.MethodPost{
		return c.Redirect(http.StatusFound, "/reg")
	}
	getUsernameReg := c.FormValue("usernameReg")
	getPasswordReg := c.FormValue("passwordReg")
	// Проверка инфы с базы даннных 
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil{
		log.Printf("Error: %v",err)
		return c.Render(http.StatusOK, "authorization", map[string]interface{}{
			"Title": "Authorization",
        	"Error": "Database connection error",
		})
	}
	defer conn.Close(context.Background())
	rows, err := conn.Query(context.Background(), "SELECT username, password FROM data_user")
	if err != nil{log.Fatal(err)}
	defer rows.Close()

	for rows.Next(){
		if err := rows.Scan(&person.Username, &person.Password); err != nil{log.Fatal(err)}
		if getUsernameReg == person.Username || getPasswordReg == person.Password{
			data := struct{Error string}{Error: "Password or login already exists"}
			return c.Render(http.StatusOK, "registration", data)
		}
	}// проверка инфы с таблиц базы данных
	WriteSQL(getUsernameReg, getPasswordReg)
	return c.Render(http.StatusOK, "registration", nil)
}
func AuthPage(c echo.Context) error{ 
	if c.Request().Method != http.MethodPost {
        return c.Redirect(http.StatusFound, "/auth")
    }
	getUsernameAuth := c.FormValue("username")
	getPasswordAuth := c.FormValue("password")
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil{
		log.Printf("Error: %v",err)
		return c.Render(http.StatusOK, "authorization", map[string]interface{}{
			"Title": "Authorization",
        	"Error": "Database connection error",
		})
	}
	defer conn.Close(context.Background())
	rows, err := conn.Query(context.Background(), "SELECT username, password FROM data_user")
	if err != nil{
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next(){
		if err := rows.Scan(&person.Username, &person.Password); err != nil{return c.Render(http.StatusOK, "authorization", map[string]interface{}{
				"Title": "Authorization",
        		"Error": "Wrong password or login",
			})
		}
		if getUsernameAuth == person.Username && getPasswordAuth == person.Password{
			return c.Redirect(http.StatusFound, "/home")
		}
	}
	return c.Render(http.StatusOK, "authorization", map[string]interface{}{
		"Title": "Authorization",
		"Error": "Wrong password or login",
	})
}
// Запись информации о клиенте в базу данных
// через енвишку для сохранности данных
func WriteSQL(username, password string) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil{
		log.Fatal(err)
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "INSERT INTO data_user (username, password) VALUES ($1, $2)", username, password) 
	if err != nil{
		log.Fatal(err)
	}
}