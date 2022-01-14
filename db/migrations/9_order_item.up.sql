CREATE TABLE IF NOT EXISTS "order_items"(
    "id" SERIAL NOT NULL,
    "order_id" UUID NOT NULL,
    "service_id" SERIAL NOT NULL,
    PRIMARY KEY("id"),
    CONSTRAINT "fk_order_id" FOREIGN KEY ("order_id") REFERENCES "orders" ("id"),
    CONSTRAINT "fk_service_id" FOREIGN KEY ("service_id") REFERENCES "services" ("id")
)WITHOUT OIDS;