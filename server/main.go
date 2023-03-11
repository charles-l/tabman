package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type TabsServer struct {
	db *sql.DB
}

func (server *TabsServer) tabsHandler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	if r.Method == http.MethodPost {
		var tabs [][2]string
		err := json.NewDecoder(r.Body).Decode(&tabs)
		if err != nil {
			log.Fatal(err)
		}

		j, err := json.Marshal(tabs)
		if err != nil {
			log.Fatal(err)
		}

		if _, err = server.db.Exec(
			`INSERT OR REPLACE INTO tabs (client_id, tabs) values (?, ?)`,
			hostname,
			j,
		); err != nil {
			log.Fatal(err)
		}
	} else if r.Method == http.MethodGet {
		row := server.db.QueryRow(`SELECT tabs FROM tabs WHERE client_id = ?`, hostname)
		var tabsStr string
		if err = row.Scan(&tabsStr); err == sql.ErrNoRows {
			w.Write([]byte("no rows"))
			return
		}

		w.Write([]byte(tabsStr))
	}
}

func byeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bye..."))
}

func main() {
	const file string = "tabs.db"
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS tabs (
            client_id TEXT PRIMARY KEY,
            tabs JSON);
            `); err != nil {
		log.Fatal(err)
	}

	server := TabsServer{db}
	mux := http.NewServeMux()
	mux.HandleFunc("/tabs/", server.tabsHandler)
	mux.HandleFunc("/b/", byeHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
