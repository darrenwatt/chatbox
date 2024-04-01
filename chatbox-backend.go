package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

// variables

// room definition
var (
	rooms      = make(map[string][]*websocket.Conn)
	roomIDLock sync.Mutex
)

// upgrade to websocket handler
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by returning true
		return true
	},
}

// main http server
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and execute the HTML template
	tmpl, err := template.ParseFiles("index.html")

	// handle errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}

/*
// websocket handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		http.Error(w, "Could not upgrade to WebSocket", http.StatusBadRequest)
		return
	}

	// Close the connection when done
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %v", err)
			return
		}
	}(conn)

	log.Println("WebSocket connection established")

	// Add the user to an available chat room
	roomID := findAvailableRoom()
	log.Printf("Found chatroom id: %v", roomID)

	// Check if the room is full
	if len(rooms[roomID]) >= 2 {
		// Room is full, inform the user or redirect to a different page
		err := conn.WriteMessage(websocket.TextMessage, []byte("Room is full. Try again later."))
		if err != nil {
			log.Printf("Error writing message: %v", err)
			return
		}
		return
	}

	rooms[roomID] = append(rooms[roomID], conn)

	// Handle incoming messages from the client
	go handleMessage(conn, roomID)

	// Keep the connection open and handle messages
	for {
		// Read message from the WebSocket connection
		_, message, err := conn.ReadMessage()
		if err != nil {
			// Handle read error (e.g., user disconnected)
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message from client %s: %s", conn.RemoteAddr(), message)

		// Broadcast message to other users in the same room
		broadcastMessage(conn, roomID, message)
	}
}
*/

// websocket handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		http.Error(w, "Could not upgrade to WebSocket", http.StatusBadRequest)
		return
	}

	// Generate a random ID for the user
	userID := strconv.Itoa(rand.Intn(1000000)) // Change this to generate a unique ID as needed

	// Close the connection when done
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			// Handle error (e.g., logging)
			log.Printf("Error closing connection: %v", err)
			return
		}
	}(conn)

	log.Printf("WebSocket connection established for user %s", userID)

	// Add the user to an available chat room
	roomID := findAvailableRoom()
	log.Printf("Found chatroom id: %v", roomID)

	// Check if the room is full
	if len(rooms[roomID]) >= 2 {
		// Room is full, inform the user or redirect to a different page
		err := conn.WriteMessage(websocket.TextMessage, []byte("Room is full. Try again later."))
		if err != nil {
			// Handle error (e.g., logging)
			log.Printf("Error writing message: %v", err)
			return
		}
		return
	}

	rooms[roomID] = append(rooms[roomID], conn)

	// Handle incoming messages and broadcast to other users in the room
	go handleMessage(conn, roomID)

	// do this all the time to keep listening and posting ...
	for {
		// Read message from the WebSocket connection
		_, message, err := conn.ReadMessage()
		if err != nil {
			// Handle read error (e.g., user disconnected)
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message from user %s: %s", userID, message)

		// Broadcast message to other users in the same room
		broadcastMessage(conn, roomID, message)
	}
}

// find chatroom
func findAvailableRoom() string {
	for id, users := range rooms {
		if len(users) < 2 {
			return id
		}
	}
	// No available rooms found, create a new one
	roomID := generateRoomID()
	rooms[roomID] = []*websocket.Conn{}
	return roomID
}

// generateRoomID generates a unique room ID.
func generateRoomID() string {
	roomIDLock.Lock()
	defer roomIDLock.Unlock()

	// this line no longer required, rand.Seed is deprecated
	// rand.Seed(time.Now().UnixNano())
	id := strconv.Itoa(rand.Intn(1000000))
	for {
		_, exists := rooms[id]
		if !exists {
			break
		}
		id = strconv.Itoa(rand.Intn(1000000))
	}
	return id
}

// handleMessage handles incoming messages from a WebSocket connection.
func handleMessage(conn *websocket.Conn, roomID string) {
	defer func() {
		// Recover from panic caused by writing to a closed connection
		if r := recover(); r != nil {
			log.Printf("Recovering from panic: %v", r)
		}
	}()

	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			// Handle read error (e.g., user disconnected)
			break
		}

		log.Printf("Received message from client %s: %s", conn.RemoteAddr(), message)

		// Broadcast message to other users in the same room
		broadcastMessage(conn, roomID, message)
	}
}

// broadcastMessage broadcasts a message to other users in the same room.
func broadcastMessage(sender *websocket.Conn, roomID string, message []byte) {
	roomIDLock.Lock()
	defer roomIDLock.Unlock()

	// Retrieve the list of connections in the room
	connections := rooms[roomID]

	// Iterate through each connection in the room
	for _, conn := range connections {
		// Skip the sender
		if conn != sender {
			// Write the message to the connection
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				// Handle error (e.g., user disconnected)
				log.Printf("Error broadcasting message to %s: %v", conn.RemoteAddr(), err)
			}
		} else {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				// Handle error (e.g., user disconnected)
				log.Printf("Error broadcasting message to %s: %v", conn.RemoteAddr(), err)
			}
		}
	}
}

func main() {
	// Define routes and corresponding handlers
	http.HandleFunc("/", indexHandler)

	// Serve static files from the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Register WebSocket handler
	http.HandleFunc("/ws", wsHandler)

	// Start the web server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
