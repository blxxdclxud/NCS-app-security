-- Users for Classic SQLi
CREATE TABLE IF NOT EXISTS users (
     id SERIAL PRIMARY KEY,
     username VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    email VARCHAR(100),
    is_admin BOOLEAN DEFAULT FALSE
    );

INSERT INTO users (username, password, email, is_admin) VALUES
    ('admin', 'supersecret123', 'admin@example.com', TRUE),
    ('user1', 'password123', 'user1@example.com', FALSE),
    ('bob', 'bob123', 'bob@example.com', FALSE);

-- Products for UNION injection
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200),
    price VARCHAR(50),
    description TEXT
    );

INSERT INTO products (name, price, description) VALUES
    ('Laptop', '$999', 'High-performance laptop'),
    ('Mouse', '$25', 'Wireless mouse'),
    ('Keyboard', '$75', 'Mechanical keyboard');
