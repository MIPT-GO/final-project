-- +goose up
-- +goose statementbegin
CREATE TABLE IF NOT EXISTS hotels (
    name TEXT PRIMARY KEY,
    email TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS rooms (
    hotel_name TEXT NOT NULL,
    number INT NOT NULL,
    cost NUMERIC(10, 2) NOT NULL, 

    PRIMARY KEY (hotel_name, number),

    FOREIGN KEY (hotel_name) REFERENCES hotels(name) ON DELETE CASCADE
);
-- +goose statementend

-- +goose down
-- +goose statementbegin
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS hotels;
-- +goose statementend