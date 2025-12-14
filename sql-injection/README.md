# Classic SQL Injection Demo

This package contains a small vulnerable web application that demonstrates classic SQL injection in a realistic but controlled environment. It shows how insecure string concatenation in SQL queries leads to authentication bypass and data leakage.

## Features

- Go web server with deliberately vulnerable endpoints
- PostgreSQL database with demo users and products
- Classic SQL injection:
    - Authentication bypass on `/login`
    - UNION-based data extraction on `/search`
- Manual exploits (curl scripts)
- Automated exploitation with sqlmap

## Architecture

- Language: Go (`net/http`, `database/sql`, `lib/pq`)
- Database: PostgreSQL (DB `vuln`, schema `public`)
- Deployment: Docker + docker-compose
- Connection via environment variables:
    - `POSTGRES_USER=admin`
    - `POSTGRES_PASSWORD=pass`
    - `POSTGRES_DB=vuln`


## Vulnerabilities

### 1. Login: Authentication Bypass

Endpoint: `POST /login`  
Parameters: `username`, `password`

Issue: SQL query is built with `fmt.Sprintf` and untrusted input:

```go
query := fmt.Sprintf(
"SELECT id, username, email, is_admin FROM users WHERE username='%s' AND password='%s'",
username, password,
)
```


Example payload:

- Username: `admin' OR '1'='1'--`
- Password: anything

Effect: The condition becomes always true and the application logs in the attacker as `admin`.

### 2. Search: UNION-Based Injection

Endpoint: `GET /search?q=...`  
Parameter: `q` (search string)

Issue: `q` is concatenated into a `LIKE` clause without escaping:

```go
query := fmt.Sprintf(
"SELECT id, name, price, description FROM products WHERE name LIKE '%%%s%%'",
searchQuery,
)
```

Example payload:

```sql
' UNION SELECT id, username, password, email FROM users--
```


Effect: The response table shows users and their passwords instead of normal product data.

## Database Schema

Database: `vuln`, schema: `public`.

Tables:

- `users(id, username, password, email, is_admin)`
    - `admin / supersecret123 / admin@example.com / true`
    - `user1 / password123 / user1@example.com / false`
    - `bob / bob123 / bob@example.com / false`
- `products(id, name, price, description)`

Initialized from `database/init.sql` when the DB container starts.

## Running the Demo

From project root:


```shell
docker-compose up sql-injection db
```


Open:

- Classic SQLi app: `http://localhost:8081`

## Manual Exploits

From `sql-injection/exploits`:

## Hardening Notes

This package is intentionally insecure. To fix the issues in real code:

- Use prepared statements / parameterized queries:

```go
db.QueryRow(
"SELECT id, username, email, is_admin FROM users WHERE username=$1 AND password=$2",
username, password,
)
```

- Validate and sanitize input where possible
- Use least-privilege database users
- Add logging and monitoring for suspicious patterns
- Consider WAF rules against common SQLi payloads