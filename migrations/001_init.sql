-- Validations handled at handler level via go-playground/validator

CREATE TABLE IF NOT EXISTS chargers (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    location    TEXT NOT NULL,
    timezone    TEXT NOT NULL DEFAULT 'UTC',
    created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE TABLE IF NOT EXISTS pricing_periods (
    id            TEXT PRIMARY KEY,
    charger_id    TEXT NOT NULL REFERENCES chargers(id) ON DELETE CASCADE,
    start_hour    INTEGER NOT NULL,
    start_minute  INTEGER NOT NULL DEFAULT 0,
    end_hour      INTEGER NOT NULL,
    end_minute    INTEGER NOT NULL DEFAULT 0,
    price_per_kwh REAL NOT NULL,
    created_at    TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_pricing_periods_charger ON pricing_periods(charger_id);

PRAGMA foreign_keys = ON;
