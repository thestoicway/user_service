-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "email" SET NOT NULL, ALTER COLUMN "password_hash" SET NOT NULL;
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "public"."users" ("email");
