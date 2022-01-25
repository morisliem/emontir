CREATE TABLE IF NOT EXISTS "user_fav_services"(
    "user_id" UUID NOT NULL,
    "service_id" INT NOT NULL, 
    PRIMARY KEY("user_id", "service_id"),
    CONSTRAINT "fk_service_id" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id")
)WITHOUT OIDS;