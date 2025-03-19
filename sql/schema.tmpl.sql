-- Enable the pgcrypto extension to use gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create ENUM type for item categories
CREATE TYPE item_type AS ENUM ('password', 'binary', 'card', 'text');

-- Create users table
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       username VARCHAR(255) UNIQUE NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password TEXT NOT NULL,
                       encryption_key TEXT NOT NULL
);

-- Create orders table
CREATE TABLE passwords (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                        name VARCHAR(255) NOT NULL,
                        login VARCHAR(255) NOT NULL,
                        password TEXT NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        CONSTRAINT unique_user_password_name UNIQUE (user_id, name)
);

CREATE TABLE binary_entries (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                        file_name VARCHAR(255) NOT NULL,
                        file_size BIGINT NOT NULL,  -- File size in bytes
                        file_url TEXT NOT NULL,  -- URL to MinIO storage
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cards (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                       encrypted_card_number TEXT NOT NULL,  -- Encrypted
                       hashed_card_number TEXT UNIQUE,
                       encrypted_expiry_date TEXT NOT NULL,  -- Encrypted
                       encrypted_cvv TEXT NOT NULL,          -- Encrypted
                       cardholder_name VARCHAR(255) NOT NULL, -- Cardholder name is stored in plaintext
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notes (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                       encrypted_content TEXT NOT NULL, -- Encrypted
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table that links different types of stored item
CREATE TABLE items (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                       type item_type NOT NULL,
                       id_resource UUID NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

DROP FUNCTION IF EXISTS update_updated_at_column();
-- Add trigger function to update the updated_at column
CREATE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for orders table to automatically update updated_at
CREATE TRIGGER set_updated_at_passwords
    BEFORE UPDATE ON passwords
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create trigger for orders table to automatically update updated_at
CREATE TRIGGER set_updated_at_cards
    BEFORE UPDATE ON cards
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create trigger for orders table to automatically update updated_at
CREATE TRIGGER set_updated_at_binary_entries
    BEFORE UPDATE ON binary_entries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create trigger for orders table to automatically update updated_at
CREATE TRIGGER set_updated_at_notes
    BEFORE UPDATE ON notes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();


DROP FUNCTION IF EXISTS add_note_to_items();
DROP FUNCTION IF EXISTS add_card_to_items();
DROP FUNCTION IF EXISTS add_binary_to_items();
DROP FUNCTION IF EXISTS add_password_to_items();

-- Function to insert notes into item
CREATE FUNCTION add_note_to_items() RETURNS TRIGGER AS $$
BEGIN
INSERT INTO items (user_id, type, id_resource)
VALUES (NEW.user_id, 'text', NEW.id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to insert cards into item
CREATE FUNCTION add_card_to_items() RETURNS TRIGGER AS $$
BEGIN
INSERT INTO items (user_id, type, id_resource)
VALUES (NEW.user_id, 'card', NEW.id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to insert binary entries into item
CREATE FUNCTION add_binary_to_items() RETURNS TRIGGER AS $$
BEGIN
INSERT INTO items (user_id, type, id_resource)
VALUES (NEW.user_id, 'binary', NEW.id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to insert passwords into item
CREATE FUNCTION add_password_to_items() RETURNS TRIGGER AS $$
BEGIN
INSERT INTO items (user_id, type, id_resource)
VALUES (NEW.user_id, 'password', NEW.id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS add_note_to_items ON notes;
DROP TRIGGER IF EXISTS add_card_to_items ON cards;
DROP TRIGGER IF EXISTS add_binary_to_items ON binary_entries;
DROP TRIGGER IF EXISTS add_password_to_items ON passwords;

-- Trigger for notes
CREATE TRIGGER add_note_to_items
    AFTER INSERT ON notes
    FOR EACH ROW EXECUTE FUNCTION add_note_to_items();

-- Trigger for cards
CREATE TRIGGER add_card_to_items
    AFTER INSERT ON cards
    FOR EACH ROW EXECUTE FUNCTION add_card_to_items();

-- Trigger for binary entries
CREATE TRIGGER add_binary_to_items
    AFTER INSERT ON binary_entries
    FOR EACH ROW EXECUTE FUNCTION add_binary_to_items();

-- Trigger for passwords
CREATE TRIGGER add_password_to_items
    AFTER INSERT ON passwords
    FOR EACH ROW EXECUTE FUNCTION add_password_to_items();