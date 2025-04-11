-- +goose Up
-- modify "passwords" table
ALTER TABLE "passwords" DROP COLUMN "name";

-- +goose Down
-- reverse: modify "passwords" table
ALTER TABLE "passwords" ADD COLUMN "name" character varying(255) NOT NULL;
