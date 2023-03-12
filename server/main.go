package main

import (
	"database/sql"
	"encoding/json"
	"html"
	"log"
	"net/http"
	"strings"
	"time"

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
			http.Error(w, "Invalid JSON in request", http.StatusBadRequest)
			return
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
			log.Printf("Failed to marshal tabs object: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		if _, err = server.db.Exec(
			`INSERT OR REPLACE INTO tabs (client_id, tabs, last_updated) values (?, ?, ?)`,
			tabs.ClientId,
			j,
			time.Now(),
		); err != nil {
			log.Printf("Failed to insert into tabs: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodGet {
		if r.URL.Query().Has("deleted") {
			// deleted sessions
			rows, err := server.db.Query(`SELECT id, client_id FROM tabs_deleted`)
			if err != nil {
				log.Printf("Failed to query tabs_deleted table: %v", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
			defer rows.Close()
			var out []map[string]any
			for rows.Next() {
				var id int
				var client_id string
				rows.Scan(&id, &client_id)
				out = append(out, map[string]any{
					"id":        id,
					"client_id": client_id,
				})
			}

			j, err := json.Marshal(out)

			w.Write([]byte(j))
		} else {
			// active sessions
			rows, err := server.db.Query(`SELECT client_id, tabs FROM tabs`)
			if err != nil {
				log.Printf("Failed to query tabs table: %v", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var allTabs []Tabs = []Tabs{}
			for rows.Next() {
				var clientId string
				var tabsStr string
				if err := rows.Scan(&clientId, &tabsStr); err != nil {
					log.Printf("failed to scan rows: %v", err)
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}

				var tabsArr [][3]string
				err := json.Unmarshal([]byte(tabsStr), &tabsArr)
				if err != nil {
					// This means we stored bad data, which shouldn't happen.
					log.Printf("BADNESS: failed to unmarshal rows stored in DB: %v", err)
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}
				allTabs = append(allTabs, Tabs{clientId, tabsArr})
			}

			j, err := json.Marshal(allTabs)

			w.Write([]byte(j))
		}
	} else if r.Method == http.MethodDelete {
		id := strings.TrimPrefix(r.URL.EscapedPath(), "/tabs/")
		log.Printf("attempting to delete '%v'", id)
		if r.URL.Query().Has("deleted") {
			// delete from tabs_deleted table
			if _, err := server.db.Exec("DELETE from tabs_deleted WHERE id = ?", id); err != nil {
				log.Printf("Error when attempting to delete session: %v", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
		} else {
			// move to the tabs_deleted table
			tx, err := server.db.Begin()
			if err != nil {
				log.Printf("Failed to create tx: %v", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
			var clientId string
			if err := tx.QueryRow("SELECT client_id FROM tabs WHERE client_id = ?", id).Scan(&clientId); err != nil {
				tx.Rollback()
				if err == sql.ErrNoRows {
					http.Error(w, "Session not found", http.StatusNotFound)
				} else {
					log.Printf("Error when attempting to find session: %v", err)
					http.Error(w, "Internal error", http.StatusInternalServerError)
				}
				return
			}
			if _, err := tx.Exec("INSERT INTO tabs_deleted (client_id, tabs) SELECT client_id, tabs FROM tabs WHERE client_id = ?", id); err != nil {
				log.Printf("Query failed: %v", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				tx.Rollback()
				return
			}
			if _, err := tx.Exec("DELETE FROM tabs WHERE client_id = ?", id); err != nil {
				log.Printf("Query failed: %v", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				tx.Rollback()
				return
			}
			tx.Commit()
		}
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
                tabs JSON,
                last_updated TEXT
            );`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS tabs_deleted (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                client_id TEXT,
                tabs JSON
            );`); err != nil {
		log.Fatal(err)
	}

	server := TabsServer{db}
	mux := http.NewServeMux()
	mux.HandleFunc("/tabs/", server.tabsHandler)
	log.Print("Starting TabMan server")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
