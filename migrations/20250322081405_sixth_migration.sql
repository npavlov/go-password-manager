-- +goose Up
-- modify "items" table
ALTER TABLE "items" ADD CONSTRAINT "items_id_resource_key" UNIQUE ("id_resource");
-- modify "metainfo" table
ALTER TABLE "metainfo" DROP CONSTRAINT "metainfo_item_id_fkey", ADD
 CONSTRAINT "metainfo_item_id_fkey" FOREIGN KEY ("item_id") REFERENCES "items" ("id_resource") ON UPDATE NO ACTION ON DELETE CASCADE;

-- +goose Down
-- reverse: modify "metainfo" table
ALTER TABLE "metainfo" DROP CONSTRAINT "metainfo_item_id_fkey", ADD
 CONSTRAINT "metainfo_item_id_fkey" FOREIGN KEY ("item_id") REFERENCES "items" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- reverse: modify "items" table
ALTER TABLE "items" DROP CONSTRAINT "items_id_resource_key";
