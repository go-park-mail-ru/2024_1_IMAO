CREATE TABLE IF NOT EXISTS public."cart"
(
    user_id   BIGINT                       NOT NULL REFERENCES public."user" (id),
	advert_id BIGINT                       NOT NULL REFERENCES public."advert" (id) ON DELETE CASCADE,
    CONSTRAINT cart_uniq_together_advert_id_user_id unique (user_id, advert_id)
);