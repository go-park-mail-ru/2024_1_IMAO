CREATE INDEX adverts_title_gin ON advert USING gin(to_tsvector('russian', title));
CREATE INDEX adverts_title_gin_english ON advert USING gin(to_tsvector('english', title));

SELECT * FROM advert WHERE to_tsvector('russian', title) @@ to_tsquery('russian', 'смартфон');

SELECT id, user_id, city_id, category_id, title, description, price, created_time, closed_time, is_used, advert_status, views
	FROM advert
	WHERE (to_tsvector(title) @@ to_tsquery(replace('за' || ':*', ' ', ' | '))) AND advert_status = 'Активно'
	ORDER BY ts_rank(to_tsvector(title), to_tsquery(replace('за' || ':*', ' ', ' | '))) DESC	 
	LIMIT 10;

SELECT ts_headline(a.title, to_tsquery(replace('ша' || ':*', ' ', ' | ')), 'StartSel=<b>,StopSel=</b>,MaxFragments=2,' || 'FragmentDelimiter=...,MaxWords=2,MinWords=1')
FROM public.advert a
WHERE (to_tsvector(a.title) @@ to_tsquery(replace('ша' || ':*', ' ', ' | '))) AND a.advert_status = 'Активно'
LIMIT 15;    

SELECT DISTINCT(LOWER(regexp_replace(ts_headline(a.title, to_tsquery(replace('ша' || ':*', ' ', ' | ')), 
								  'MaxFragments=0,' || 'FragmentDelimiter=...,MaxWords=5,MinWords=1'), '<b>|</b>', '', 'g')))
FROM public.advert a
WHERE (to_tsvector(a.title) @@ to_tsquery(replace('ша' || ':*', ' ', ' | '))) AND a.advert_status = 'Активно'
LIMIT 15; 

WITH lower_titles AS (
    SELECT DISTINCT LOWER(regexp_replace(ts_headline(a.title, to_tsquery(replace('ша' || ':*', ' ', ' | ')), 
                                  'MaxFragments=0,' || 'FragmentDelimiter=...,MaxWords=5,MinWords=1'), '<b>|</b>', '', 'g')) AS title
    FROM public.advert a
    WHERE (to_tsvector(a.title) @@ to_tsquery(replace('ша' || ':*', ' ', ' | '))) AND a.advert_status = 'Активно'
),
original_titles AS (
    SELECT DISTINCT a.title AS title
    FROM public.advert a
    WHERE (to_tsvector(a.title) @@ to_tsquery(replace('ша' || ':*', ' ', ' | '))) AND a.advert_status = 'Активно'
)
SELECT * FROM lower_titles
UNION ALL
SELECT * FROM original_titles;

WITH one_word_titles AS (
    SELECT DISTINCT LOWER(regexp_replace(ts_headline(a.title, to_tsquery(replace('сма' || ':*', ' ', ' | ')), 
                                  'MaxFragments=0,' || 'FragmentDelimiter=...,MaxWords=2,MinWords=1'), '<b>|</b>', '', 'g')) AS title
    FROM public.advert a
    WHERE (to_tsvector(a.title) @@ to_tsquery(replace('сма' || ':*', ' ', ' | '))) AND a.advert_status = 'Активно'
),
two_word_titles AS (
    SELECT DISTINCT LOWER(regexp_replace(ts_headline(a.title, to_tsquery(replace('сма' || ':*', ' ', ' | ')), 
                                  'MaxFragments=1,' || 'FragmentDelimiter=...,MaxWords=2,MinWords=1'), '<b>|</b>', '', 'g')) AS title
    FROM public.advert a
    WHERE (to_tsvector(a.title) @@ to_tsquery(replace('сма' || ':*', ' ', ' | '))) AND a.advert_status = 'Активно'
)
SELECT * FROM one_word_titles
UNION 
SELECT * FROM two_word_titles
LIMIT 8;

CREATE INDEX trgm_title_idx ON advert USING gist (title gist_trgm_ops);

CREATE EXTENSION pg_trgm;