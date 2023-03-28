CREATE TABLE stores (
    id serial PRIMARY KEY,
    some_field VARCHAR ( 255 ) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON stores
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();