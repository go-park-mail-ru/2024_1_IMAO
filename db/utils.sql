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