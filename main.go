package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nyc-app-maps/api"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)


var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}


func main() {
	if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Println("WebSocket upgrade error:", err)
            return
        }
        defer conn.Close()
        watchFiles(conn)
    })

	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message": "Hello from Golang API!"}`)
	})

	http.HandleFunc("/generate", generateHandler)

	port := ":8080"
	fmt.Printf("Server running at http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}

func watchFiles(conn *websocket.Conn) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    err = watcher.Add("./static")
    if err != nil {
        log.Fatal(err)
    }

    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            if event.Op&fsnotify.Write == fsnotify.Write {
                log.Println("File changed:", event.Name)
                if err := conn.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
                    log.Println("WebSocket error:", err)
                    return
                }
            }
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            log.Println("Watcher error:", err)
        }
    }
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := api.GenerateContent()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result interface{}
	if len(resp.Candidates) > 0 {
		result = resp.Candidates[0]
	} else {
		result = []interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
    
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}