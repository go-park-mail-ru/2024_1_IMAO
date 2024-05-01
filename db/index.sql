CREATE INDEX adverts_title_gin ON advert USING gin(to_tsvector('russian', title));
CREATE INDEX adverts_title_gin_english ON advert USING gin(to_tsvector('english', title));

SELECT * FROM advert WHERE to_tsvector('russian', title) @@ to_tsquery('russian', 'смартфон');

SELECT id, user_id, city_id, category_id, title, description, price, created_time, closed_time, is_used, advert_status, views
	FROM advert
	WHERE (to_tsvector(title) @@ to_tsquery(replace('за' || ':*', ' ', ' | '))) AND advert_status = 'Активно'
	ORDER BY ts_rank(to_tsvector(title), to_tsquery(replace('за' || ':*', ' ', ' | '))) DESC	 
	LIMIT 10;