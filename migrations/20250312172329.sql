-- Set comment to column: "code" on table: "countries"
COMMENT ON COLUMN "countries"."code" IS 'ISO 3166-1 alpha-2 code';
-- Set comment to column: "runtime" on table: "movies"
COMMENT ON COLUMN "movies"."runtime" IS 'Duration in minutes';
-- Modify "sessions" table
ALTER TABLE "sessions" DROP COLUMN "token", ADD COLUMN "access_token" text NOT NULL, ADD COLUMN "refresh_token" text NOT NULL, ADD COLUMN "refresh_expiry" timestamptz NOT NULL;
-- Create index "idx_sessions_access_token" to table: "sessions"
CREATE UNIQUE INDEX "idx_sessions_access_token" ON "sessions" ("access_token");
-- Create index "idx_sessions_refresh_token" to table: "sessions"
CREATE UNIQUE INDEX "idx_sessions_refresh_token" ON "sessions" ("refresh_token");
