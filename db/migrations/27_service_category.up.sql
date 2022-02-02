CREATE TABLE IF NOT EXISTS "service_categories"(
    "service_id" INT NOT NULL,
    "category" VARCHAR(128) NOT NULL,
    PRIMARY KEY ("service_id", "category"),
    CONSTRAINT "fk_service_id" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_category" FOREIGN KEY ("category") REFERENCES "categories" ("category")
);