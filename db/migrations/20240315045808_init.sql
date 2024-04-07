-- +goose Up
-- create "todos" table
CREATE TABLE "public"."todos" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "title" text NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_todos_deleted_at" to table: "todos"
CREATE INDEX "idx_todos_deleted_at" ON "public"."todos" ("deleted_at");

-- +goose Down
-- reverse: create index "idx_todos_deleted_at" to table: "todos"
DROP INDEX "public"."idx_todos_deleted_at";
-- reverse: create "todos" table
DROP TABLE "public"."todos";
