CREATE EXTENSION pgcrypto;

CREATE TABLE applications (
    id         UUID NOT NULL DEFAULT gen_random_uuid(),
    first_name TEXT NOT NULL,
    last_name  TEXT NOT NULL,
    status     TEXT NOT NULL,

    created_at TIMESTAMP,
    updated_at TIMESTAMP,

    PRIMARY KEY (id)
);

CREATE INDEX applications_status_idx ON applications USING hash (status);
CREATE INDEX applications_created_at_idx ON applications USING btree (status);
