# Service Layer

This layer contains the fundamental business logic of this service.

## Testing

### Unit / Whitebox Tests

All the essential unit tests to be covered:

- [x] Create a new record.
- [x] List all the records with supported filters.
- [x] Get a record by it's ID.
- [x] Update a record with new options.
- [x] Delete a record.

### Integration / Blackbox Tests

To write itnegration tests, you would typically want to mock this layer's interfaces and consume them outside the package.

**To generate mock files of the service interface, use the following commands:**

1. Install mockgen.
    ```
    go install go.uber.org/mock/mockgen@latest
    ```
1. Generate mocks.
    ```
    mockgen -destination=mock.go -source=service.go -package=service
    ```

This will generate the file `mock.go` which will contain your mock service. You can import it in your tests with:

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

  SUT(m)
}
```