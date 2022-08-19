CREATE TABLE IF NOT EXISTS users(
    user_id     uuid            DEFAULT uuid_generate_v4 (),
    login       VARCHAR(32)     NOT NULL UNIQUE, 
    password    VARCHAR         NOT NULL,
    created_at  TIMESTAMP       WITHOUT TIME ZONE DEFAULT(now() AT TIME ZONE('Europe/Moscow'))
);
