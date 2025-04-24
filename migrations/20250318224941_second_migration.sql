-- +goose Up
-- modify "cards" table
ALTER TABLE "cards" ADD COLUMN "hashed_card_number" text NULL, ADD CONSTRAINT "cards_hashed_card_number_key" UNIQUE ("hashed_card_number");

-- +goose Down
-- reverse: modify "cards" table
ALTER TABLE "cards" DROP CONSTRAINT "cards_hashed_card_number_key", DROP COLUMN "hashed_card_number";
