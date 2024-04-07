-- +goose Up
-- create "records" table
CREATE TABLE "public"."records" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "title" text NOT NULL,
  PRIMARY KEY ("id")
);

-- +goose Down
-- reverse: create "records" table
DROP TABLE "public"."records";
