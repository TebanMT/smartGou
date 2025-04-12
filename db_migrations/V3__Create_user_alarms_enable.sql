
-- Create a table to store the alarms
CREATE TABLE reminders (
    reminder_id UUID NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    name_es VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    description_es VARCHAR(255) NOT NULL,
    description_en VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create a trigger to update the created_at and updated_at fields
CREATE OR REPLACE FUNCTION before_insert_reminders_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_reminders
BEFORE INSERT ON reminders
FOR EACH ROW
EXECUTE FUNCTION before_insert_reminders_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_reminders_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_reminders
BEFORE UPDATE ON reminders
FOR EACH ROW
EXECUTE FUNCTION before_update_reminders_update_timestamp();

-- seed data
INSERT INTO reminders (name_es, name_en, description_es, description_en, is_active) VALUES ('Gastos diarios', 'Daily Expenses', 'Recibe una notificación cada vez que gastas dinero', 'Receive a notification every time you spend money', TRUE);
INSERT INTO reminders (name_es, name_en, description_es, description_en, is_active) VALUES ('Recordatorio de facturas', 'Bill Reminder', 'Le notificaremos cuando sea el momento de pagar sus facturas', 'We will nudge you when it is time to pay your bills', TRUE);
INSERT INTO reminders (name_es, name_en, description_es, description_en, is_active) VALUES ('Resumen semanal', 'Weekly Summary', 'Le notificaremos cuando esté listo su resumen semanal', 'We will nudge you when your weekly summary is ready', TRUE);

-- Create a table to store the user alarms enable
CREATE TABLE user_reminders (
    user_reminders_id UUID NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    reminder_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_reminders_user_id FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT fk_user_reminders_reminder_id FOREIGN KEY (reminder_id) REFERENCES reminders(reminder_id)
);

-- Create a trigger to update the created_at and updated_at fields
CREATE OR REPLACE FUNCTION before_insert_user_reminders_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_user_reminders
BEFORE INSERT ON user_reminders
FOR EACH ROW
EXECUTE FUNCTION before_insert_user_reminders_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_user_reminders_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_user_reminders
BEFORE UPDATE ON user_reminders
FOR EACH ROW
EXECUTE FUNCTION before_update_user_reminders_update_timestamp();
