CREATE SCHEMA IF NOT EXISTS booking;

CREATE TABLE IF NOT EXISTS booking.reservations (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    hotel TEXT NOT NULL,
    number BIGINT,
    created_at TIMESTAMPTZ DEFAULT now(),
);

CREATE INDEX IF NOT EXISTS idx_reservations_email ON booking.reservations (email);
CREATE INDEX IF NOT EXISTS idx_reservations_hotel ON booking.reservations (hotel);
CREATE INDEX IF NOT EXISTS idx_reservations_number ON booking.reservations (number);
