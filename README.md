# Application Security: SQL Injection and Buffer Overflow Vulnerabilities---

## I. Goal and Tasks of the Project

### Project Goal

Demonstrate three critical application security vulnerabilities through practical proof-of-concept implementations. The project shows how classic SQL injection, blind SQL injection, and stack-based buffer overflow can be exploited to compromise systems and extract sensitive data. Each vulnerability is deployed in an isolated Docker environment with working exploits and detailed documentation of attack vectors and defense mechanisms.

### Project Scope

The project covers three vulnerability categories from the Common Weakness Enumeration list. First, classic SQL injection demonstrates direct authentication bypass and data access through UNION-based attacks. Second, blind SQL injection shows indirect data extraction when application responses do not reveal query results. Third, buffer overflow demonstrates memory corruption leading to arbitrary code execution.

All vulnerabilities are containerized using Docker for reproducible testing environments. Automated and manual exploitation methods are provided for each vulnerability type. The implementation includes database setup, web application code, exploit scripts, and comprehensive documentation of both attack and defense strategies.

## Quick Start

This project demonstrates various application security vulnerabilities. Each component can be run using Docker Compose.

### Buffer Overflow

```bash
    docker-compose up buffer-overflow
```

See `buffer-overflow/README.md` for detailed instructions.

### SQL Injection

```bash
    docker-compose up sql-injection
```

See `sql-injection/README.md` for detailed instructions.

### Blind SQL Injection

```bash
    docker-compose up blind-sql-injection
```

See `blind-sql-injection/README.md` for detailed instructions.
