CREATE OR REPLACE FUNCTION increment_views()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE advert
    SET views = views + 1
    WHERE id = NEW.advert_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION increment_cart_adverts_number()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE profile
    SET cart_adverts_number = cart_adverts_number + 1
    WHERE user_id = NEW.user_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION decrement_cart_adverts_number()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE profile
    SET cart_adverts_number = cart_adverts_number - 1
    WHERE user_id = OLD.user_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION increment_fav_adverts_number()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE profile
    SET fav_adverts_number = fav_adverts_number + 1
    WHERE user_id = NEW.user_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION decrement_fav_adverts_number()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE profile
    SET fav_adverts_number = fav_adverts_number - 1
    WHERE user_id = OLD.user_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION set_advert_status_to_sold()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE advert
    SET adverts_status = 'Продано'
    WHERE id = NEW.advert_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

--============= RK4 =================

CREATE OR REPLACE FUNCTION update_advert_promotion_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE advert
    SET is_promoted = true, promotion_start = now(), promotion_duration = NEW.promotion_duration
    WHERE id = NEW.advert_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


