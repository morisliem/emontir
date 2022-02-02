CREATE TABLE IF NOT EXISTS "cart_items"(
    "id" SERIAL NOT NULL,
    "cart_id" UUID NOT NULL,
    "service_id" SERIAL NOT NULL,
    PRIMARY KEY("id"),
    CONSTRAINT "fk_cart_id" FOREIGN KEY ("cart_id") REFERENCES "carts" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT "fk_service_id" FOREIGN KEY ("service_id") REFERENCES "services" ("id")
);