-- Modify "users" table
ALTER TABLE "public"."users" ADD CONSTRAINT "chk_users_email" CHECK (length(email) > 0);
