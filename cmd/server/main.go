
package main

import (
	"dnd_mj_app/internal/api"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func main() {
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Pages principales
	http.HandleFunc("/mj", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/mj.html")
	})
	http.HandleFunc("/player", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/player.html")
	})
	http.HandleFunc("/music", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/music.html")
	})

	// Page d'accueil
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/mj", http.StatusFound)
	})

	// API
	http.HandleFunc("/api/image", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.GetCurrentImage(w, r)
		} else if r.Method == http.MethodPost {
			api.SetCurrentImage(w, r)
		} else {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/upload-image", api.UploadImage)
	http.HandleFunc("/api/images", api.GetImageList)

	// WebSocket
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erreur WebSocket Upgrade:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Client WS déconnecté: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Erreur envoi WS: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
