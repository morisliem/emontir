ALTER TABLE "order_items" 
    DROP CONSTRAINT IF EXISTS "fk_order_id",
    ADD CONSTRAINT "fk_order_id" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON DELETE CASCADE;
