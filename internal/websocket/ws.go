package websocket

import (
	"context"
	_"fmt"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// Определение структуры сообщения
type Message struct{ 
	Username string `json:"username"`
	Content string 	`json:"content"`
	Time string 	`json:"time"`
}
// Определяем структуру клиента 
type Client struct { 
	conn *websocket.Conn
	send chan Message
	stopPing chan bool
}

var (
	upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool{
		return true
	},}
	clients    = make(map[*Client]bool)
	clientsMux = &sync.Mutex{}
	broadcast  = make(chan Message, 100) // Общий канал 

	redisClient *redis.Client
	maxMessages int64 = 250
)
// Инициальзация Redis через переменные окружения 
func init(){
	// address := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	password := os.Getenv("REDIS_PASSWORD")
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: password,
		DB: 0,
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil{
		log.Printf("Невозможно подключиться к Redis!%v", err)
	}else{
		log.Println("Победа! Подключение к Redis!")
	}
}
// Сохранение сообщений в Redis
func SaveMessages(msg Message){
	// Сначала преобразуем сообщение в JSON
	msgJSON, err := json.Marshal(msg) 
	if err != nil{
		log.Printf("Невозможно переписать в JSON сообщение: %v", err)
	}
	// Затем закидываем сериализованнное сообщение в Redis
	if err := redisClient.LPush(context.Background(), os.Getenv("REDIS_KEY"), msgJSON).Err(); err != nil{
		log.Printf("Невозможно записать данные в Redis :%v", err)
	}
	// Делим сообщение 
	if err := redisClient.LTrim(context.Background(), os.Getenv("REDIS_KEY"), 0, maxMessages-1).Err(); err != nil{
		log.Printf("Невозможно разделить сообщения: %v", err)
	}
}

func SendHistory(ws *websocket.Conn){
	messages, err := redisClient.LRange(context.Background(), os.Getenv("REDIS_KEY"), 0, maxMessages-1).Result()
	if err != nil{
		log.Printf("%v", err)
	}
	if len(messages) == 0 {
		log.Println("История сообщений пуста")
	}
	for i := len(messages)-1; i >= 0; i--{
		var msg Message
		if err := json.Unmarshal([]byte(messages[i]), &msg); err != nil{
			log.Printf("%v", err)
		}
		if err := ws.WriteJSON(msg); err != nil{
			log.Printf("%v", err)
		}
	}
}

// Функция, сохраняющая историю чата
func GetHistory(c echo.Context){
	messages, err := redisClient.LRange(context.Background(), os.Getenv("REDIS_KEY"), 0, maxMessages-1).Result()
	if err != nil{
		log.Printf("%v", err)
	}  
	var res []Message
	for i := len(messages)-1; i >=0; i--{
		var msg Message
		if err := json.Unmarshal([]byte(messages[i]), &msg); err != nil{
			log.Printf("Ошибка при десериаизации данных: %v", err)
			continue
		}
		res = append(res, msg)
	} 

	c.Response().Header().Set("Content-Type", "application/json")
	json.NewEncoder(c.Response()).Encode(res)
}


func HandleConnections(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return nil
	}
	defer ws.Close()

	ws.SetReadDeadline(time.Now().Add(180 * time.Second)) 
    ws.SetPongHandler(func(string) error { 
        ws.SetReadDeadline(time.Now().Add(180 * time.Second)) 
        return nil 
    })	
	client := &Client{ // Создаем клиента
		conn: ws,
		send: make(chan Message, 100),
		stopPing: make(chan bool),
	}
	go PingClient(client)
	
	RegisterClient(client) // Регистрируем клиента
	defer UnregisterClient(client) // Удаляем клиента при отключении
	SendHistory(ws)

// Читаем сообщения от клиента
	for { 
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Read error: %v", err)
			}
			break
		}
		msg.Time = time.Now().Format("15:04") // Добавляем время к сообщению
		log.Printf("Получено сообщение: %s", msg.Content)

		SaveMessages(msg)
		broadcast <- msg // Отправляем в канал широковещания
	}
	return nil
}

func PingClient(client *Client){
	ticker := time.NewTicker(30 * time.Second) 
    defer ticker.Stop()
// бесконечный цикл
    for { 
        select {
        case <-ticker.C:
            client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                log.Printf("Ping error: %v", err)
                return
            }
        case <-client.stopPing: // Остановить если канал закрыт
            return
        }
    }
}
// Регистрация клиента
func RegisterClient(client *Client) { 
	clientsMux.Lock()
	defer clientsMux.Unlock()
	clients[client] = true
	newClientMessage := Message{ // Уведомление о приходе нового клиента
		Username: "System",
		Content: "Пользователь зашел в чат",
		Time: 	time.Now().Format("15:04"),
	}	
	broadcast <- newClientMessage // Отправка сообщения в общий канал
	log.Printf("Новый клиент подключен. Всего клиентов: %d", len(clients))
	go client.WritePump() // Запуск горутины для отправки сообщений клиенту
}
// Удаление клиента
func UnregisterClient(client *Client) { 
	clientsMux.Lock()
	defer clientsMux.Unlock()

	close(client.stopPing)

	delete(clients, client)
	log.Printf("Клиент отключен. Осталось клиентов: %d", len(clients))

	leaveMsg := Message{ // Отправляем уведомление о выходе клиента
		Username: "System",
		Content:  "Пользователь вышел из чата",
		Time:     time.Now().Format("15:04"),
	}
	broadcast <- leaveMsg // Отправка сообщения в общий канал
	close(client.send)
}
// Отправка сообщений клиенту
func (c *Client) WritePump() { 
	for {
		msg, ok := <-c.send
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})  // Канал закрыт
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err := c.conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Write error: %v", err)
				return
			}
	}
}
// Обработчик широковещательных сообщений
func HandleMessages() {  
	for {
		msg := <- broadcast
		clientsMux.Lock()
		for client := range clients {
			select {
				case client.send <- msg:  // Сообщение отправлено в канал клиента
				default:
					log.Printf("Канал клиента переполнен, сообщение пропущено")
			}
		}
		clientsMux.Unlock()
	}
}