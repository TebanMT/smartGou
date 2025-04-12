create table user_usage_purpose (
    usage_purpose_id UUID NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    description_es VARCHAR(255) NOT NULL,
    description_en VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE OR REPLACE FUNCTION before_insert_user_usage_purpose_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_user_usage_purpose
BEFORE INSERT ON user_usage_purpose
FOR EACH ROW
EXECUTE FUNCTION before_insert_user_usage_purpose_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_user_usage_purpose_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_user_usage_purpose
BEFORE UPDATE ON user_usage_purpose
FOR EACH ROW
EXECUTE FUNCTION before_update_user_usage_purpose_update_timestamp();

-- Add a column to the users table
ALTER TABLE users ADD COLUMN user_usage_purpose_id UUID;
-- Add a foreign key constraint to the users table
ALTER TABLE users ADD CONSTRAINT fk_user_usage_purpose_id FOREIGN KEY (user_usage_purpose_id) REFERENCES user_usage_purpose(usage_purpose_id);

-- seed data
INSERT INTO user_usage_purpose (description_es, description_en) VALUES ('Hacer un presupuesto', 'Make a budget');
INSERT INTO user_usage_purpose (description_es, description_en) VALUES ('Seguir mis gastos', 'Track my spending');
INSERT INTO user_usage_purpose (description_es, description_en) VALUES ('Ahorrar para un objetivo', 'Save for a goal');
INSERT INTO user_usage_purpose (description_es, description_en) VALUES ('Pagar mis deudas', 'Pay off debts');
INSERT INTO user_usage_purpose (description_es, description_en) VALUES ('Compartir con mi pareja', 'Share with my partner');
INSERT INTO user_usage_purpose (description_es, description_en) VALUES ('Sincronizar entre dispositivos', 'Sync between devices');
INSERT INTO user_usage_purpose (description_es, description_en) VALUES ('Otro', 'Other');
