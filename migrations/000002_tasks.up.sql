CREATE TABLE tasks (
                       id          BIGSERIAL PRIMARY KEY,
                       owner_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       title       TEXT NOT NULL,
                       description TEXT DEFAULT '',
                       status      BOOLEAN NOT NULL DEFAULT false,
                       created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
                       updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
