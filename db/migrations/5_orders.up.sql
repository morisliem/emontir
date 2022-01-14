CREATE TABLE IF NOT EXISTS "orders" (
  "id" UUID NOT NULL,
  "description" VARCHAR(256),
  "total_price" FLOAT,
  "location" VARCHAR(256),
  "created_at" TIMESTAMP,
  "status" VARCHAR(64),
  "user_id" UUID NOT NULL,
  "motor_cycle_brand_name" VARCHAR(128) NOT NULL,
  "time_slot" VARCHAR(36),
  "date" DATE,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE
) WITHOUT OIDS;

CREATE INDEX IF NOT EXISTS "order_location" ON "orders" ("location");
CREATE INDEX IF NOT EXISTS "order_date" ON "orders" ("date");
CREATE INDEX IF NOT EXISTS "order_time_slot" ON "orders" ("time_slot");
CREATE INDEX IF NOT EXISTS "order_motor_cycle_brand_name" ON "orders" ("motor_cycle_brand_name");