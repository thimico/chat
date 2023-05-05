package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Define the structure for a ChatRoom
type ChatRoom struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

// Define the structure for a Message
type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

var chatRooms = make(map[string]*ChatRoom)
var lock = sync.RWMutex{}

// Upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	http.HandleFunc("/", handleConnections)

	log.Println("Server started on :8080")
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	username := r.URL.Query().Get("username")

	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error upgrading to WebSocket: %v", err)
	}

	defer ws.Close()

	lock.Lock()
	// Create a new chat room if it does not exist
	if _, ok := chatRooms[roomID]; !ok {
		chatRooms[roomID] = &ChatRoom{
			clients:   make(map[*websocket.Conn]bool),
			broadcast: make(chan Message),
		}
		go handleBroadcasts(roomID)
	}

	// Add the user to the chat room
	chatRooms[roomID].clients[ws] = true
	lock.Unlock()

	// Handle user messages
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(chatRooms[roomID].clients, ws)
			break
		}
		msg.Username = username
		handleCommand(roomID, &msg)
		chatRooms[roomID].broadcast <- msg
	}
}

func handleBroadcasts(roomID string) {
	for {
		msg := <-chatRooms[roomID].broadcast
		for client := range chatRooms[roomID].clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(chatRooms[roomID].clients, client)
			}
		}
	}
}

func handleCommand(roomID string, msg *Message) {
	command := strings.Split(msg.Text, "=")
	if len(command) == 2 && command[0] == "/stock" {
		go fetchStockDataAndSend(roomID, command[1])
	}
}

// fetchStockData retrieves stock data from the API and returns a formatted stock quote message.
func fetchStockData(stockCode string) (string, error) {
	// Build the API URL
	apiURL := fmt.Sprintf(os.Getenv("STOCK_API_BASE_URL")+os.Getenv("STOCK_API_PARAMS"), stockCode)

	// Send the API request
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Error sending API request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the API response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading API response body: %v", err)
		return "", err
	}

	// Parse the CSV data
	reader := csv.NewReader(strings.NewReader(string(body)))
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error parsing CSV data: %v", err)
		return "", err
	}

	// Check if the CSV data contains the expected number of rows and columns
	if len(records) < 2 || len(records[1]) < 5 {
		log.Printf("Invalid CSV data format")
		return "", fmt.Errorf("Invalid CSV data format")
	}

	// Extract the stock price from the CSV data
	stockPriceStr := records[1][4]

	// Check if the stock price is available
	if stockPriceStr == "N/D" {
		return "", fmt.Errorf("Stock price data is not available for the requested stock")
	}

	// Convert the stock price to a float
	stockPrice, err := strconv.ParseFloat(stockPriceStr, 64)
	if err != nil {
		log.Printf("Error parsing stock price: %v", err)
		return "", err
	}

	// Format and return the stock quote message
	stockQuote := fmt.Sprintf("%s quote is $%.2f per share", stockCode, stockPrice)
	return stockQuote, nil
}

func fetchStockDataAndSend(roomID, stockCode string) {
	_, err := fetchStockData(stockCode)
	if err != nil {
		log.Printf("Error parsing CSV data: %v", err)
		return
	}

}
