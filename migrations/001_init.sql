CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username TEXT UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       balance INT DEFAULT 1000 CHECK (balance >= 0)
);

CREATE TABLE merch (
                       id SERIAL PRIMARY KEY,
                       name TEXT UNIQUE NOT NULL,
                       price INT NOT NULL CHECK (price >= 0)
);

CREATE TABLE inventory (
                           id SERIAL PRIMARY KEY,
                           user_id INT REFERENCES users(id) ON DELETE CASCADE,
                           merch_id INT REFERENCES merch(id) ON DELETE CASCADE,
                           quantity INT DEFAULT 1 CHECK (quantity >= 0)
);

CREATE TABLE transactions (
                              id SERIAL PRIMARY KEY,
                              from_user INT REFERENCES users(id) ON DELETE SET NULL,
                              to_user INT REFERENCES users(id) ON DELETE SET NULL,
                              amount INT NOT NULL CHECK (amount > 0),
                              timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO merch (name, price) VALUES
                                    ('t-shirt', 80),
                                    ('cup', 20),
                                    ('book', 50),
                                    ('pen', 10),
                                    ('powerbank', 200),
                                    ('hoody', 300),
                                    ('umbrella', 200),
                                    ('socks', 10),
                                    ('wallet', 50),
                                    ('pink-hoody', 500);
