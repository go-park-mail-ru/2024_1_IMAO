CREATE TRIGGER view_insert_trigger
AFTER INSERT ON view
FOR EACH ROW
EXECUTE PROCEDURE increment_views();
