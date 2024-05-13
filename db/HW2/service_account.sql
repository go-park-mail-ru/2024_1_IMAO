CREATE USER vol4ok_service_account WITH PASSWORD 'vol4ok-Password123!';

CREATE ROLE vol4ok_service_account_role WITH LOGIN;

GRANT USAGE ON SCHEMA public TO vol4ok_service_account_role;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO vol4ok_service_account_role;
GRANT INSERT ON advert, favourite, cart, advert_image, answer, blacklist, complaint, "order", profile, "user", "subscription",
survey, user_survey, "view" TO vol4ok_service_account_role;
GRANT DELETE ON favourite, cart, advert_image, "subscription", blacklist TO vol4ok_service_account_role;
GRANT UPDATE ON advert, profile, "user" TO vol4ok_service_account_role;

GRANT vol4ok_service_account_role TO vol4ok_service_account;