CREATE TABLE IF NOT EXISTS "mechanics"(
    "id" SERIAL NOT NULL,
    "name" VARCHAR(128) NOT NULL,
    "phone_number" VARCHAR(16) NOT NULL,
    "is_available" BOOLEAN NOT NULL,
    "completed_service" INT NOT NULL,
    "picture" VARCHAR(256),
    PRIMARY KEY("id")
)WITHOUT OIDS;

INSERT INTO "mechanics" ("name", "phone_number", "is_available", "completed_service", "picture")
VALUES
  ('pratama', '+6281111111110', FALSE, '10', 'https://www.floridacareercollege.edu/wp-content/uploads/sites/4/2020/06/3-Reasons-Why-Being-a-Mechanic-Could-Be-An-Amazing-Career-Florida-Career-College.jpeg'),
  ('nendi', '+6281111111111', TRUE, '1','https://thrivesjournal.com/wp-content/uploads/2021/06/Ask-A-Mechanic-1200x900-1.jpeg'),
  ('adi', '+6281111111112', TRUE, '19','https://www.autoguru.com.au/Content/images/directories/mechanics.jpg'),
  ('dimas', '+6281111111113', TRUE, '3',''),
  ('andi', '+6281111111114', FALSE, '15','https://www.oneeducation.org.uk/wp-content/uploads/2017/09/Car-Mechanic-Training.jpg'),
  ('budi', '+6281111111115', TRUE, '4',''),
  ('jono', '+6281111111116', TRUE, '9','https://www.opencolleges.edu.au/static/oca/media/patchwork/pages/Mechanic_HERO_B.jpg'),
  ('aldo', '+6281111111117', TRUE, '18','https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/2021-popularmechanics-guiseppesgarage-series-ep2-oil-change-web-00-01-23-19-still013-1637180560.jpg'),
  ('tomi', '+6281111111118', TRUE, '22','https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT5AIpDx8hrHcpLyj4MgYUCtwbzIcNsXL8Lj3axqKT_C2SWdqmVPCssP-oE5wGeGeVRyBI&usqp=CAU'),
  ('tono', '+6281111111119', TRUE, '17','https://www.checkatrade.com/blog/wp-content/uploads/2021/07/Feature-mechanic-hourly-rate-uk.jpg')
ON CONFLICT DO NOTHING;