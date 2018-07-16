CREATE TABLE google_results (
    id bigserial primary key,
    created_at TIMESTAMPTZ,
    detection_confidence DECIMAL(4) NOT NULL,
    blurred DECIMAL(4) NOT NULL,
    joy DECIMAL(4) NOT NULL,
    sorrow DECIMAL(4) NOT NULL,
    surprise DECIMAL(4) NOT NULL,
    image bytea NOT NULL,
    user_id integer REFERENCES users(id) NOT NULL
);
