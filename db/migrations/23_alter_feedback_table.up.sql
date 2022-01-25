ALTER TABLE "feedbacks" 
    ADD COLUMN "is_reviewed" BOOLEAN,
    ADD CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;