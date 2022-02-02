CREATE TABLE IF NOT EXISTS "users" (
  "id" UUID NOT NULL,
  "name" VARCHAR(128) NOT NULL,
  "email" VARCHAR(128) UNIQUE NOT NULL,
  "password" VARCHAR(60) NOT NULL,
  "address" VARCHAR(256),
  "phone_num" VARCHAR(15),
  "is_active" BOOLEAN NOT NULL,
  PRIMARY KEY ("id")
);