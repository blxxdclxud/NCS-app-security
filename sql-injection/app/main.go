package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blxxdclxud/NCS-app-security/sql-injection/handlers"
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

	// Проверка подключения
	if err = db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	log.Println("Connected to database")

	// Роутинг
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", handlers.VulnerableLoginHandler(db))
	http.HandleFunc("/search", handlers.VulnerableSearchHandler(db))

	log.Println("Server running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head><title>Classic SQL Injection Demo</title></head>
	<body>
		<h1>Classic SQL Injection Demo</h1>
		<ul>
			<li><a href="/login">Login Form (Vulnerable)</a></li>
			<li><a href="/search">Product Search (Vulnerable to UNION)</a></li>
		</ul>
	</body>
	</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
