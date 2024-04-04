# XYZ Service

**To generate mock files of the service interface, use the following commands:**

1. Install mockgen.
    ```
    go install go.uber.org/mock/mockgen@latest
    ```
1. Generate mocks.
    ```
    mockgen -destination=mock.go -source=service.go -package=todo
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

  SUT(m)
}
```