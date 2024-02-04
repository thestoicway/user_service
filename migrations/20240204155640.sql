-- Create "profiles" table
CREATE TABLE "public"."profiles" (
  "user_id" uuid NULL,
  "name" text NULL,
  "birth_date" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL
);
-- Create index "idx_profiles_deleted_at" to table: "profiles"
CREATE INDEX "idx_profiles_deleted_at" ON "public"."profiles" ("deleted_at");
