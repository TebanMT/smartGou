-- Create the budgets table
CREATE TABLE budgets (
    budget_id UUID NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 4) NOT NULL,
    color VARCHAR(10) NOT NULL DEFAULT '#000000',
    icon VARCHAR(255) NOT NULL DEFAULT '💰',
    limit_amount DECIMAL(10, 4) NOT NULL DEFAULT 0,
    date DATE NOT NULL,
    deadline DATE NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_budgets_user_id FOREIGN KEY (user_id) REFERENCES users(user_id)
);


-- Create a trigger to update the created_at and updated_at fields
CREATE OR REPLACE FUNCTION before_insert_budgets_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_budgets
BEFORE INSERT ON budgets
FOR EACH ROW
EXECUTE FUNCTION before_insert_budgets_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_budgets_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_budgets
BEFORE UPDATE ON budgets
FOR EACH ROW
EXECUTE FUNCTION before_update_budgets_update_timestamp();

-- Create the meta categories table
CREATE TABLE meta_categories (
    meta_category_id UUID NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    name_en VARCHAR(255) NOT NULL,
    name_es VARCHAR(255) NOT NULL,
    icon VARCHAR(255) NOT NULL DEFAULT '💰',
    color VARCHAR(10) NOT NULL DEFAULT '#000000',
    description TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



-- Create a trigger to update the created_at and updated_at fields
CREATE OR REPLACE FUNCTION before_insert_meta_categories_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_meta_categories
BEFORE INSERT ON meta_categories
FOR EACH ROW
EXECUTE FUNCTION before_insert_meta_categories_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_meta_categories_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_meta_categories
BEFORE UPDATE ON meta_categories
FOR EACH ROW
EXECUTE FUNCTION before_update_meta_categories_update_timestamp();

-- seed the meta categories table
INSERT INTO meta_categories (name_en, name_es, icon, color, description) VALUES
('Food', 'Comida', '🍔', '#FF0000', 'Food is a necessity for survival and enjoyment.'),
('Groceries', 'Mercado', '🛍️', '#FF0000', 'Groceries are a necessity for survival and enjoyment.'),
('Transportation', 'Transporte', '🚗', '#00FF00', 'Transportation is a necessity for survival and enjoyment.'),
('Housing', 'Hogar', '🏠', '#000000', 'Housing is a necessity for survival and enjoyment.'),
('Lifestyle', 'Vida Social', '🎉', '#0000FF', 'Lifestyle is a necessity for survival and enjoyment.'),
('Beauty', 'Belleza', '💄', '#000000', 'Beauty is a necessity for survival and enjoyment.'),
('Healthcare', 'Salud', '🏥', '#000000', 'Healthcare is a necessity for survival and enjoyment.'),
('Education', 'Educación', '🎓', '#000000', 'Education is a necessity for survival and enjoyment.'),
('Gifts', 'Regalos', '🎁', '#000000', 'Gifts are a necessity for survival and enjoyment.'),
('Pets', 'Mascotas', '🐶', '#000000', 'Pets are a necessity for survival and enjoyment.'),
('Subscriptions', 'Suscripciones', '📺', '#000000', 'Subscriptions are a necessity for survival and enjoyment.'),
('Travel', 'Viajes', '✈️', '#000000', 'Travel is a necessity for survival and enjoyment.'),
('Luxury', 'Lujo', '💎', '#000000', 'Luxury is a necessity for survival and enjoyment.'),
('Kids', 'Niños', '👶', '#000000', 'Kids are a necessity for survival and enjoyment.'),
('Gambling', 'Juego', '🎲', '#000000', 'Gambling is a necessity for survival and enjoyment.'),
('Savings', 'Ahorro', '💰', '#000000', 'Savings are a necessity for survival and enjoyment.'),
('Gaming', 'Videojuegos', '🎮', '#000000', 'Gaming is a necessity for survival and enjoyment.'),
('Other', 'Otros', '💰', '#000000', 'Other is a necessity for survival and enjoyment.');

-- Create the spending categories table
CREATE TABLE spending_categories (
    spending_category_id UUID NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    name_en VARCHAR(255) NOT NULL,
    name_es VARCHAR(255) NOT NULL,
    icon VARCHAR(255) NOT NULL DEFAULT '💰',
    color VARCHAR(10) NOT NULL DEFAULT '#000000',
    description TEXT NULL,
    meta_category_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_spending_categories_meta_category_id FOREIGN KEY (meta_category_id) REFERENCES meta_categories(meta_category_id)
);



-- Create a trigger to update the created_at and updated_at fields
CREATE OR REPLACE FUNCTION before_insert_spending_categories_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_spending_categories
BEFORE INSERT ON spending_categories
FOR EACH ROW
EXECUTE FUNCTION before_insert_spending_categories_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_spending_categories_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_spending_categories
BEFORE UPDATE ON spending_categories
FOR EACH ROW
EXECUTE FUNCTION before_update_spending_categories_update_timestamp();


-- seed the spending categories table
INSERT INTO spending_categories (name_en, name_es, meta_category_id, icon, color, description) VALUES
-- spending categories for meta category food
('Groceries', 'Mercado', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Food'), '🛍️', '#FF0000', 'Groceries are a necessity for survival and enjoyment.'),
('Restaurants', 'Restaurantes', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Food'), '🍴', '#FF0000', 'Restaurants are a necessity for survival and enjoyment.'),

-- spending categories for meta category transportation
('Fuel', 'Combustible', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Transportation'), '🚗', '#00FF00', 'Fuel is a necessity for survival and enjoyment.'),
('Public Transportation', 'Transporte Público', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Transportation'), '🚌', '#00FF00', 'Public Transportation is a necessity for survival and enjoyment.'),
('Parking', 'Estacionamiento', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Transportation'), '🅿️', '#00FF00', 'Parking is a necessity for survival and enjoyment.'),

-- spending categories for meta category housing
('Rent', 'Alquiler', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Housing'), '🏠', '#000000', 'Rent is a necessity for survival and enjoyment.'),
('Mortgage', 'Hipoteca', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Housing'), '🏠', '#000000', 'Mortgage is a necessity for survival and enjoyment.'),
('Utilities', 'Servicios', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Housing'), '🔌', '#000000', 'Utilities is a necessity for survival and enjoyment.'),
('Maintenance', 'Mantenimiento', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Housing'), '🔌', '#000000', 'Maintenance is a necessity for survival and enjoyment.'),
('Home Improvement', 'Mejora del Hogar', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Housing'), '🔌', '#000000', 'Home Improvement is a necessity for survival and enjoyment.'),
('Internet', 'Internet', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Housing'), '🔌', '#000000', 'Internet is a necessity for survival and enjoyment.'),
('Cable', 'Cable', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Housing'), '🔌', '#000000', 'Cable is a necessity for survival and enjoyment.'),
('Water', 'Agua', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Housing'), '🔌', '#000000', 'Water is a necessity for survival and enjoyment.'),
('Electricity', 'Electricidad', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Housing'), '🔌', '#000000', 'Electricity is a necessity for survival and enjoyment.'),
('Gas', 'Gas', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Housing'), '🔌', '#000000', 'Gas is a necessity for survival and enjoyment.'),
('Telephone', 'Teléfono', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Housing'), '🔌', '#000000', 'Telephone is a necessity for survival and enjoyment.'),

-- spending categories for meta category lifestyle
('Shopping', 'Compras', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Lifestyle'), '🛍️', '#FF00FF', 'Shopping is a necessity for survival and enjoyment.'),
('Utilities', 'Servicios', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Lifestyle'), '🔌', '#000000', 'Utilities is a necessity for survival and enjoyment.'),
('Gym', 'Gimnasio', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Lifestyle'), '🏃‍♂️', '#000000', 'Gym is a necessity for survival and enjoyment.'),
('Yoga', 'Yoga', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Lifestyle'), '🧘‍♂️', '#000000', 'Yoga is a necessity for survival and enjoyment.'),
('Meditation', 'Meditación', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Lifestyle'), '🧘‍♂️', '#000000', 'Meditation is a necessity for survival and enjoyment.'),
('Massage', 'Masaje', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Lifestyle'), '💆‍♂️', '#000000', 'Massage is a necessity for survival and enjoyment.'),
('Haircut', 'Corte de pelo', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Lifestyle'), '💆‍♂️', '#000000', 'Haircut is a necessity for survival and enjoyment.'),
('Clothing', 'Ropa', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Lifestyle'), '👕', '#000000', 'Clothing is a necessity for survival and enjoyment.'),
('Shoes', 'Zapatos', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Lifestyle'), '👟', '#000000', 'Shoes are a necessity for survival and enjoyment.'),
('Accessories', 'Accesorios', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Lifestyle'), '👜', '#000000', 'Accessories are a necessity for survival and enjoyment.'),

-- spending categories for meta category beauty
('Skincare', 'Cuidado de la piel', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Beauty'), '💄', '#000000', 'Skincare are a necessity for survival and enjoyment.'),
('Beauty', 'Belleza', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Beauty'), '💄', '#000000', 'Beauty are a necessity for survival and enjoyment.'),
('Hair', 'Pelo', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Beauty'), '💄', '#000000', 'Hair are a necessity for survival and enjoyment.'),
('Nails', 'Uñas', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Beauty'), '💄', '#000000', 'Nails are a necessity for survival and enjoyment.'),
('Makeup', 'Maquillaje', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Beauty'), '💄', '#000000', 'Makeup are a necessity for survival and enjoyment.'),
('Fragrance', 'Fragancia', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Beauty'), '💄', '#000000', 'Fragrance are a necessity for survival and enjoyment.'),

-- spending categories for meta category healthcare
('Insurance', 'Seguro', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Healthcare'), '🏥', '#000000', 'Insurance is a necessity for survival and enjoyment.'),
('Medication', 'Medicamento', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Healthcare'), '🏥', '#000000', 'Medication is a necessity for survival and enjoyment.'),
('Doctor', 'Doctor', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Healthcare'), '🏥', '#000000', 'Doctor is a necessity for survival and enjoyment.'),
('Pharmacy', 'Farmacia', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Healthcare'), '🏥', '#000000', 'Pharmacy is a necessity for survival and enjoyment.'),
('Dental', 'Dental', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Healthcare'), '🏥', '#000000', 'Dental is a necessity for survival and enjoyment.'),
('Vision', 'Vision', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Healthcare'), '🏥', '#000000', 'Vision is a necessity for survival and enjoyment.'),
('Medical Emergency', 'Emergencia Médica', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Healthcare'), '🏥', '#000000', 'Medical Emergency is a necessity for survival and enjoyment.'),

-- spending categories for meta category education
('Books', 'Libros', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Education'), '🎓', '#000000', 'Books are a necessity for survival and enjoyment.'),
('Courses', 'Cursos', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Education'), '🎓', '#000000', 'Courses are a necessity for survival and enjoyment.'),
('Tuition', 'Matrícula', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Education'), '🎓', '#000000', 'Tuition is a necessity for survival and enjoyment.'),
('School', 'Escuela', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Education'), '🎓', '#000000', 'School is a necessity for survival and enjoyment.'),
('Online Courses', 'Cursos Online', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Education'), '🎓', '#000000', 'Online Courses are a necessity for survival and enjoyment.'),
('Certifications', 'Certificaciones', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Education'), '🎓', '#000000', 'Certifications are a necessity for survival and enjoyment.'),

-- spending categories for meta category gifts
('Gift Cards', 'Tarjetas de Regalo', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gifts'), '🎁', '#000000', 'Gift Cards are a necessity for survival and enjoyment.'),
('Gift Baskets', 'Cestas de Regalo', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gifts'), '🎁', '#000000', 'Gift Baskets are a necessity for survival and enjoyment.'),
('Gift Certificates', 'Certificados de Regalo', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gifts'), '🎁', '#000000', 'Gift Certificates are a necessity for survival and enjoyment.'),
('Gift Wrap', 'Wrap de Regalo', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gifts'), '🎁', '#000000', 'Gift Wrap are a necessity for survival and enjoyment.'),
('Gift Bags', 'Bolsas de Regalo', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gifts'), '🎁', '#000000', 'Gift Bags are a necessity for survival and enjoyment.'),
('Gift Tags', 'Etiquetas de Regalo', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gifts'), '🎁', '#000000', 'Gift Tags are a necessity for survival and enjoyment.'),
('Gift Boxes', 'Cajas de Regalo', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gifts'), '🎁', '#000000', 'Gift Boxes are a necessity for survival and enjoyment.'),
('Gift Bowls', 'Cajas de Regalo', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gifts'), '🎁', '#000000', 'Gift Bowls are a necessity for survival and enjoyment.'),

-- spending categories for meta category pets
('Pet Food', 'Comida para Mascotas', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Pets'), '🐶', '#000000', 'Pet Food are a necessity for survival and enjoyment.'),
('Pet Supplies', 'Accesorios para Mascotas', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Pets'), '🐶', '#000000', 'Pet Supplies are a necessity for survival and enjoyment.'),
('Pet Grooming', 'Corte de Pelo para Mascotas', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Pets'), '🐶', '#000000', 'Pet Grooming are a necessity for survival and enjoyment.'),
('Pet Care', 'Cuidado de Mascotas', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Pets'), '🐶', '#000000', 'Pet Care are a necessity for survival and enjoyment.'),

-- spending categories for meta category subscriptions
('Streaming Services', 'Servicios de Streaming', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'Streaming Services are a necessity for survival and enjoyment.'),
('Netflix', 'Netflix', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'Netflix are a necessity for survival and enjoyment.'),
('Amazon Prime', 'Amazon Prime', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'Amazon Prime are a necessity for survival and enjoyment.'),
('Disney+', 'Disney+', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'Disney+ are a necessity for survival and enjoyment.'),
('Hulu', 'Hulu', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'Hulu are a necessity for survival and enjoyment.'),
('Apple Music', 'Apple Music', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'Apple Music are a necessity for survival and enjoyment.'),
('Apple TV+', 'Apple TV+', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'Apple TV+ are a necessity for survival and enjoyment.'),
('Apple Podcasts', 'Apple Podcasts', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'Apple Podcasts are a necessity for survival and enjoyment.'),
('Apple News+', 'Apple News+', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'Apple News+ are a necessity for survival and enjoyment.'),
('Apple Arcade', 'Apple Arcade', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'Apple Arcade are a necessity for survival and enjoyment.'),
('Spotify', 'Spotify', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'Spotify are a necessity for survival and enjoyment.'),
('Tinder', 'Tinder', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'Tinder are a necessity for survival and enjoyment.'),
('HBO Max', 'HBO Max', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Subscriptions'), '📺', '#000000', 'HBO Max are a necessity for survival and enjoyment.'),

-- spending categories for meta category travel
('Hotels', 'Hoteles', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Travel'), '✈️', '#000000', 'Hotels are a necessity for survival and enjoyment.'),
('Flights', 'Vuelos', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Travel'), '✈️', '#000000', 'Flights are a necessity for survival and enjoyment.'),
('Car Rentals', 'Alquiler de Autos', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Travel'), '✈️', '#000000', 'Car Rentals are a necessity for survival and enjoyment.'),

-- spending categories for meta category luxury
('Jewelry', 'Joyas', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Luxury'), '💎', '#000000', 'Jewelry is a necessity for survival and enjoyment.'),

-- spending categories for meta category kids
('Toys', 'Juguetes', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Kids'), '👶', '#000000', 'Toys are a necessity for survival and enjoyment.'),
('School Supplies', 'Materiales Escolares', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Kids'), '👶', '#000000', 'School Supplies are a necessity for survival and enjoyment.'),

-- spending categories for meta category gambling
('Casino', 'Casino', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gambling'), '🎲', '#000000', 'Casino is a necessity for survival and enjoyment.'),
('Lottery', 'Lotería', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gambling'), '🎲', '#000000', 'Lottery is a necessity for survival and enjoyment.'),
('Sports Betting', 'Apuestas Deportivas', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gambling'), '🎲', '#000000', 'Sports Betting is a necessity for survival and enjoyment.'),

-- spending categories for meta category savings
('Emergency Fund', 'Fondo de Emergencia', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Savings'), '💰', '#000000', 'Emergency Fund is a necessity for survival and enjoyment.'),
('Retirement', 'Retiro', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Savings'), '💰', '#000000', 'Retirement is a necessity for survival and enjoyment.'),
('Investments', 'Inversiones', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Savings'), '💰', '#000000', 'Investments are a necessity for survival and enjoyment.'),

-- spending categories for meta category gaming
('Gaming Consoles', 'Consolas de Videojuegos', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gaming'), '🎮', '#000000', 'Gaming Consoles are a necessity for survival and enjoyment.'),
('Gaming Accessories', 'Accesorios de Videojuegos', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gaming'), '🎮', '#000000', 'Gaming Accessories are a necessity for survival and enjoyment.'),
('Gaming Software', 'Software de Videojuegos', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gaming'), '🎮', '#000000', 'Gaming Software are a necessity for survival and enjoyment.'),
('Video Games', 'Videojuegos', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gaming'), '🎮', '#000000', 'Video Games are a necessity for survival and enjoyment.'),
('Online Games', 'Videojuegos en línea', (SELECT meta_category_id FROM meta_categories WHERE name_en = 'Gaming'), '🎮', '#000000', 'Online Games are a necessity for survival and enjoyment.');




-- Create the budget categories table
CREATE TABLE budget_categories (
    budget_category_id UUID NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID NOT NULL,
    spending_category_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_budget_categories_budget_id FOREIGN KEY (budget_id) REFERENCES budgets(budget_id),
    CONSTRAINT fk_budget_categories_spending_category_id FOREIGN KEY (spending_category_id) REFERENCES spending_categories(spending_category_id)
);

-- Create a trigger to update the created_at and updated_at fields
CREATE OR REPLACE FUNCTION before_insert_budget_categories_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_budget_categories
BEFORE INSERT ON budget_categories
FOR EACH ROW
EXECUTE FUNCTION before_insert_budget_categories_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_budget_categories_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_budget_categories
BEFORE UPDATE ON budget_categories
FOR EACH ROW
EXECUTE FUNCTION before_update_budget_categories_update_timestamp();


-- create the user categories preferences table
CREATE TABLE user_categories_preferences (
    user_category_preference_id UUID NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    spending_category_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_categories_preferences_user_id FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT fk_user_categories_preferences_spending_category_id FOREIGN KEY (spending_category_id) REFERENCES spending_categories(spending_category_id)
);

-- Create a trigger to update the created_at and updated_at fields
CREATE OR REPLACE FUNCTION before_insert_user_categories_preferences_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_user_categories_preferences
BEFORE INSERT ON user_categories_preferences
FOR EACH ROW
EXECUTE FUNCTION before_insert_user_categories_preferences_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_user_categories_preferences_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_user_categories_preferences
BEFORE UPDATE ON user_categories_preferences
FOR EACH ROW
EXECUTE FUNCTION before_update_user_categories_preferences_update_timestamp();
