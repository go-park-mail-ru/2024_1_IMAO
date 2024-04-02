DROP TABLE IF EXISTS public."advert" CASCADE;
DROP TABLE IF EXISTS public."user" CASCADE;
DROP TABLE IF EXISTS public."profile" CASCADE;
DROP TABLE IF EXISTS public."order" CASCADE;
DROP TABLE IF EXISTS public."image" CASCADE;
DROP TABLE IF EXISTS public."category" CASCADE;
DROP TABLE IF EXISTS public."city" CASCADE;
DROP TABLE IF EXISTS public."status" CASCADE;
DROP TABLE IF EXISTS public."view" CASCADE;
DROP TABLE IF EXISTS public."favourite" CASCADE;
DROP TABLE IF EXISTS public."review" CASCADE;
DROP TABLE IF EXISTS public."cart" CASCADE;
DROP TABLE IF EXISTS public."subscription" CASCADE;

DROP SEQUENCE IF EXISTS advert_id_seq;
DROP SEQUENCE IF EXISTS user_id_seq;
DROP SEQUENCE IF EXISTS profile_id_seq;
DROP SEQUENCE IF EXISTS order_id_seq;
DROP SEQUENCE IF EXISTS image_id_seq;
DROP SEQUENCE IF EXISTS category_id_seq;
DROP SEQUENCE IF EXISTS city_id_seq;
DROP SEQUENCE IF EXISTS status_id_seq;
DROP SEQUENCE IF EXISTS view_id_seq;
DROP SEQUENCE IF EXISTS favourite_id_seq;
DROP SEQUENCE IF EXISTS review_id_seq;
DROP SEQUENCE IF EXISTS cart_id_seq;
DROP SEQUENCE IF EXISTS subscription_id_seq;

CREATE SEQUENCE IF NOT EXISTS advert_id_seq;
CREATE SEQUENCE IF NOT EXISTS user_id_seq;
CREATE SEQUENCE IF NOT EXISTS profile_id_seq;
CREATE SEQUENCE IF NOT EXISTS order_id_seq;
CREATE SEQUENCE IF NOT EXISTS image_id_seq;
CREATE SEQUENCE IF NOT EXISTS category_id_seq;
CREATE SEQUENCE IF NOT EXISTS city_id_seq;
CREATE SEQUENCE IF NOT EXISTS status_id_seq;
CREATE SEQUENCE IF NOT EXISTS view_id_seq;
CREATE SEQUENCE IF NOT EXISTS favourite_id_seq;
CREATE SEQUENCE IF NOT EXISTS review_id_seq;
CREATE SEQUENCE IF NOT EXISTS cart_id_seq;
CREATE SEQUENCE IF NOT EXISTS subscription_id_seq;


CREATE TABLE IF NOT EXISTS public."advert"
(
    id              BIGINT                   DEFAULT NEXTVAL('advert_id_seq'::regclass) NOT NULL PRIMARY KEY,
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
    status_id       BIGINT                   DEFAULT 1                                   NOT NULL REFERENCES public."status" (id)
);


CREATE TABLE IF NOT EXISTS public."user"
(
    id              BIGINT              DEFAULT NEXTVAL('user_id_seq'::regclass) NOT NULL PRIMARY KEY,
    email           TEXT UNIQUE                                                  NOT NULL CHECK (email <> '')
        CONSTRAINT max_len_email CHECK (LENGTH(email) <= 256),
    password_hash   TEXT                                                         NOT NULL CHECK (password_hash <> '')
        CONSTRAINT max_len_password_hash CHECK (LENGTH(password_hash) <= 256)
);


CREATE TABLE IF NOT EXISTS public."profile"
(
    id         BIGINT                DEFAULT NEXTVAL('profile_id_seq'::regclass) NOT NULL PRIMARY KEY,
    user_id    BIGINT                                                            NOT NULL REFERENCES public."user" (id),
    city_id    BIGINT                                                            NOT NULL REFERENCES public."city" (id),
    email      TEXT UNIQUE                                                       NOT NULL CHECK (email <> '')
        CONSTRAINT max_len_email CHECK (LENGTH(email) <= 256),
    phone      TEXT UNIQUE DEFAULT NULL
        CONSTRAINT max_len_phone CHECK (LENGTH(phone) <= 18),
    name       TEXT DEFAULT NULL
        CONSTRAINT max_len_name CHECK (LENGTH(name) <= 256),
    surname       TEXT DEFAULT NULL
        CONSTRAINT max_len_surname CHECK (LENGTH(surname) <= 256),    
    regtime TIMESTAMP WITH TIME ZONE DEFAULT NOW()                               NOT NULL,
    verified        BOOLEAN                  DEFAULT FALSE                       NOT NULL,
    avatar_url     TEXT UNIQUE
        CONSTRAINT max_len_avatar_url CHECK (LENGTH(avatar_url) <= 256),
);

CREATE TABLE IF NOT EXISTS public."order"
(
    id         BIGINT                   DEFAULT NEXTVAL('order_id_seq'::regclass) NOT NULL PRIMARY KEY,
    user_id    BIGINT                                                             NOT NULL REFERENCES public."user" (id),
    advert_id  BIGINT                                                             NOT NULL REFERENCES public."advert" (id),
    status_id  BIGINT                   DEFAULT 1                                 NOT NULL REFERENCES public."status" (id)
    created_time TIMESTAMP WITH TIME ZONE DEFAULT NOW()                           NOT NULL,
    updated_time TIMESTAMP WITH TIME ZONE DEFAULT NOW()                           NOT NULL
        CONSTRAINT updated_time_is_after_created_time CHECK (updated_time >= created_time),
    closed_time  TIMESTAMP WITH TIME ZONE DEFAULT NOW()                           NOT NULL
        CONSTRAINT closed_time_is_after_created_time CHECK (closed_time >= created_time),
    phone          TEXT                                                           NOT NULL
        CONSTRAINT max_len_phone CHECK (LENGTH(phone) <= 18),  
    name           TEXT                                                           NOT NULL
        CONSTRAINT max_len_name CHECK (LENGTH(name) <= 256),
    surname       TEXT                                                            NOT NULL
        CONSTRAINT max_len_surname CHECK (LENGTH(surname) <= 256), 
    patronymic    TEXT                                                            NOT NULL
        CONSTRAINT max_len_patronymic CHECK (LENGTH(patronymic) <= 256),  
    email         TEXT                                                            NOT NULL
        CONSTRAINT max_len_patronymic CHECK (LENGTH(email) <= 256)             
    
);

CREATE TABLE IF NOT EXISTS public."image"
(
    id         BIGINT DEFAULT NEXTVAL('image_id_seq'::regclass) NOT NULL PRIMARY KEY,
    url        TEXT                                             NOT NULL CHECK (url <> '')
        CONSTRAINT max_len_url CHECK (LENGTH(url) <= 256),
    advert_id  BIGINT                                           NOT NULL REFERENCES public."advert" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public."category"
(
    id        BIGINT DEFAULT NEXTVAL('category_id_seq'::regclass) NOT NULL PRIMARY KEY,
    name      TEXT UNIQUE                                         NOT NULL CHECK (name <> '')
        CONSTRAINT max_len_name CHECK (LENGTH(name) <= 256),
    translation   TEXT UNIQUE                                     NOT NULL CHECK (translation <> '')
        CONSTRAINT max_len_translation CHECK (LENGTH(translation) <= 256)
);

CREATE TABLE IF NOT EXISTS public."city"
(
    id              BIGINT            DEFAULT NEXTVAL('city_id_seq'::regclass) NOT NULL PRIMARY KEY,
    name            TEXT                                                       NOT NULL CHECK (name <> '')
        CONSTRAINT max_len_name CHECK (LENGTH(name) <= 256)
    translation   TEXT UNIQUE                                     NOT NULL CHECK (translation <> '')
        CONSTRAINT max_len_translation CHECK (LENGTH(translation) <= 256)    
);

CREATE TABLE IF NOT EXISTS public."status"
(
    id              BIGINT            DEFAULT NEXTVAL('status_id_seq'::regclass) NOT NULL PRIMARY KEY,
    name            TEXT                                                         NOT NULL CHECK (name <> '')
        CONSTRAINT max_len_name CHECK (LENGTH(name) <= 256)
    
);

CREATE TABLE IF NOT EXISTS public."view"
(
    id              BIGINT      DEFAULT NEXTVAL('view_id_seq'::regclass) NOT NULL PRIMARY KEY,
    user_id         BIGINT                                               NOT NULL REFERENCES public."user" (id),
    advert_id      BIGINT                                               NOT NULL REFERENCES public."advert" (id) ON DELETE CASCADE,
    CONSTRAINT uniq_together_advert_id_user_id unique (user_id, advert_id)
);

CREATE TABLE IF NOT EXISTS public."favourite"
(
    id        BIGINT DEFAULT NEXTVAL('favourite_id_seq'::regclass) NOT NULL PRIMARY KEY,
    user_id   BIGINT                                               NOT NULL REFERENCES public."user" (id),
    advert_id BIGINT                                                NOT NULL REFERENCES public."advert" (id) ON DELETE CASCADE,
    CONSTRAINT uniq_together_advert_id_user_id unique (user_id, advert_id)
);

CREATE TABLE IF NOT EXISTS public."review"
(
    id        BIGINT DEFAULT NEXTVAL('review_id_seq'::regclass)    NOT NULL PRIMARY KEY,
    user_id   BIGINT                                               NOT NULL REFERENCES public."user" (id),
    advert_id BIGINT                                               NOT NULL REFERENCES public."advert" (id) ON DELETE CASCADE,
    review    TEXT DEFAULT NULL
        CONSTRAINT max_len_review CHECK (LENGTH(review) <= 256),
    created_time    TIMESTAMP WITH TIME ZONE DEFAULT NOW()         NOT NULL, 
    rating    SMALLINT DEFAULT 1
        CONSTRAINT rating_interval CHECK (rating >= 1 and rating <= 5),   
    CONSTRAINT uniq_together_advert_id_user_id unique (user_id, advert_id)
);

CREATE TABLE IF NOT EXISTS public."cart"
(
    id        BIGINT DEFAULT NEXTVAL('cart'::regclass) NOT NULL PRIMARY KEY,
    user_id   BIGINT                       CONSTRAINT uniq_together_advert_id_user_id unique (user_id, advert_id)
ChatGPT
Ограничение uniq_together_advert_id_user_id гарантирует уникальность комбинаций значений в столбцах user_id и advert_id. Это означает, что в таблице не может быть двух или более записей с одинаковым значением user_id и advert_id, что предотвращает возможные дубликаты в связи между пользователем и объявлением.                        NOT NULL REFERENCES public."user" (id),
    advert_id BIGINT                                               NOT NULL REFERENCES public."advert" (id) ON DELETE CASCADE,
    CONSTRAINT uniq_together_advert_id_user_id unique (user_id, advert_id)
);

CREATE TABLE IF NOT EXISTS public."subscription"
(
    id        BIGINT DEFAULT NEXTVAL('subscription'::regclass) NOT NULL PRIMARY KEY,
    user_id_subscriber   BIGINT                                               NOT NULL REFERENCES public."user" (id),
    user_id_merchant     BIGINT                                               NOT NULL REFERENCES public."user" (id),
    CONSTRAINT uniq_together_user_id_user_id unique (user_id_subscriber, user_id_merchant)
);
