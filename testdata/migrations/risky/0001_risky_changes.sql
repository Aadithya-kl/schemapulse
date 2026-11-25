CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL
);

CREATE INDEX idx_users_email ON users (email);

ALTER TABLE users ADD COLUMN legacy_code TEXT NOT NULL;

ALTER TABLE users ALTER COLUMN email TYPE VARCHAR(320);

ALTER TABLE users ADD CONSTRAINT fk_users_org FOREIGN KEY (org_id) REFERENCES organizations (id);

ALTER TABLE users DROP COLUMN legacy_code;

DROP TABLE old_sessions;