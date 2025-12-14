package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func BlindProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing id", http.StatusBadRequest)
			return
		}

		query := fmt.Sprintf(
			"SELECT name, price FROM products WHERE id=%s",
			id,
		)

		start := time.Now()
		row := db.QueryRow(query)

		var name string
		var price string
		err := row.Scan(&name, &price)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		duration := time.Since(start)

		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(
			fmt.Sprintf("Product: %s, price: %s (query time: %v)", name, price, duration),
		))
	}
}
