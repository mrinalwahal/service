# Database layer

- To compare the models with database schema, add the models to `/scripts/loader.go`. This will help generate the migrations for your models.
- To generate a new migrations, run `/scripts/migrate.sh [name_of_your_migration]`. For example: `/scripts/migrate.sh init` will generate a new migration in `/migrations` directory.
- To compare the state/status of the migrations against the database schema, run `/scripts/status.sh`. This will print which migrations are pending to be applied and which ones have been applied.
- To apply any pending migrations, run `/scripts/apply.sh`.

**To generate mock files of the service interface, use the following commands:**

1. Install mockgen.
    ```
    go install go.uber.org/mock/mockgen@latest
    ```
1. Generate mocks.
    ```
    mockgen -destination=mock.go -source=service.go -package=record
    ```

This will generate the file `mock.go` which will contains your mock service. You can import it in your tests with:

```
func TestFoo(t *testing.T) {
  ctrl := gomock.NewController(t)

  m := NewMockService(ctrl)

  // Asserts that the first and only call to Bar() is passed 99.
  // Anything else will fail.
  m.
    EXPECT().
    Bar(gomock.Eq(99)).
    Return(101)

  YourUnitTest(m)
}
```