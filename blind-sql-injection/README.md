# Blind SQL Injection 

This topic demonstrates boolean-based and time-based blind SQL injection in a Go web application with PostgreSQL.

## Structure

- `app`: Go HTTP server with vulnerable `/handlers/login` and `/handlers/product` endpoints.
- `exploits`: Python scripts and a sqlmap wrapper to exploit blind SQLi.

## Running

1. Start database and app via `docker-compose.yml` at project root.
2. Build and run the blind-sql-injection app on port 8082.
3. Visit `http://localhost:8082` to see the demo.
4. Run exploits from the `exploits/` directory.