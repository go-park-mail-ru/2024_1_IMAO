CREATE TABLE IF NOT EXISTS public.blacklist
(
    user_id_blocker     BIGINT                                               NOT NULL REFERENCES public."user" (id),
    user_id_blocked     BIGINT                                               NOT NULL REFERENCES public."user" (id),
    CONSTRAINT blacklist_uniq_together_user_id_user_id unique (user_id_blocker, user_id_blocked)
);