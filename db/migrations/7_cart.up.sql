CREATE TABLE IF NOT EXISTS "carts"(
    "id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "date" DATE NOT NULL,
    "description" VARCHAR(256),
    "time_slot" VARCHAR(36) NOT NULL,
    "user_address_id" UUID,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_user_address_id" FOREIGN KEY ("user_address_id") REFERENCES "user_addresses" ("id") ON DELETE CASCADE
) WITHOUT OIDS;