-- =============================== v0 ============================

WITH promoted_adverts AS (
    SELECT id, title, description, price, is_promoted
    FROM public.advert
    WHERE is_promoted = TRUE
	ORDER BY id
	OFFSET 25
	LIMIT 5
), non_promoted_adverts AS (
    SELECT id, title, description, price, is_promoted
    FROM public.advert
    WHERE is_promoted = FALSE
	ORDER BY id
	OFFSET 75
)
SELECT * FROM promoted_adverts
UNION ALL
SELECT * FROM non_promoted_adverts
ORDER BY is_promoted DESC, id ASC
LIMIT 20;

WITH promoted_adverts AS (
    SELECT id, title, description, price, is_promoted
    FROM public.advert
    WHERE is_promoted = TRUE
	ORDER BY id
	OFFSET (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) * (1 / (div((SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE), 5)) * 0)
	LIMIT 5
), non_promoted_adverts AS (
    SELECT id, title, description, price, is_promoted
    FROM public.advert
    WHERE is_promoted = FALSE
	ORDER BY id
	OFFSET 75
)
SELECT * FROM promoted_adverts
UNION ALL
SELECT * FROM non_promoted_adverts
ORDER BY is_promoted DESC, id ASC
LIMIT 20;


=========================================

WITH promoted_adverts AS (
    SELECT id, title, description, price, is_promoted
    FROM public.advert
    WHERE is_promoted = TRUE
	ORDER BY id
	OFFSET (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) * (1 / (div((SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE), 5)) * 6)
	LIMIT 5
), non_promoted_adverts AS (
    SELECT id, title, description, price, is_promoted
    FROM public.advert
    WHERE is_promoted = FALSE
	ORDER BY id
	OFFSET 15 * 6 + 5 * div((SELECT 
    CASE 
        WHEN (6 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) < 0 THEN 0
        ELSE (6 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE)
    END), 5)
)
SELECT * FROM promoted_adverts
UNION ALL
SELECT * FROM non_promoted_adverts
ORDER BY is_promoted DESC, id ASC
LIMIT 20;


-- =============================== v1 ============================

WITH promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = TRUE AND a.advert_status = 'Активно' AND c.translation = $2
		ORDER BY id
		OFFSET 5 * $1
		LIMIT 5
	), non_promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = FALSE AND a.advert_status = 'Активно' AND c.translation = $2
		ORDER BY id
		OFFSET 15 * $1 + 5 * div((SELECT 
		CASE 
			WHEN ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) < 0 THEN 0
			ELSE ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE)
		END), 5) + (SELECT
			CASE
				WHEN (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) - $1 * 5 < 0 
						THEN 5 - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) % 5
				ELSE 0
			END)
	)
	SELECT * FROM promoted_adverts
	UNION ALL
	SELECT * FROM non_promoted_adverts
	ORDER BY is_promoted DESC, id ASC
	LIMIT 20;


