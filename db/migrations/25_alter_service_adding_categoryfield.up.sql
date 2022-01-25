ALTER TABLE "services" 
    ADD COLUMN "category" VARCHAR(128);

UPDATE "services" SET "category" = 'kaburator' WHERE "id" = 1;
UPDATE "services" SET "category" = 'rem, popular' WHERE "id" = 2;
UPDATE "services" SET "category" = 'oli' WHERE "id" = 3;
UPDATE "services" SET "category" = 'lampu, popular' WHERE "id" = 4;
UPDATE "services" SET "category" = 'lampu' WHERE "id" = 5;
UPDATE "services" SET "category" = 'lampu' WHERE "id" = 6;
UPDATE "services" SET "category" = 'ban' WHERE "id" = 7;
UPDATE "services" SET "category" = 'lampu' WHERE "id" = 8;
UPDATE "services" SET "category" = 'knalpot' WHERE "id" = 9;
UPDATE "services" SET "category" = 'monthly_services' WHERE "id" = 10;
UPDATE "services" SET "category" = 'monthly_services, popular' WHERE "id" = 11;
UPDATE "services" SET "category" = 'velg' WHERE "id" = 12;
UPDATE "services" SET "category" = 'cleaning, popular' WHERE "id" = 13;
UPDATE "services" SET "category" = 'ban' WHERE "id" = 14;
UPDATE "services" SET "category" = 'ban, popular' WHERE "id" = 15;
