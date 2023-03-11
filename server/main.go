package main

import (
	"database/sql"
	"encoding/json"
	"html"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type TabsServer struct {
	db *sql.DB
}

type Tabs struct {
	ClientId string      `json:"client_id"`
	Tabs     [][3]string `json:"tabs"`
}

func (server *TabsServer) tabsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var tabs Tabs
		err := json.NewDecoder(r.Body).Decode(&tabs)
		if err != nil {
			log.Fatal(err)
		}

		if tabs.Tabs == nil {
			http.Error(w, "tabs cannot be null", http.StatusBadRequest)
			return
		}

		for i, t := range tabs.Tabs {
			tabs.Tabs[i] = [3]string{html.EscapeString(t[0]), html.EscapeString(t[1]), html.EscapeString(t[2])}
		}

		j, err := json.Marshal(tabs.Tabs)
		if err != nil {
			log.Fatal(err)
		}

		if _, err = server.db.Exec(
			`INSERT OR REPLACE INTO tabs (client_id, tabs) values (?, ?)`,
			tabs.ClientId,
			j,
		); err != nil {
			log.Fatal(err)
		}
	} else if r.Method == http.MethodGet {
		rows, err := server.db.Query(`SELECT client_id, tabs FROM tabs`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var allTabs []Tabs
		for rows.Next() {
			var clientId string
			var tabsStr string
			if err := rows.Scan(&clientId, &tabsStr); err == sql.ErrNoRows {
				// FIXME: set status code
				w.Write([]byte(`{"error": "no rows"}`))
				return
			}

			var tabsArr [][3]string
			err := json.Unmarshal([]byte(tabsStr), &tabsArr)
			if err != nil {
				log.Fatal(err)
			}
			allTabs = append(allTabs, Tabs{clientId, tabsArr})
		}

		j, err := json.Marshal(allTabs)

		w.Write([]byte(j))
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
