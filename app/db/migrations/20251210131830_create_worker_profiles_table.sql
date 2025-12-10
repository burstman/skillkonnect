-- +goose Up
CREATE TABLE IF NOT EXISTS worker_profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    user_id INTEGER NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    profession VARCHAR(255) NOT NULL,
    rating REAL DEFAULT 0.0,
    distance REAL DEFAULT 0.0,
    reviews INTEGER DEFAULT 0,
    price REAL DEFAULT 0.0,
    available BOOLEAN DEFAULT true,
    latitude REAL DEFAULT 0.0,
    longitude REAL DEFAULT 0.0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_worker_profiles_user_id ON worker_profiles(user_id);
CREATE INDEX idx_worker_profiles_deleted_at ON worker_profiles(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS worker_profiles;
