## Classic vs blind SQL injection

Classic (error/union‑based) SQL injection: the attacker can directly see the effect of the injection in the HTTP response. For example, a UNION‑based payload returns extra rows in an HTML table, or a verbose database error reveals table and column names. One or a few requests are enough to dump sensitive data.

Blind SQL injection: the application does not show errors or query results. The attacker only observes whether the request is treated as success or failure (boolean‑based) or how long the response takes (time‑based). Data is extracted bit‑by‑bit or character‑by‑character through many requests, but an attacker can still fully recover values like the admin password.

In this project, classic SQLi is demonstrated in task 1 (login bypass and UNION‑based extraction of the users table), while blind SQLi in task 2 recovers the admin password from the same table using boolean‑based and time‑based techniques.

## Character‑by‑character extraction with time delays

In time‑based blind SQL injection, the attacker makes the database “sleep” when a condition is true and respond immediately when it is false. By measuring response times, a script can reconstruct secret data one character at a time.

Typical workflow:

The attacker crafts a payload that checks one character of the password, for example (PostgreSQL style):
1; SELECT CASE WHEN (SUBSTRING(password,1,1)='s') THEN pg_sleep(3) ELSE pg_sleep(0) END FROM users WHERE username='admin'--.

A script sends this request and measures the response time:

if the response is slower than a threshold (for example, > 2 seconds), the condition is treated as true, so the first character of the password is s;

if the response is fast, the condition is false, and the script tries the next character in the alphabet.

The process repeats for positions 1, 2, 3, … and for all candidate characters. Eventually, the full password string is recovered, even though the application never shows the password in any response.

Your 2_time_blind_sqli.py script implements exactly this logic against the time‑based vulnerable endpoint and reconstructs the admin user’s password.

## Defenses and monitoring

Code‑level defenses
Prepared statements / parameterized queries. All database access should use parameter binding instead of string concatenation. In Go, this means using calls like db.QueryRow("SELECT ... WHERE username=$1 AND password=$2", username, password) instead of fmt.Sprintf(...). This prevents user input from changing the structure of the SQL query.

Least‑privilege database accounts. The application should connect with a DB user that only has the minimal required permissions (no DROP, CREATE, or access to unrelated schemas). Even if SQL injection occurs, the impact is limited.

Input validation and normalization. For parameters that should be numeric or from a fixed set of values (IDs, pagination, sort order), validate and cast them explicitly and reject anything that does not match the expected format.

Infrastructure‑level defenses
Web Application Firewall (WAF). A WAF can block common SQLi patterns such as ' OR 1=1, UNION SELECT, pg_sleep, SLEEP, WAITFOR DELAY, and similar keywords, and can also enforce rate limits for suspicious endpoints.

Rate limiting and CAPTCHAs. Blind SQLi needs hundreds or thousands of nearly identical requests. Per‑IP and per‑account rate limiting, or CAPTCHAs on sensitive actions (login, password reset, search) slow down such attacks dramatically and make them easier to detect.

Logging and monitoring
Logging suspicious input patterns. Log requests where parameters contain typical SQLi markers: single quotes, SQL comments (--, /* */), and SQL keywords such as UNION, SELECT, pg_sleep, SLEEP, WAITFOR. These logs help detect SQLi attempts early, even if they fail.

Latency anomaly detection. Time‑based blind SQLi produces characteristic traffic: many similar requests to one endpoint with periodic spikes in response time. Collecting latency metrics per route and raising alerts on unusual response‑time patterns can reveal ongoing time‑based attacks.