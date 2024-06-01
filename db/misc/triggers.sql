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

CREATE OR REPLACE TRIGGER advert_status_update_trigger
AFTER INSERT ON "order"
FOR EACH ROW
EXECUTE PROCEDURE set_advert_status_to_sold();

--============= RK4 ===========================
CREATE TRIGGER update_advert_promotion
AFTER UPDATE ON payments
FOR EACH ROW
EXECUTE PROCEDURE update_advert_promotion_func()

--============= DEFENCE ===========================

CREATE OR REPLACE TRIGGER profile_rating_update_trigger
AFTER INSERT ON review
FOR EACH ROW
EXECUTE PROCEDURE update_profile_rating();
