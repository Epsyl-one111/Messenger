package handlers

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"Messanger/internal/database"
	"Messanger/internal/mail"
	"Messanger/internal/websocket"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct{
	templates *template.Template // Структура для шаблонов
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error{
	return t.templates.ExecuteTemplate(w, name, data) // Метод для рендера шаблонов 
}

func MiddlewareSessions(next echo.HandlerFunc) echo.HandlerFunc{ // Создаем Middleware для сессии
	return func(c echo.Context) error{
		var store = sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))
		// store.Options = &sessions.Options{
		// 	HttpOnly: true,

		// }
		session, _ := store.Get(c.Request(), "genesis-auth")
		c.Set("session", session)
		return next(c)
	}	
}

func HandleRequests(){
	e := echo.New()
	
	e.Use(MiddlewareSessions)
	e.Use(middleware.Logger()) // Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete}, // CORS для разрешения с браузера 
	}))
	
	e.Static("/web/static", "web/static") // Создаем для статических файлов CSS

	templates, err := template.ParseFiles( // Обработка HTML-файлов (ну тупо страниц)
		"web/templates/footer.html",
	    "web/templates/header.html",
		"web/templates/chat.html",
		"web/templates/contacts.html",
		"web/templates/sendingcode.html",
		"web/templates/sidebar.html",
	    "web/templates/authorization.html",
	    "web/templates/home.html",
		"web/templates/about.html",
		"web/templates/entermail.html",
		"web/templates/registration.html",
	); 
	if err != nil {log.Fatalf("Ошибка загрузки шаблонов:%v", err)}
	
	e.Renderer = &Template{templates: templates} // Рендер шаблонов 

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/auth")
	})

	e.GET("/home", homePage)
	e.GET("/contacts", contactsPage)
	e.GET("/about", aboutPage)
	e.GET("/chat", chatPage)

	e.GET("/auth", showAuthPage)
	e.POST("/auth/post", database.AuthPage)

	e.GET("/reg", showRegPage)
	e.POST("/reg/post", database.RegPage)

	e.GET("/entermail", showEnterMail)
	e.POST("/entermail/post", mail.SendWithGomail)

	e.GET("/checkingcode", showCheckCode)
	e.POST("/checkingcode/post", mail.CheckCode)
	
	e.GET("/ws", func(c echo.Context) error {
		websocket.HandleConnections(c.Response(), c.Request())
		return nil
	})
	go websocket.HandleMessages()
	
	e.Logger.Fatal(e.Start("0.0.0.0:8080")) // Хост для показа всем интерфейсам
}

// Домашняя страница
func homePage(c echo.Context) error{ 
	return c.Render(http.StatusOK, "home", map[string]interface{}{
		"Title": "Home page",
	})
}
// Страница о самом Месседжере
func aboutPage(c echo.Context) error{ 
	return c.Render(http.StatusOK, "about", map[string]interface{}{
		"Title": "About",
	})
}
func contactsPage(c echo.Context) error{
	return c.Render(http.StatusOK, "contacts", map[string]interface{}{
		"Title": "Contacts",
	})
}
// Страница чата
func chatPage(c echo.Context) error{ 
	return c.Render(http.StatusOK, "chat", map[string]interface{}{
		"Title": "Chat",
	})
}
// Функция, показывающая страницу регистрации
func showRegPage(c echo.Context) error{ 
	return c.Render(http.StatusOK, "registration", map[string]interface{}{
        "Title": "Registration",
        "Error": "", 
    })
}
func showEnterMail(c echo.Context) error{
	return c.Render(http.StatusOK, "entermail", map[string]interface{}{
        "Title": "Registration",
        "Error": "", 
    })
}
func showCheckCode(c echo.Context) error{
	return c.Render(http.StatusOK, "sendingcode", map[string]interface{}{
        "Title": "Registration",
        "Error": "Wrong mail!", 
    })
}
// Функция, показывающая страницу авторизации
func showAuthPage(c echo.Context) error { 
	 return c.Render(http.StatusOK, "authorization", map[string]interface{}{
        "Title": "Authorization",
        "Error": "", 
    })
}
