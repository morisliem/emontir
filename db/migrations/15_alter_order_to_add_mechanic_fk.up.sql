ALTER TABLE "orders" 
    DROP COLUMN IF EXISTS "status",
    ADD COLUMN "status_order" VARCHAR(64),
    ADD COLUMN "status_detail" VARCHAR(64),
    ADD COLUMN "mechanic_id" INT,
    ADD CONSTRAINT "fk_mechanic_id" FOREIGN KEY ("mechanic_id") REFERENCES "mechanics" ("id");

ALTER TABLE "order_items"
    DROP CONSTRAINT IF EXISTS "fk_service_id",
    DROP COLUMN IF EXISTS "service_id",
    ADD COLUMN "service_id" INT NOT NULL,
    ADD CONSTRAINT "fk_service_id" FOREIGN KEY ("service_id") REFERENCES "services" ("id");