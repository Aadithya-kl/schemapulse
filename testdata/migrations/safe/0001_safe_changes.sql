CREATE TABLE organizations (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    org_id BIGINT,
    status TEXT NOT NULL DEFAULT 'active'
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users (email);

ALTER TABLE users
    ADD CONSTRAINT fk_users_org FOREIGN KEY (org_id) REFERENCES organizations (id) NOT VALID;

ALTER TABLE users VALIDATE CONSTRAINT fk_users_org;