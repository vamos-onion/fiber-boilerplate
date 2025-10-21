-- must run with superuser
-- CREATE EXTENSION pgcrypto;

-- GLOBAL FUNCTIONS
CREATE
OR REPLACE FUNCTION fn_set_modified_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.modified_at
= now();
RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TYPE enum_gender AS ENUM ('M', 'F');

CREATE TABLE appuser
(
    id          bigserial   NOT NULL PRIMARY KEY,
    uuid        uuid        NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    created_at  timestamptz NOT NULL        DEFAULT now(),
    modified_at timestamptz NOT NULL        DEFAULT now(),
--
    name        varchar(64) NOT NULL,
    birthday    date        NOT NULL,
    gender      enum_gender NOT NULL,
    withdraw    boolean     NOT NULL        DEFAULT false
);
CREATE TRIGGER tr_appuser_update_modified_at
    BEFORE UPDATE
    ON appuser
    FOR EACH ROW
    EXECUTE PROCEDURE fn_set_modified_at();

CREATE TABLE array_test
(
    id                  bigserial   NOT NULL PRIMARY KEY,
    uuid                uuid        NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    created_at          timestamptz NOT NULL        DEFAULT now(),
    modified_at         timestamptz NOT NULL        DEFAULT now(),
--
    varchar_array_field varchar(64)[],
    text_array_field    text[],
    int_array_field     int[],
    float_array_field   float[],
    bool_array_field    boolean[]
);
CREATE TRIGGER tr_array_test_update_modified_at
    BEFORE UPDATE
    ON array_test
    FOR EACH ROW
    EXECUTE PROCEDURE fn_set_modified_at();