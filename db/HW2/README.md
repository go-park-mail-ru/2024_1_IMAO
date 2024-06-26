# Защита от SQL Injections

В нашем проекте мы используем библиотеку jackc/pgx/v5, которая поддерживает параметризированные запросы с использованием плейсхолдеров `$1`, `$2` и т.д.

Пример кода:
```
SQLUpdateProfileAvatarURL := `
	INSERT INTO advert_image (url, advert_id, url_resized)
	VALUES 
    ($1, $2, $3)
	RETURNING url;`
```
запись работы [sqlmap](https://sqlmap.org/) с помощью [asciinema](https://asciinema.org/) находится в файле **demo.cast**

# Работа с БД через сервисную учетную запись

Скрипт для создания сервисной учётной записи находится в файле **service_account.sql**

В конфигурации приложения параметры подключения к СУБД исправлены на работу через сервисную учетную запись

```
const DATABASE_URL string = "postgres://vol4ok_service_account:vol4ok-Password123!@postgres:5432/IMAO_VOL4OK_2024"
```

> Скрипт создания пользователя и прав необходимо сохранить в репозитории, в отдельном скрипте (придумайте куда было бы удонее и логичнее его положить и аргументируйте свой выбор).

Мы решили положить этот файл в ту же папку, где лежит всё остальное, связанное с ДЗ2, чтобы ментору не пришлось далеко ходить

# Пулл соединений и параметры соедений

Мы реализовали connection pool на строне приложения с помощью библиотеки jackc/pgx/v5/pool 

Исходный код функции, создающий экземпляр `*pgxpool.Config` находится в [файле](https://github.com/go-park-mail-ru/2024_1_IMAO/blob/main/internal/pkg/server/repository/pgxpool.go), создание экземпляра `*pgxpool.Pool` находится в [файле](https://github.com/go-park-mail-ru/2024_1_IMAO/blob/main/internal/pkg/server/app.go) на 45 строке

> Вопрос в том, как правильно сбалансировать количество содениений в connection pool и значение max_connections в postgresql.conf

Значение max_connections в postgresql.conf должно быть чуть больше, чем максимальное количество содениений в connection pool. Таким образом, всегда будет несколько доступных соединений для прямого подключения для обслуживания и мониторинга системы

[ссылка](https://wiki.postgresql.org/wiki/Number_Of_Database_Connections)

> Также необходимо правильно настроить параметр listen_addresses. Подумайте, настройте и аргументируйте ментору, какие адреса в вашем случае здесь должны быть указаны

В listen_addresses по-хорошему должны быть указаны адреса всех микросервисов бекенда. Но у нас не получилось указать в качестве этого параметра ничего кроме `'*'` и `'127.0.0.1/24'`, так что мы оставили `'*'`, что соответствует прослушиванию всех доступных интерфейсов IPv4 и IPv6.

Исходя из того, как мы понимаем параметр listen_addresses, в нем прописывается не "кому разрешено подключаться", а на "на каких интерфейсах PostgreSQL должен принимать подключения". 

А настройки того, "кому разрешено подключаться" мы прописали в файле **pg_hba.conf**


# Настройка параметров сервера и клиента

## Таймауты

Мы поставили таймаунт в 30 секунд, потому что исходя из бизнес логики ни один запрос на бекенд не должен выполняться более 30 секунд.

```
statement_timeout = 30s 
lock_timeout = 30s
```
> Например, подумайте о граничных случаях - например, есть ли смысл ставить ограничение, в 1 минуту?

Смысла ставить ограничение в 1 минуту нет, потому что, как уже было написано выше, ни один запрос на бекенд не должен выполняться более 30 секунд.

> Будет ли у такого ограничения какая-то польза с точки зрения UX? 

Нет, не будет

> А с точки зрения защиты вашей системы от DOS атак?

В случае DDoS-атаки, злоумышленники могут попытаться выполнить очень долгие или бесконечные запросы, чтобы замедлить работу системы или вывести ее из строя. Установка statement_timeout и lock_timeout может помочь предотвратить выполнение таких запросов.



## Логгирование и протоколирование медленных запросов

Вот эти параметры из `postgresql.conf` отвечают за логирование (точнее сказать, только эти параметры мы редактировали)

```
log_line_prefix = '%t [%p]: '
logging_collector = on 
log_directory = 'log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_min_duration_statement = 5
log_error_verbosity = verbose
```

> А в какое именно (время, после которого запрос считается медленным) - вам необходимо придумать исходя из бизнес требований вашего приложения.

Мы поставили это значение равным 5ms, потому что у нас в принципе нет запросов, которые бы выполнялись дольше 6ms 


## Параметры потребления ресурсов

> Необходимо изучить конфигурируемые параметры и выбрать 6 любых параметров, которые на ваш взгляд важны для вашего приложения. Изменить их значение от значения по-умолчанию.
> Выбор параметров и значений аргументировать.

+ shared_buffers = 256MB: Этот параметр определяет количество памяти, выделенной для кэширования данных и индексов. За счёт увеличения этого значения (относительно значения по умолчанию) мы расчитываем улучшить производительность за счет уменьшения времени доступа к данным.
+ work_mem = 16MB: Этот параметр определяет максимальное количество памяти, которое PostgreSQL может использовать для выполнения операций сортировки, хеширования и других операций, требующих временного хранения данных. За счёт увеличения этого значения (относительно значения по умолчанию) мы расчитываем ускорить выполнение запросов, особенно тех, которые включают в себя большие объёмы сортировки или группировки. 
+ maintenance_work_mem = 64MB: Этот параметр определяет максимальное количество памяти, которое PostgreSQL может использовать для выполнения операций обслуживания базы данных, таких как VACUUM, CREATE INDEX и ANALYZE. За счёт увеличения этого значения (относительно значения по умолчанию) мы расчитываем ускорить эти операции, особенно на больших таблицах.
+ temp_buffers = 16MB: Этот параметр определяет количество памяти, выделенной для кэширования временных файлов. За счёт увеличения этого значения (относительно значения по умолчанию) мы расчитываем улучшить производительность за счет уменьшения времени доступа к временным файлам.
+ min_dynamic_shared_memory = 16MB: Этот параметр определяет минимальное количество памяти, которое PostgreSQL будет пытаться динамически выделить для общего кэша. За счёт увеличения этого значения (относительно значения по умолчанию) мы расчитываем обеспечить некоторую гибкость в управлении памятью, позволяя PostgreSQL адаптироваться к изменениям в нагрузке.
+ max_stack_depth = 4MB: Этот параметр определяет максимальную глубину стека для функций PostgreSQL. За счёт увеличения этого значения (относительно значения по умолчанию) мы расчитываем ускорить обработку сложных запросов с большим количеством вложенных вызовов функций. 
