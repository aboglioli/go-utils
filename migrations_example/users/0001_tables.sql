CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    name varchar(64) NOT NULL,
    username varchar(32) NOT NULL,
    password varchar(256) NOT NULL
)
