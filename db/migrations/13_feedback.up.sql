CREATE TABLE IF NOT EXISTS "feedbacks"(
    "id" SERIAL NOT NULL,
    "feedback" VARCHAR(300),
    "rating" FLOAT NOT NULL,
    "service_id" INT NOT NULL,
    "user_id" UUID NOT NULL,
    "order_id" UUID NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT uniq UNIQUE ("service_id","order_id"),
    CONSTRAINT "fk_service_id" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_order_id" FOREIGN KEY ("order_id") REFERENCES "orders" ("id")
)WITHOUT OIDS;