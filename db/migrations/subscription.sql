CREATE TABLE IF NOT EXISTS public.subscription
(
    user_id_subscriber   BIGINT                                               NOT NULL REFERENCES public."user" (id),
    user_id_merchant     BIGINT                                               NOT NULL REFERENCES public."user" (id),
    CONSTRAINT subscription_uniq_together_user_id_user_id unique (user_id_subscriber, user_id_merchant)
);