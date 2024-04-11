CREATE TABLE IF NOT EXISTS public."category"
(
    id        BIGINT                                              GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name      TEXT UNIQUE                                         NOT NULL CHECK (name <> '')
        CONSTRAINT max_len_name CHECK (LENGTH(name) <= 256),
    translation   TEXT UNIQUE                                     NOT NULL CHECK (translation <> '')
        CONSTRAINT max_len_translation CHECK (LENGTH(translation) <= 256)
);