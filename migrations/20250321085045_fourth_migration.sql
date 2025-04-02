-- +goose Up
-- create "metainfo" table
CREATE TABLE "metainfo" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "item_id" uuid NULL,
  "key" character varying(255) NOT NULL,
  "value" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "unique_item_key" UNIQUE ("item_id", "key"),
  CONSTRAINT "metainfo_item_id_fkey" FOREIGN KEY ("item_id") REFERENCES "items" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- +goose Down
-- reverse: create "metainfo" table
DROP TABLE "metainfo";
