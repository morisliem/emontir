CREATE TABLE IF NOT EXISTS "services" (
  "id" SERIAL NOT NULL,
  "title" VARCHAR(128) UNIQUE NOT NULL,
  "description" VARCHAR(256),
  "rating" FLOAT NOT NULL,
  "price" FLOAT NOT NULL,
  "picture" VARCHAR(128),
  PRIMARY KEY ("id")
) WITHOUT OIDS;

CREATE INDEX IF NOT EXISTS "services_rating" ON "services" ("rating");
CREATE INDEX IF NOT EXISTS "services_price" ON "services" ("price");
CREATE INDEX IF NOT EXISTS trgm_idx on "services" using gin("title" gin_trgm_ops);

INSERT INTO "services" ("title", "description", "rating", "price", "picture") 
VALUES 
  ('pembersihan kaburator', 'pembersih karborator terbaik in town', '4.3', '250000', 'a.jpeg'),
  ('ganti kampas rem', NULL, '4.2', '90000', 'https://media.suara.com/pictures/970x544/2019/04/01/61779-rem-motor.jpg'),
  ('ganti oli', NULL, '4.5', '130000', 'a.jpeg'),
  ('ganti lampu depan', 'Ganti lampu depan termurah, tercepat dan terbaik se Indonesia', '4.1', '263500', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQT9yqkjEGGDUt9pzKLsnaHzcZk0xJdyCm1XQ&usqp=CAU'),
  ('ganti lampu sein', 'selesai dalam 5 menit', '4.3', '213000', '890023500.jpg'),
  ('ganti lampu rem', NULL, '4.3', '330500', 'a.jpeg'),
  ('ganti ban tubeless', 'Twee shabby chic taiyaki flannel, enamel pin venmo vape four loko. Hexagon kale chips typewriter kitsch 8-bit organic plaid small batch keffiyeh ethical banh mi narwhal echo park cronut.', '4.8', '420000', 'a.jpeg'),
  ('ganti lampu belakang',NULL,'5','250000','a.jpeg'),
  ('ganti knalpot',NULL,'4.4','525000','a.jpeg'),
  ('service berkala',NULL,'4.8','700000','a.jpeg'),
  ('service jangka panjang',NULL,'4.6','500000','a.jpeg'),
  ('ganti velg',NULL,'4','800000','a.jpeg'),
  ('cuci motor',NULL,'4.6','120000','a.jpeg'),
  ('tambal ban',NULL,'4.2','20000','a.jpeg'),
  ('pompa ban',NULL,'4.4','20000','a.jpeg')
ON CONFLICT DO NOTHING;