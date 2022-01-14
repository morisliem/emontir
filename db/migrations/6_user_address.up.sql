CREATE TABLE IF NOT EXISTS "user_addresses"(
    "id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "label" VARCHAR(128),
    "address" VARCHAR(256),
    "address_detail" VARCHAR(256),
    "phone_num" VARCHAR(16),
    "recipient_name" VARCHAR(128),
    "latitude" VARCHAR(128),
    "longitude" VARCHAR(128),
    "created_at" TIMESTAMP NOT NULL,
    PRIMARY KEY("id"),
    CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE
) WITHOUT OIDS;