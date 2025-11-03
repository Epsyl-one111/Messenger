package mail

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"

	"github.com/go-gomail/gomail"
	"github.com/joho/godotenv"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
) 

func SendWithGomail(c echo.Context) error {
	if c.Request().Method != http.MethodPost{
		return c.Redirect(http.StatusFound, "/checkingcode")
	}
	if err := godotenv.Load(); err != nil{
		fmt.Print("Can't find .env file")
	}
	WorkMail := os.Getenv("WORK")
    ToMail := os.Getenv("TO")
    m := gomail.NewMessage()
	
    m.SetHeader("From", WorkMail)
	getMail := c.FormValue("Mail")
	m.SetHeader("To", getMail)
    m.SetAddressHeader("From", WorkMail, "Genesis Info Sender")
    m.SetHeader("Subject", "Регистрация в Месседжере Genesis")    
    
	txtRand := rand.Text()

    m.SetBody("text/html", "Никому не говорите Ваш код: " + txtRand)
	
    d := gomail.NewDialer("smtp.gmail.com", 587, WorkMail, "ehau xgwr kvgm smem")
    if err := d.DialAndSend(m); err != nil {
        return err
    }

	session := c.Get("session").(*sessions.Session) // Используются сессии для хранения данных 
	session.Values["ver_code"] = txtRand // Назначаем значение переменной сессии (в данном случае рандомная строка)
	session.Values["email"] = ToMail

	if err := session.Save(c.Request(), c.Response()); err != nil{
		return err
	}

	return c.Render(http.StatusOK, "sendingcode", map[string]interface{}{
		"Title": "Registration",
		"Error": "",
	})
}

func CheckCode(c echo.Context) error{
	getCode := c.FormValue("GmailCode") 
	session := c.Get("session").(*sessions.Session)
	storeCode, ok := session.Values["ver_code"].(string)
	if !ok{
		return c.Render(http.StatusOK, "sendingcode", map[string]interface{}{
			"Title": "Registration",
			"Error": "Session is over!",
		})
	}
	if storeCode == getCode{
		delete(session.Values, "ver_code") // После успешной верификации кода, сессия удаляется, чтобы не занимать память
		session.Save(c.Request(), c.Response()) // Сессия сохраняется 
		return c.Render(http.StatusOK, "registration", map[string]interface{}{
			"Title": "Registration",
			"Error": "",
		})
	}
	return c.Render(http.StatusOK, "authorization", map[string]interface{}{
		"Title": "Registration",
		"Error": "Wrong Code!",
	})
}