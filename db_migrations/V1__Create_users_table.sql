
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    username VARCHAR(100),
    name VARCHAR(100),
    last_name VARCHAR(100),
    second_last_name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    dailing_code VARCHAR(100),
    phone VARCHAR(100) UNIQUE,
    password VARCHAR(250),
    verified_phone BOOLEAN DEFAULT FALSE,
    verified_email BOOLEAN DEFAULT FALSE,
    is_onboarding_completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION before_insert_user_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_user
BEFORE INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION before_insert_user_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_user_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_user
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION before_update_user_update_timestamp();