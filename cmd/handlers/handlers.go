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
	"Messanger/web/handlers"

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
	// Создаем для статических файлов CSS
	e.Static("/web/static", "web/static") 
	// Обработка HTML-файлов (страниц)
	templates, err := template.ParseFiles( 
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
	// Рендер шаблонов 
	e.Renderer = &Template{templates: templates} 

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/auth")
	})

	e.GET("/api/history", func(c echo.Context) error{
		websocket.GetHistory(c)
		return nil
	})

	e.GET("/home", handlers.HomePage)
	e.GET("/contacts", handlers.ContactsPage)
	e.GET("/about", handlers.AboutPage)
	e.GET("/chat", handlers.ChatPage)

	e.GET("/auth", handlers.ShowAuthPage)
	e.POST("/auth/post", database.AuthPage)

	e.GET("/reg", handlers.ShowRegPage)
	e.POST("/reg/post", database.RegPage)

	e.GET("/entermail", handlers.ShowEnterMail)
	e.POST("/entermail/post", mail.SendWithGomail)

	e.GET("/checkingcode", handlers.ShowCheckCode)
	e.POST("/checkingcode/post", mail.CheckCode)
	
	e.GET("/ws", func(c echo.Context) error {
		websocket.HandleConnections(c)
		return nil
	})

	go websocket.HandleMessages()
	
	e.Logger.Fatal(e.Start("0.0.0.0:8080")) // Хост для показа всем интерфейсам
}
