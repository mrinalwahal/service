-- +goose Up
-- drop index "idx_todos_deleted_at" from table: "todos"
DROP INDEX "public"."idx_todos_deleted_at";
-- modify "todos" table
ALTER TABLE "public"."todos" ALTER COLUMN "deleted_at" DROP NOT NULL;

-- +goose Down
-- reverse: modify "todos" table
ALTER TABLE "public"."todos" ALTER COLUMN "deleted_at" SET NOT NULL;
-- reverse: drop index "idx_todos_deleted_at" from table: "todos"
CREATE INDEX "idx_todos_deleted_at" ON "public"."todos" ("deleted_at");
