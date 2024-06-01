UPDATE public.advert
	SET   promotion_start= now(), promotion_duration='28 days'
	WHERE is_promoted = true;

-- ============== =========================

UPDATE public.advert
	SET  price_history = ARRAY['{"updated_time":"' || DATE_TRUNC('second', created_time)::TIMESTAMP WITHOUT TIME ZONE || '", "new_price":' || price || '}']::jsonb[]
	WHERE id >= 1; 

-- ============== =========================

UPDATE public.advert
	SET   is_promoted=false, promotion_start=null, promotion_duration=null
	WHERE promotion_start + promotion_duration < now();  



-- ============== =========================

(SELECT COALESCE(AVG(r.rating), 0) AS average_rating FROM public.review r INNER JOIN advert a ON r.advert_id = a.id WHERE a.user_id = 4)


UPDATE public.profile
	SET  rating = (SELECT COALESCE(AVG(r.rating), 0) AS average_rating FROM public.review r INNER JOIN advert a ON r.advert_id = a.id WHERE a.user_id = profile.user_id)
	WHERE profile.user_id > 0;

(SELECT COALESCE(r.rating, 0) AS rating FROM public.advert a LEFT JOIN public.review r ON r.advert_id = a.id WHERE a.id = ord.advert_id)
