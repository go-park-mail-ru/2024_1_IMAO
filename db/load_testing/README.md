vegeta attack -duration=5s -rate=1000 -targets=http://127.0.0.1:8008
echo "GET http://127.0.0.1:8008" | vegeta attack -name=50qps -rate=50 -duration=5s > results.50qps.bin
echo "GET http://localhost:8080/api/adverts/?count=20&startId=2000" | vegeta attack -name=50qps -rate=50 -duration=3s > results_50qps.csv
echo "GET http://localhost:8080/api/adverts/?count=20&startId=2000" | vegeta attack -name=100qps -rate=100 -duration=3s > results_100qps.csv
echo "GET http://localhost:8080/api/adverts/?count=20&startId=2000" | vegeta attack -name=50qps -rate=50 -duration=3s --output results.bin; 
vegeta report results.bin

echo "GET http://www.vol-4-ok.ru:8080/api/adverts/?userId=2&deleted=1" | vegeta attack -name=50qps -rate=50 -duration=3s --output results.bin;

echo "GET http://www.vol-4-ok.ru:8080/api/adverts/?userId=2&deleted=1&count=20&startId=1" | vegeta attack -name=50qps -rate=10 -duration=15s --output results.bin;
echo "GET http://www.vol-4-ok.ru:8080/api/profile/3" | vegeta attack -name=50qps -rate=10 -duration=15s --output results.bin;

# Итерация №1

## Денормализация таблицы advert

Данные о просмотрах и избранном хранятся в таблицах view и favourite соответственно. И для того, чтобы получить информацию о том, сколько просмотров у объявления или сколько человек добавило объявление в избранное, необходимо проходиться по всей таблице и считать записи.

Для того, чтобы каждый раз не высчитывать эту информацию, было принято решение денормализовать таблицу advert и добавить ей 2 новых атрибута: favourites_number и views. В favourites_number будет храниться количество пользователей, которое добавило объявление в избранное, а во views количество просмотров на объявлении. favourites_number будет инкрементироваться (декрементироваться) при добавлении (удалении) записи в (из) таблицу(таблицы) favourite с помощью триггера. views будет инкрементироваться при добавлении записи в таблицу view с помощью триггера.

Исходные коды триггеров и функций, исполняемых триггерами, находятся в папке **db/misc**

Пример логов запроса **ДО** денормализации :

*кэш очищенный*
```
2024-05-13 09:55:41.813 UTC [93] LOG:  00000: duration: 104.874 ms  statement: SELECT 
	a.id, 
	a.user_id,
	a.city_id, 
	c.name AS city_name, 
	c.translation AS city_translation, 
	a.category_id, 
	cat.name AS category_name, 
	cat.translation AS category_translation, 
	a.title, 
	a.description, 
	a.price, 
	a.created_time, 
	a.closed_time, 
	a.is_used,
	a.advert_status,
	(SELECT COUNT(*) FROM public.view WHERE advert_id = a.id) AS view_count,
	(SELECT COUNT(*) FROM public.favourite WHERE advert_id = a.id) AS favorite_count
	FROM 
	public.advert a
	LEFT JOIN 
	public.city c ON a.city_id = c.id
	LEFT JOIN 
	public.category cat ON a.category_id = cat.id
	WHERE a.id = 1;
2024-05-13 09:55:41.813 UTC [93] LOCATION:  exec_simple_query, postgres.c:1369
```

*кэш прогретый*
```
2024-05-13 10:06:24.637 UTC [29] LOG:  00000: duration: 39.014 ms  statement: SELECT 
	a.id, 
	a.user_id,
	a.city_id, 
	c.name AS city_name, 
	c.translation AS city_translation, 
	a.category_id, 
	cat.name AS category_name, 
	cat.translation AS category_translation, 
	a.title, 
	a.description, 
	a.price, 
	a.created_time, 
	a.closed_time, 
	a.is_used,
	a.advert_status,
	(SELECT COUNT(*) FROM public.view WHERE advert_id = a.id) AS view_count,
	(SELECT COUNT(*) FROM public.favourite WHERE advert_id = a.id) AS favorite_count
	FROM 
	public.advert a
	LEFT JOIN 
	public.city c ON a.city_id = c.id
	LEFT JOIN 
	public.category cat ON a.category_id = cat.id
	WHERE a.id = 1;
2024-05-13 10:06:24.637 UTC [29] LOCATION:  exec_simple_query, postgres.c:1369
```

Пример лога запроса **ПОСЛЕ** денормализации:

*кэш прогретый*
```
2024-05-13 10:11:41.790 UTC [29] LOG:  00000: duration: 2.706 ms  statement: SELECT 
	a.id, 
	a.user_id,
	a.city_id, 
	c.name AS city_name, 
	c.translation AS city_translation, 
	a.category_id, 
	cat.name AS category_name, 
	cat.translation AS category_translation, 
	a.title, 
	a.description, 
	a.price, 
	a.created_time, 
	a.closed_time, 
	a.is_used,
	a.advert_status,
	a.advert_status,
	a.favourites_number
	FROM 
	public.advert a
	LEFT JOIN 
	public.city c ON a.city_id = c.id
	LEFT JOIN 
	public.category cat ON a.category_id = cat.id
	WHERE a.id = 1;
2024-05-13 10:11:41.790 UTC [29] LOCATION:  exec_simple_query, postgres.c:1369
```

## Денормализация таблицы profile

Данные о корзинах и избранном хранятся в таблицах cart и favourite соответственно. И для того, чтобы получить информацию о том, сколько товаров сейчас находится в корзине у пользователя или сколько товаров сейчас находится в избранном у пользователя, необходимо проходиться по всей таблице и считать записи. Для того, чтобы каждый раз не высчитывать эту информацию, было принято решение денормализовать таблицу profile и добавить ей 2 новых атрибута: cart_adverts_number и fav_adverts_number. В cart_adverts_number будет храниться количество объявлений, которое находится у пользователя в корзине, а во fav_adverts_number количество объявлений, которое находится у пользователя в избранном. cart_adverts_number будет инкрементироваться (декрементироваться) при добавлении (удалении) записи в (из) таблицу(таблицы) cart с помощью триггера. fav_adverts_number будет инкрементироваться (декрементироваться) при добавлении (удалении) записи в (из) таблицу(таблицы) favourite с помощью триггера.

Исходные коды триггеров и функций, исполняемых триггерами, находятся в папке **db/misc**