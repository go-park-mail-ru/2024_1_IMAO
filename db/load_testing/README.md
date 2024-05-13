# Заполнение базы данных через API Endpoint’а, который создает основную сущность

Заполнение происходило в 2 захода

Для заполнения использовался скрипт *db/perf_test/vegeta_formdata/main.go*

*1 заход*

```
Requests      [total, rate, throughput]  450000, 125.00, 125.00
Duration      [total, attack, wait]      1h0m0.0125987s, 59m59.9898183s, 22.7804ms
Latencies     [mean, 50, 95, 99, max]    45.159402ms, 16.977245ms, 63.762902ms, 1.004549585s, 2.2006796s
Bytes In      [total, mean]              273949774, 608.78
Bytes Out     [total, mean]              535498810, 1190.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               0:1  200:449999
Error Set:
```

*2 заход*

```
Requests      [total, rate, throughput]  375000, 125.00, 125.00
Duration      [total, attack, wait]      50m0.0028684s, 49m59.9894925s, 13.3759ms
Latencies     [mean, 50, 95, 99, max]    33.20195ms, 16.757371ms, 57.527088ms, 466.685613ms, 2.1270584s
Bytes In      [total, mean]              228292146, 608.78
Bytes Out     [total, mean]              446250000, 1190.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:375000
Error Set:
```

Пример логов **INSERT** запроса к **"основной сущности"** :

```
2024-05-13 16:27:22.750 UTC [1095] LOG:  00000: duration: 1.062 ms  execute stmtcache_f2b7dd0433adb4737ab4529798d38b8bad2e584ed07b49d2: WITH ins AS (
			INSERT INTO advert (user_id, city_id, category_id, title, description, price, is_used, phone)
			SELECT
				$1,
				city.id,
				category.id,
				$2,
				$3,
				$4,
				$5,
				$8
			FROM
				city
			JOIN
				category ON city.name = $6 AND category.translation = $7
			RETURNING 
				advert.id, 
				advert.user_id,
				advert.city_id, 
				advert.category_id, 
				advert.title, 
				advert.description,
				advert.created_time,
				advert.closed_time, 
				advert.price, 
				advert.is_used
		)
		SELECT ins.*, c.name AS city_name, c.translation AS city_translation, cat.name AS category_name, cat.translation AS category_translation
		FROM ins
		LEFT JOIN public.city c ON ins.city_id = c.id
		LEFT JOIN public.category cat ON ins.category_id = cat.id;
2024-05-13 16:27:22.750 UTC [1095] DETAIL:  parameters: $1 = '2', $2 = 'Дз 3 по базам данных', $3 = 'Данное объявление создано в рамках Дз 3 по базам данных', $4 = '777', $5 = 't', $6 = 'Москва', $7 = 'handmade', $8 = '7 777 777 77 77'
2024-05-13 16:27:22.750 UTC [1095] LOCATION:  exec_execute_message, postgres.c:2313
2024-05-13 16:27:22.813 UTC [1095] LOG:  00000: duration: 1.028 ms  statement: commit
2024-05-13 16:27:22.813 UTC [1095] LOCATION:  exec_simple_query, postgres.c:1369
```

Для моделирования реальной нагрузки также было необходимо заполнить таблицу **favourite**, а так как для этой таблицы стоит уникальность на пару значений *user_id*, *advert_id*, то с текущим количеством пользователей (а именно 10) не было возможности обеспечить необходимое количество сочетаний для заполнение таблицы **favourite** в 1 млн записей, так что было принято решение также добавить в таблицы **user** и **favourite** по 2 тыс записей. После чего можно было добвить в таблицу **favourite** 1 млн записей.

Скрипты для локального заполнения таблиц **user**, **favourite** и **favourite** находятся в файле *db/perf_test/local_fill_db/main.go*

# Итерация №0

## Проведение замеров

Пример логов **SELECT** запроса к **"основной сущности"** :

<u>*кэш очищенный*</u>
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

<u>*кэш прогретый*</u>
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

## Проведение нагрузочного текстирования с помощью vegeta


<u>*5 запросов в секунду*</u>
```
Requests      [total, rate, throughput]         50, 5.10, 5.08
Duration      [total, attack, wait]             9.834s, 9.795s, 38.725ms
Latencies     [min, mean, 50, 90, 95, 99, max]  37.516ms, 40.522ms, 38.727ms, 40.122ms, 40.884ms, 115.75ms, 115.75ms
Bytes In      [total, mean]                     26250, 525.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:50
Error Set:
```

<u>*10 запросов в секунду*</u>
```
Requests      [total, rate, throughput]         100, 10.10, 9.96
Duration      [total, attack, wait]             9.936s, 9.901s, 35.06ms
Latencies     [min, mean, 50, 90, 95, 99, max]  20.962ms, 37.474ms, 36.533ms, 38.026ms, 38.511ms, 90.946ms, 142.794ms
Bytes In      [total, mean]                     52171, 521.71
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           99.00%
Status Codes  [code:count]                      200:99  404:1
Error Set:
404 Not Found
```

<u>*15 запросов в секунду*</u>
```
Requests      [total, rate, throughput]         150, 15.12, 14.97
Duration      [total, attack, wait]             9.955s, 9.921s, 34.107ms
Latencies     [min, mean, 50, 90, 95, 99, max]  10.111ms, 34.053ms, 33.595ms, 34.445ms, 34.638ms, 35.498ms, 107.032ms
Bytes In      [total, mean]                     78421, 522.81
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           99.33%
Status Codes  [code:count]                      200:149  404:1
Error Set:
404 Not Found
```

<u>*20 запросов в секунду*</u>
```
Requests      [total, rate, throughput]         200, 20.08, 9.90
Duration      [total, attack, wait]             9.996s, 9.961s, 34.713ms
Latencies     [min, mean, 50, 90, 95, 99, max]  40.5µs, 17.356ms, 16.788ms, 33.835ms, 34.172ms, 48.152ms, 73.724ms
Bytes In      [total, mean]                     71771, 358.86
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           49.50%
Status Codes  [code:count]                      200:99  404:101
Error Set:
404 Not Found
```

<u>*25 запросов в секунду*</u>
```
Requests      [total, rate, throughput]         250, 25.10, 4.80
Duration      [total, attack, wait]             9.992s, 9.959s, 32.92ms
Latencies     [min, mean, 50, 90, 95, 99, max]  359.4µs, 7.431ms, 547.133µs, 33.829ms, 34.724ms, 37.695ms, 73.582ms
Bytes In      [total, mean]                     64792, 259.17
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           19.20%
Status Codes  [code:count]                      200:48  404:202
Error Set:
404 Not Found
```

# Итерация №1

## Денормализация таблицы advert

Данные о просмотрах и избранном хранятся в таблицах view и favourite соответственно. И для того, чтобы получить информацию о том, сколько просмотров у объявления или сколько человек добавило объявление в избранное, необходимо проходиться по всей таблице и считать записи.

Для того, чтобы каждый раз не высчитывать эту информацию, было принято решение денормализовать таблицу advert и добавить ей 2 новых атрибута: favourites_number и views. В favourites_number будет храниться количество пользователей, которое добавило объявление в избранное, а во views количество просмотров на объявлении. favourites_number будет инкрементироваться (декрементироваться) при добавлении (удалении) записи в (из) таблицу(таблицы) favourite с помощью триггера. views будет инкрементироваться при добавлении записи в таблицу view с помощью триггера.

Исходные коды триггеров и функций, исполняемых триггерами, находятся в папке **db/misc**

Пример лога запроса **ПОСЛЕ** денормализации:

<u>*кэш прогретый*</u>
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

### Нагрузочное тестирование **ПОСЛЕ** денормализации :

<u>*250 запросов в секунду*</u>
```
Requests      [total, rate, throughput]         3750, 250.08, 216.03
Duration      [total, attack, wait]             17.358s, 14.995s, 2.363s
Latencies     [min, mean, 50, 90, 95, 99, max]  367.233ms, 4.972s, 4.767s, 7.347s, 8.021s, 8.931s, 9.493s
Bytes In      [total, mean]                     1657054, 441.88
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:3750
Error Set:
```

<u>*500 запросов в секунду*</u>
```
Requests      [total, rate, throughput]         4999, 499.95, 497.67
Duration      [total, attack, wait]             10.001s, 9.999s, 1.517ms
Latencies     [min, mean, 50, 90, 95, 99, max]  167.9µs, 2.647ms, 2.099ms, 3.294ms, 3.777ms, 5.595ms, 90.215ms
Bytes In      [total, mean]                     2617237, 523.55
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           99.56%
Status Codes  [code:count]                      200:4977  404:22
Error Set:
404 Not Found
```

<u>*1000 запросов в секунду*</u>
```
Requests      [total, rate, throughput]         10000, 1001.31, 996.18
Duration      [total, attack, wait]             9.989s, 9.987s, 2.262ms
Latencies     [min, mean, 50, 90, 95, 99, max]  517µs, 2.434ms, 1.869ms, 2.441ms, 2.671ms, 12.665ms, 126.454ms
Bytes In      [total, mean]                     5233879, 523.39
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           99.51%
Status Codes  [code:count]                      200:9951  404:49
Error Set:
404 Not Found
```

## Денормализация таблицы profile

Данные о корзинах и избранном хранятся в таблицах cart и favourite соответственно. И для того, чтобы получить информацию о том, сколько товаров сейчас находится в корзине у пользователя или сколько товаров сейчас находится в избранном у пользователя, необходимо проходиться по всей таблице и считать записи. Для того, чтобы каждый раз не высчитывать эту информацию, было принято решение денормализовать таблицу profile и добавить ей 2 новых атрибута: cart_adverts_number и fav_adverts_number. В cart_adverts_number будет храниться количество объявлений, которое находится у пользователя в корзине, а во fav_adverts_number количество объявлений, которое находится у пользователя в избранном. cart_adverts_number будет инкрементироваться (декрементироваться) при добавлении (удалении) записи в (из) таблицу(таблицы) cart с помощью триггера. fav_adverts_number будет инкрементироваться (декрементироваться) при добавлении (удалении) записи в (из) таблицу(таблицы) favourite с помощью триггера.

Исходные коды триггеров и функций, исполняемых триггерами, находятся в папке **db/misc**