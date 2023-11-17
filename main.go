package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go.etcd.io/bbolt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var db *bbolt.DB

func initDB() error {
	var err error
	db, err = bbolt.Open("requestLogs.db", 0600, nil)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Requests"))
		return err
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
	}
	body := string(bodyBytes)

	var headers []string
	for name, values := range r.Header {
		for _, value := range values {
			headers = append(headers, fmt.Sprintf("%v: %v", name, value))
		}
	}

	data := RequestData{
		Method:    r.Method,
		URL:       r.URL.String(),
		Path:      r.URL.Path,
		Host:      r.Host,
		Timestamp: time.Now().Format(time.RFC3339),
		Body:      body,
		Headers:   strings.Join(headers, ", "),
	}

	if err = SaveRequest(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	broadcast <- data.String()

	w.Write([]byte("Request received"))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		log.Printf("Received a websocket message: %v", msg)
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go handleMessages()
	http.HandleFunc("/ws", handleConnections)

	http.HandleFunc("/api/requests/delete", DeleteAllRequestsHandler)

	http.HandleFunc("/api/requests/all", GetAllRequests)

	http.HandleFunc("/requests/", requestHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Println("Listening on :8000...")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
