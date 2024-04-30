CREATE TRIGGER view_insert_trigger
AFTER INSERT ON view
FOR EACH ROW
EXECUTE PROCEDURE increment_views();

CREATE TRIGGER cart_insert_trigger
AFTER INSERT ON cart
FOR EACH ROW
EXECUTE PROCEDURE increment_cart_adverts_number();

CREATE TRIGGER cart_delete_trigger
AFTER DELETE ON cart
FOR EACH ROW
EXECUTE PROCEDURE decrement_cart_adverts_number();

CREATE OR REPLACE TRIGGER fav_insert_trigger
AFTER INSERT ON favourite
FOR EACH ROW
EXECUTE PROCEDURE increment_fav_adverts_number();

CREATE OR REPLACE TRIGGER fav_delete_trigger
AFTER DELETE ON favourite
FOR EACH ROW
EXECUTE PROCEDURE decrement_fav_adverts_number();
