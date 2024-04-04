-- +goose Up
-- modify "todos" table
ALTER TABLE "public"."todos" ALTER COLUMN "deleted_at" SET NOT NULL;

-- +goose Down
-- reverse: modify "todos" table
ALTER TABLE "public"."todos" ALTER COLUMN "deleted_at" DROP NOT NULL;
