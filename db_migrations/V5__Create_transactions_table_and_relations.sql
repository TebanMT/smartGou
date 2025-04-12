-- create transactions table
CREATE TYPE transaction_type AS ENUM ('INCOME', 'EXPENSE', 'TRANSFER');
CREATE TABLE transactions (
    transaction_id UUID NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    amount DECIMAL(10, 4) NOT NULL,
    note TEXT NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    type transaction_type NOT NULL,
    category_id uuid NOT NULL,
    budget_id uuid NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_transactions_user_id FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT fk_transactions_category_id FOREIGN KEY (category_id) REFERENCES spending_categories(spending_category_id),
    CONSTRAINT fk_transactions_budget_id FOREIGN KEY (budget_id) REFERENCES budgets(budget_id)
);


-- Create a trigger to update the created_at and updated_at fields
CREATE OR REPLACE FUNCTION before_insert_transactions_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_transactions
BEFORE INSERT ON transactions
FOR EACH ROW
EXECUTE FUNCTION before_insert_transactions_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_transactions_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_transactions
BEFORE UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION before_update_transactions_update_timestamp();

-- create the expenses table
CREATE TABLE expenses (
    expense_id SERIAL PRIMARY KEY,
    transaction_id uuid NOT NULL,
    merchant TEXT,
    receipt_url TEXT,
    CONSTRAINT fk_expenses_transaction_id FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id)
);

-- Create a trigger to update the created_at and updated_at fields
CREATE OR REPLACE FUNCTION before_insert_expenses_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_expenses
BEFORE INSERT ON expenses
FOR EACH ROW
EXECUTE FUNCTION before_insert_expenses_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_expenses_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_expenses
BEFORE UPDATE ON expenses
FOR EACH ROW
EXECUTE FUNCTION before_update_expenses_update_timestamp();

-- create the income table
CREATE TABLE income (
    income_id SERIAL PRIMARY KEY,
    transaction_id uuid NOT NULL,
    source TEXT,
    CONSTRAINT fk_income_transaction_id FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id)
);

-- Create a trigger to update the created_at and updated_at fields
CREATE OR REPLACE FUNCTION before_insert_income_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_income
BEFORE INSERT ON income
FOR EACH ROW
EXECUTE FUNCTION before_insert_income_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_income_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_income
BEFORE UPDATE ON income
FOR EACH ROW
EXECUTE FUNCTION before_update_income_update_timestamp();


-- create the accounts table
CREATE TABLE accounts (
    account_id UUID NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    balance DECIMAL(10, 4) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create a trigger to update the created_at and updated_at fields
CREATE OR REPLACE FUNCTION before_insert_accounts_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_accounts
BEFORE INSERT ON accounts
FOR EACH ROW
EXECUTE FUNCTION before_insert_accounts_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_accounts_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_accounts
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE FUNCTION before_update_accounts_update_timestamp();

-- create the transfers table
CREATE TABLE transfers (
    transfer_id UUID NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id uuid NOT NULL,
    from_account_id uuid NOT NULL,
    to_account_id uuid NOT NULL,
    CONSTRAINT fk_transfers_transaction_id FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id),
    CONSTRAINT fk_transfers_from_account_id FOREIGN KEY (from_account_id) REFERENCES accounts(account_id),
    CONSTRAINT fk_transfers_to_account_id FOREIGN KEY (to_account_id) REFERENCES accounts(account_id)
);

-- Create a trigger to update the created_at and updated_at fields
CREATE OR REPLACE FUNCTION before_insert_transfers_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_transfers
BEFORE INSERT ON transfers
FOR EACH ROW
EXECUTE FUNCTION before_insert_transfers_update_timestamp();

CREATE OR REPLACE FUNCTION before_update_transfers_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_transfers
BEFORE UPDATE ON transfers
FOR EACH ROW
EXECUTE FUNCTION before_update_transfers_update_timestamp();

