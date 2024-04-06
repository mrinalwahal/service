# Database layer

- To compare the models with database schema, add the models to `/scripts/loader.go`. This will help generate the migrations for your models.
- To generate a new migrations, run `/scripts/migrate.sh [name_of_your_migration]`. For example: `/scripts/migrate.sh init` will generate a new migration in `/migrations` directory.
- To compare the state/status of the migrations against the database schema, run `/scripts/status.sh`. This will print which migrations are pending to be applied and which ones have been applied.
- To apply any pending migrations, run `/scripts/apply.sh`.
