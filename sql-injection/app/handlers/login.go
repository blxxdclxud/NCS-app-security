package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

func VulnerableLoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl := `
					<!DOCTYPE html>
					<html>
					<head><title>Login - Vulnerable</title></head>
					<body>
						<h1>Login Form (Vulnerable to SQLi)</h1>
						<form method="POST">
							<label>Username:</label><br>
							<input type="text" name="username" placeholder="admin"><br><br>
							<label>Password:</label><br>
							<input type="password" name="password" placeholder="password"><br><br>
							<button type="submit">Login</button>
						</form>
						<hr>
						<p><b>Try authentication bypass:</b></p>
						<code>Username: admin' OR '1'='1'--</code><br>
						<code>Password: anything</code>
					</body>
					</html>`
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(tmpl))
			return
		}

		// POST request
		username := r.FormValue("username")
		password := r.FormValue("password")

		// VULN: strings concatenation
		query := fmt.Sprintf(
			"SELECT id, username, email, is_admin FROM users WHERE username='%s' AND password='%s'",
			username, password,
		)

		fmt.Println("Executing query:", query)

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), 500)
			return
		}
		defer rows.Close()

		if rows.Next() {
			var id int
			var uname, email string
			var isAdmin bool
			rows.Scan(&id, &uname, &email, &isAdmin)

			response := fmt.Sprintf(`
					<!DOCTYPE html>
					<html>
					<head><title>Success</title></head>
					<body>
						<h1>✅ Login Successful!</h1>
						<p><b>User ID:</b> %d</p>
						<p><b>Username:</b> %s</p>
						<p><b>Email:</b> %s</p>
						<p><b>Admin:</b> %v</p>
						<a href="/classic/login">Back to login</a>
					</body>
					</html>`, id, uname, email, isAdmin)

			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(response))
		} else {
			http.Error(w, "❌ Invalid credentials", 401)
		}
	}
}
