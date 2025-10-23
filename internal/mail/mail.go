package mail

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"

	"github.com/go-gomail/gomail"
	"github.com/joho/godotenv"
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
    m.SetHeader("To", ToMail)
    m.SetAddressHeader("From", WorkMail, "Genesis Info Sender")
    m.SetHeader("Subject", "Регистрация в Месседжере Genesis")    
    txtRand := rand.Text()
	
    m.SetBody("text/html", "Никому не говорите Ваш код: " + txtRand)
	
    d := gomail.NewDialer("smtp.gmail.com", 587, WorkMail, "ehau xgwr kvgm smem")
    if err := d.DialAndSend(m); err != nil {
        return err
    }
	// return c.Redirect(http.StatusFound, "/checkingcode")
	return c.Render(http.StatusOK, "checking_code", map[string]interface{}{
		"Title": "Registration",
		"Error": "",
	})
}

func CheckCode(c echo.Context, txtRand string) error{
	getCode := c.FormValue("GmailCode")
	if txtRand == getCode{
		return c.Redirect(http.StatusFound, "/reg")
	}
	return nil
}