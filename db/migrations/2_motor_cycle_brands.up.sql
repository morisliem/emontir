CREATE TABLE IF NOT EXISTS "motor_cycle_brands" (
  "name" VARCHAR(128),
  PRIMARY KEY("name")
) WITHOUT OIDS;

INSERT INTO "motor_cycle_brands" ("name")
VALUES 
  ('honda beat'),
  ('honda scoopy'),
  ('honda pcx'),
  ('yamaha mio'),
  ('yamaha nmax'),
  ('suzuki')
ON CONFLICT ("name") DO NOTHING;