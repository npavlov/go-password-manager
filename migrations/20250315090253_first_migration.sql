-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- create enum type "item_type"
CREATE TYPE "item_type" AS ENUM ('password', 'binary', 'card', 'text');

-- create "users" table
CREATE TABLE "users" (
                         "id" uuid NOT NULL DEFAULT gen_random_uuid(),
                         "username" character varying(255) NOT NULL,
                         "email" character varying(255) NOT NULL,
                         "password" text NOT NULL,
                         "encryption_key" text NOT NULL,
                         PRIMARY KEY ("id"),
                         CONSTRAINT "users_email_key" UNIQUE ("email"),
                         CONSTRAINT "users_username_key" UNIQUE ("username")
);

-- create "binary_entries" table
CREATE TABLE "binary_entries" (
                                  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
                                  "user_id" uuid NULL,
                                  "file_name" character varying(255) NOT NULL,
                                  "file_size" bigint NOT NULL,
                                  "file_type" character varying(255) NOT NULL,
                                  "file_url" text NOT NULL,
                                  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                  PRIMARY KEY ("id"),
                                  CONSTRAINT "binary_entries_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- create "cards" table
CREATE TABLE "cards" (
                         "id" uuid NOT NULL DEFAULT gen_random_uuid(),
                         "user_id" uuid NULL,
                         "encrypted_card_number" text NOT NULL,
                         "encrypted_expiry_date" text NOT NULL,
                         "encrypted_cvv" text NOT NULL,
                         "cardholder_name" character varying(255) NOT NULL,
                         "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         PRIMARY KEY ("id"),
                         CONSTRAINT "cards_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- create "item" table
CREATE TABLE "items" (
                         "id" uuid NOT NULL DEFAULT gen_random_uuid(),
                         "user_id" uuid NULL,
                         "type" "item_type" NOT NULL,
                         "id_resource" uuid NOT NULL,
                         "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         PRIMARY KEY ("id"),
                         CONSTRAINT "items_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- create "notes" table
CREATE TABLE "notes" (
                         "id" uuid NOT NULL DEFAULT gen_random_uuid(),
                         "user_id" uuid NULL,
                         "encrypted_content" text NOT NULL,
                         "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         PRIMARY KEY ("id"),
                         CONSTRAINT "notes_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- create "passwords" table
CREATE TABLE "passwords" (
                             "id" uuid NOT NULL DEFAULT gen_random_uuid(),
                             "user_id" uuid NULL,
                             "name" character varying(255) NOT NULL,
                             "login" character varying(255) NOT NULL,
                             "password" text NOT NULL,
                             "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             PRIMARY KEY ("id"),
                             CONSTRAINT "unique_user_password_name" UNIQUE ("user_id", "name"),
                             CONSTRAINT "passwords_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Create the update_updated_at_column function
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- create triggers to update updated_at column
CREATE TRIGGER set_updated_at_passwords
    BEFORE UPDATE ON passwords
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_updated_at_cards
    BEFORE UPDATE ON cards
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_updated_at_binary_entries
    BEFORE UPDATE ON binary_entries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_updated_at_notes
    BEFORE UPDATE ON notes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- create function to add item to the item table
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION add_note_to_items()
RETURNS TRIGGER AS $$
BEGIN
INSERT INTO items (user_id, type, id_resource)
VALUES (NEW.user_id, 'text', NEW.id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION add_card_to_items()
RETURNS TRIGGER AS $$
BEGIN
INSERT INTO items (user_id, type, id_resource)
VALUES (NEW.user_id, 'card', NEW.id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION add_binary_to_items()
RETURNS TRIGGER AS $$
BEGIN
INSERT INTO items (user_id, type, id_resource)
VALUES (NEW.user_id, 'binary', NEW.id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION add_password_to_items()
RETURNS TRIGGER AS $$
BEGIN
INSERT INTO items (user_id, type, id_resource)
VALUES (NEW.user_id, 'password', NEW.id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- create triggers to add item to the item table
CREATE TRIGGER add_note_to_items
    AFTER INSERT ON notes
    FOR EACH ROW
    EXECUTE FUNCTION add_note_to_items();

CREATE TRIGGER add_card_to_items
    AFTER INSERT ON cards
    FOR EACH ROW
    EXECUTE FUNCTION add_card_to_items();

CREATE TRIGGER add_binary_to_items
    AFTER INSERT ON binary_entries
    FOR EACH ROW
    EXECUTE FUNCTION add_binary_to_items();

CREATE TRIGGER add_password_to_items
    AFTER INSERT ON passwords
    FOR EACH ROW
    EXECUTE FUNCTION add_password_to_items();

-- +goose Down
-- drop triggers
DROP TRIGGER IF EXISTS add_note_to_items ON notes;
DROP TRIGGER IF EXISTS add_card_to_items ON cards;
DROP TRIGGER IF EXISTS add_binary_to_items ON binary_entries;
DROP TRIGGER IF EXISTS add_password_to_items ON passwords;

DROP TRIGGER IF EXISTS set_updated_at_passwords ON passwords;
DROP TRIGGER IF EXISTS set_updated_at_cards ON cards;
DROP TRIGGER IF EXISTS set_updated_at_binary_entries ON binary_entries;
DROP TRIGGER IF EXISTS set_updated_at_notes ON notes;

-- drop functions
DROP FUNCTION IF EXISTS add_note_to_items();
DROP FUNCTION IF EXISTS add_card_to_items();
DROP FUNCTION IF EXISTS add_binary_to_items();
DROP FUNCTION IF EXISTS add_password_to_items();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- drop tables
DROP TABLE "passwords";
DROP TABLE "notes";
DROP TABLE "items";
DROP TABLE "cards";
DROP TABLE "binary_entries";
DROP TABLE "users";

-- drop enum type
DROP TYPE "item_type";