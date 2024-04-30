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
