package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Домашняя страница
func HomePage(c echo.Context) error{ 
	return c.Render(http.StatusOK, "home", map[string]interface{}{
		"Title": "Home page",
	})
}
// Страница о самом Месседжере
func AboutPage(c echo.Context) error{ 
	return c.Render(http.StatusOK, "about", map[string]interface{}{
		"Title": "About",
	})
}
func ContactsPage(c echo.Context) error{
	return c.Render(http.StatusOK, "contacts", map[string]interface{}{
		"Title": "Contacts",
	})
}
// Страница чата
func ChatPage(c echo.Context) error{ 
	return c.Render(http.StatusOK, "chat", map[string]interface{}{
		"Title": "Chat",
	})
}
// Функция, показывающая страницу регистрации
func ShowRegPage(c echo.Context) error{ 
	return c.Render(http.StatusOK, "registration", map[string]interface{}{
        "Title": "Registration",
        "Error": "", 
    })
}
func ShowEnterMail(c echo.Context) error{
	return c.Render(http.StatusOK, "entermail", map[string]interface{}{
        "Title": "Registration",
        "Error": "", 
    })
}
func ShowCheckCode(c echo.Context) error{
	return c.Render(http.StatusOK, "sendingcode", map[string]interface{}{
        "Title": "Registration",
        "Error": "Wrong mail!", 
    })
}
// Функция, показывающая страницу авторизации
func ShowAuthPage(c echo.Context) error { 
	 return c.Render(http.StatusOK, "authorization", map[string]interface{}{
        "Title": "Authorization",
        "Error": "", 
    })
}