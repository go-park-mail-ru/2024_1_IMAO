CREATE OR REPLACE FUNCTION increment_views()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE advert
    SET views = views + 1
    WHERE id = NEW.advert_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
