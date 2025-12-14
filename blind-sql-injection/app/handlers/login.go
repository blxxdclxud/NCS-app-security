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
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)

			html := `
			<!DOCTYPE html>
			<html>
			<head><title>Login Successful</title></head>
			<body>
				<h1>Login Successful!</h1>
				<p><b>Username:</b> admin</p>
				<p><b>Note:</b> This page is returned even though the application
				uses a blind (boolean-based) SQL injection point. The attacker
				normally discovers the correct password by observing only whether
				the response is success or failure.</p>
				<a href="/blind/login">Back to login</a>
			</body>
			</html>`

			_, _ = w.Write([]byte(html))
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	}
}
