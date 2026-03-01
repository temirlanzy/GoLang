package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	var err error

	for i := 1; i <= 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Database connected!")
				break
			}
		}
		log.Println("Waiting for database...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to database")
	}

	http.HandleFunc("/users", usersHandler)

	log.Println("Starting the Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		rows, _ := db.Query("SELECT id, name FROM users")
		defer rows.Close()

		for rows.Next() {
			var id int
			var name string
			rows.Scan(&id, &name)
			fmt.Fprintf(w, "%d - %s\n", id, name)
		}

	case http.MethodPost:
		name := r.URL.Query().Get("name")
		db.Exec("INSERT INTO users(name) VALUES($1)", name)
		fmt.Fprintln(w, "User added")

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		db.Exec("DELETE FROM users WHERE id=$1", id)
		fmt.Fprintln(w, "User deleted")

	case http.MethodPut:
		id := r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")
		db.Exec("UPDATE users SET name=$1 WHERE id=$2", name, id)
		fmt.Fprintln(w, "User updated")
	}
}
