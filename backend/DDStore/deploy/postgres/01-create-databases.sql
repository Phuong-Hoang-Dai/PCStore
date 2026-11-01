-- Runs once, on the first start of an empty postgres volume.
-- ddstore_users is created by POSTGRES_DB; each service owns its own database.
CREATE DATABASE ddstore_orders;
