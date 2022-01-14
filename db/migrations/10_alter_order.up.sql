DROP INDEX IF EXISTS "order_location";
ALTER TABLE "orders" 
    DROP COLUMN IF EXISTS "location",
    ADD COLUMN "user_address_id" UUID NOT NULL,
    ADD CONSTRAINT "fk_user_address_id" FOREIGN KEY ("user_address_id") REFERENCES "user_addresses" ("id");
