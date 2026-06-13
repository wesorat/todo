-- перенести в новую sql, сделать хэш добавить индексы
CREATE TABLE refresh (
    id BIGSERIAL PRIMARY KEY,
    refresh_hash VARCHAR(255) UNIQUE NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE
)