package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

func VulnerableSearchHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Query().Get("q") == "" {
			tmpl := `
				<!DOCTYPE html>
				<html>
				<head><title>Product Search</title></head>
				<body>
					<h1>Product Search (Vulnerable)</h1>
					<form method="GET">
						<input type="text" name="q" placeholder="Search products..." size="50">
						<button type="submit">Search</button>
					</form>
					<hr>
					<p><b>Try UNION injection:</b></p>
					<code>' UNION SELECT id, username, password, email FROM users--</code>
				</body>
				</html>`
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(tmpl))
			return
		}

		searchQuery := r.URL.Query().Get("q")

		// VULN:
		rows, err := db.Query(fmt.Sprintf(
			"SELECT id, name, price, description FROM products WHERE name LIKE '%%%s%%'",
			searchQuery,
		))

		// SECURE:
		//searchPattern := "%" + searchQuery + "%"
		//
		//rows, err := db.Query(
		//	"SELECT id, name, price, description FROM products WHERE name LIKE $1",
		//	searchPattern,
		//)

		if err != nil {
			http.Error(w, "Database error: "+err.Error(), 500)
			return
		}
		defer rows.Close()

		result := `<!DOCTYPE html><html><head><title>Results</title></head><body><h1>Search Results</h1><table border="1"><tr><th>ID</th><th>Name</th><th>Price</th><th>Description</th></tr>`

		for rows.Next() {
			var id int
			var name, price, description string
			rows.Scan(&id, &name, &price, &description)
			result += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td></tr>", id, name, price, description)
		}

		result += `</table><br><a href="/classic/search">New search</a></body></html>`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(result))
	}
}
