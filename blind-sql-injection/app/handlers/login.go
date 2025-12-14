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
			w.Write([]byte(tmpl))
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Intentionally vulnerable (blind SQLi)
		query := fmt.Sprintf(
			"SELECT id, username, email, is_admin FROM users WHERE username='%s' AND password='%s'",
			username, password,
		)

		row := db.QueryRow(query)

		var (
			id       int
			uname    string
			email    string
			isAdmin  bool
		)

		err := row.Scan(&id, &uname, &email, &isAdmin)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Same format as image
		html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head><title>Success</title></head>
		<body>
			<h1>Login Successful!</h1>

			<p><b>User ID:</b> %d</p>
			<p><b>Username:</b> %s</p>
			<p><b>Email:</b> %s</p>
			<p><b>Admin:</b> %v</p>

			<a href="/blind/login">Back to login</a>
		</body>
		</html>`,
			id, uname, email, isAdmin,
		)

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}
}