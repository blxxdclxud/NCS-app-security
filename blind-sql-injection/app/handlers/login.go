package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

func BlindLoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := `
				<!DOCTYPE html>
				<html>
				<head><title>Blind Login</title></head>
				<body>
					<h1>Boolean-based Blind SQLi Login</h1>
					<form method="POST">
						<label>Username:</label><br>
						<input type="text" name="username" value="admin"><br><br>
						<label>Password:</label><br>
						<input type="password" name="password"><br><br>
						<button type="submit">Login</button>
					</form>
					<hr>
					<p>This form is vulnerable to boolean-based blind SQL injection.</p>
				</body>
				</html>`
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(tmpl))
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		query := fmt.Sprintf(
			"SELECT id FROM users WHERE username='%s' AND password='%s'",
			username, password,
		)

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		if rows.Next() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	}
}
