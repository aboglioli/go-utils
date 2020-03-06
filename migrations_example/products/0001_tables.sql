CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    sku integer NOT NULL,
    name varchar(64) NOT NULL,
    description varchar(32) NOT NULL,
)
