CREATE TYPE order_status AS ENUM ('В корзине', 'В обработке', 'Оплачено', 'Закрыт');
CREATE TYPE advert_status AS ENUM ('Скрыто', 'Удалено', 'Активно', 'Продано');
CREATE TYPE complain_type AS ENUM ('Спам', 'Мошенничество', 'Запрещённые материалы', 'Продажа запрещенных товаров', 'Другое');