CREATE TABLE IF NOT EXISTS public."advert"
(
    id              BIGINT                                                               GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id         BIGINT                                                               NOT NULL REFERENCES public."user" (id),
    city_id         BIGINT                                                               NOT NULL REFERENCES public."city" (id),
    category_id     BIGINT                                                               NOT NULL REFERENCES public."category" (id),
    title           TEXT                                                                 NOT NULL CHECK (title <> '')
        CONSTRAINT max_len_title CHECK (LENGTH(title) <= 256),
    description     TEXT                                                                 NOT NULL CHECK (description <> '')
        CONSTRAINT max_len_description CHECK (LENGTH(description) <= 2000),
    price           BIGINT                   DEFAULT 0                                   NOT NULL
        CONSTRAINT not_negative_price CHECK (price >= 0),
    created_time    TIMESTAMP WITH TIME ZONE DEFAULT NOW()                               NOT NULL,
    closed_time     TIMESTAMP WITH TIME ZONE DEFAULT NOW()                               NOT NULL,
        CONSTRAINT closed_time_is_after_created_time CHECK (closed_time >= created_time),
    is_used         BOOLEAN                  DEFAULT FALSE                               NOT NULL,
    advert_status   advert_status            DEFAULT 'Активно'                           NOT NULL 
);