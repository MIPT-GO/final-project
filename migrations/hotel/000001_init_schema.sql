-- +goose up
-- +goose statementbegin
CREATE TABLE IF NOT EXISTS "Hotels" (
    name TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS "Rooms" (
    hotel_name TEXT NOT NULL,
    number INT NOT NULL,
    cost NUMERIC(10, 2) NOT NULL, 

    PRIMARY KEY (hotel_name, number),

    FOREIGN KEY (hotel_name) REFERENCES "Hotels"(name) ON DELETE CASCADE
);
-- +goose statementend

-- +goose down
-- +goose statementbegin
DROP TABLE IF EXISTS "Rooms";
DROP TABLE IF EXISTS "Hotels";
-- +goose statementend