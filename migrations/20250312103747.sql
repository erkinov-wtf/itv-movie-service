-- Create "countries" table
CREATE TABLE "countries" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "code" character varying(2) NOT NULL,
  "continent" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_countries_code" to table: "countries"
CREATE UNIQUE INDEX "idx_countries_code" ON "countries" ("code");
-- Create index "idx_countries_name" to table: "countries"
CREATE UNIQUE INDEX "idx_countries_name" ON "countries" ("name");
-- Set comment to column: "code" on table: "countries"
COMMENT ON COLUMN "countries"."code" IS 'ISO 3166-1 alpha-2 code';
-- Create "languages" table
CREATE TABLE "languages" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "code" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_languages_code" to table: "languages"
CREATE UNIQUE INDEX "idx_languages_code" ON "languages" ("code");
-- Create index "idx_languages_deleted_at" to table: "languages"
CREATE INDEX "idx_languages_deleted_at" ON "languages" ("deleted_at");
-- Create index "idx_languages_name" to table: "languages"
CREATE INDEX "idx_languages_name" ON "languages" ("name");
-- Create "movies" table
CREATE TABLE "movies" (
  "id" uuid NOT NULL,
  "title" text NOT NULL,
  "director" text NULL,
  "year" integer NULL,
  "plot" text NULL,
  "runtime" integer NULL,
  "rating" numeric(3,1) NULL DEFAULT 0,
  "poster_url" text NULL,
  "trailer_url" text NULL,
  "release_date" date NULL,
  "language" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_movies_language" FOREIGN KEY ("language") REFERENCES "languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_movies_director" to table: "movies"
CREATE INDEX "idx_movies_director" ON "movies" ("director");
-- Create index "idx_movies_title" to table: "movies"
CREATE INDEX "idx_movies_title" ON "movies" ("title");
-- Create index "idx_movies_year" to table: "movies"
CREATE INDEX "idx_movies_year" ON "movies" ("year");
-- Set comment to column: "runtime" on table: "movies"
COMMENT ON COLUMN "movies"."runtime" IS 'Duration in minutes';
-- Create "movie_countries" table
CREATE TABLE "movie_countries" (
  "movie_id" uuid NOT NULL,
  "country_id" uuid NOT NULL,
  PRIMARY KEY ("movie_id", "country_id"),
  CONSTRAINT "fk_movie_countries_country" FOREIGN KEY ("country_id") REFERENCES "countries" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_movie_countries_movie" FOREIGN KEY ("movie_id") REFERENCES "movies" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "genres" table
CREATE TABLE "genres" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_genres_name" to table: "genres"
CREATE UNIQUE INDEX "idx_genres_name" ON "genres" ("name");
-- Create "movie_genres" table
CREATE TABLE "movie_genres" (
  "movie_id" uuid NOT NULL,
  "genre_id" uuid NOT NULL,
  PRIMARY KEY ("movie_id", "genre_id"),
  CONSTRAINT "fk_movie_genres_genre" FOREIGN KEY ("genre_id") REFERENCES "genres" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_movie_genres_movie" FOREIGN KEY ("movie_id") REFERENCES "movies" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "first_name" text NOT NULL,
  "last_name" text NOT NULL,
  "username" text NOT NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "role" text NULL DEFAULT 'user',
  "active" boolean NULL DEFAULT true,
  "last_login_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username");
-- Create "sessions" table
CREATE TABLE "sessions" (
  "id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "token" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "user_agent" text NULL,
  "ip_address" text NULL,
  "is_revoked" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_sessions" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_sessions_token" to table: "sessions"
CREATE UNIQUE INDEX "idx_sessions_token" ON "sessions" ("token");
-- Create index "idx_sessions_user_id" to table: "sessions"
CREATE INDEX "idx_sessions_user_id" ON "sessions" ("user_id");
