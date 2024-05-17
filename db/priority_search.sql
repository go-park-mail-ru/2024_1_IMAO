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



