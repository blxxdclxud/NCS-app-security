package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Anvar823/NCS-app-security/blind-sql-injection/handlers"
	_ "github.com/lib/pq"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "db"
	}

	connStr := fmt.Sprintf(
		"host=%s port=5432 user=admin password=pass dbname=vuln sslmode=disable",
		dbHost,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	log.Println("Connected to database (blind-sql-injection)")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/blind/login", handlers.BlindLoginHandler(db))
	http.HandleFunc("/blind/product", handlers.BlindProductHandler(db))

	log.Println("Server running on http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head><title>Blind SQL Injection Demo</title></head>
	<body>
		<h1>Blind SQL Injection Demo</h1>
		<ul>
			<li><a href="/blind/login">Boolean-based blind login</a></li>
			<li><a href="/blind/product?id=1">Time-based product lookup</a></li>
		</ul>
	</body>
	</html>`
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(html))
}