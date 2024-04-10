-- +goose Up
-- create "records" table
CREATE TABLE "public"."records" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "title" text NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_records_title" CHECK (length(title) > 0)
);

-- +goose Down
-- reverse: create "records" table
DROP TABLE "public"."records";
