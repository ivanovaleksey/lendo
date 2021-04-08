CREATE EXTENSION pgcrypto;

CREATE TABLE jobs (
    id          UUID NOT NULL DEFAULT gen_random_uuid(),
    application JSON NOT NULL,
    status      TEXT NOT NULL,

    created_at     TIMESTAMP NOT NULL DEFAULT now(),
    updated_at     TIMESTAMP NOT NULL DEFAULT now(),

    PRIMARY KEY (id)
);

CREATE INDEX jobs_status_idx ON jobs USING hash (status);
CREATE INDEX jobs_created_at_idx ON jobs USING btree (status);
