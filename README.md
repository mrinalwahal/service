# Records Service Boilerplate

## Design

- [Google Cloud Design Guide](https://cloud.google.com/apis/design).
- Naming conventions by [Google](https://cloud.google.com/apis/design/naming_convention#product_names).

### Connection Pooling
Only open a connection to your database w/ `gorm.Open()` just once in your code and pass it everywhere from a global variable. Gorm's underlying `sql.DB()` interface will automatically use connection pooling for every transaction. To use connection pooling, it is important to configure the following values:

```
sqlDB.SetConnMaxIdleTime(time.Minute * 5)
sqlDB.SetConnMaxLifetime(time.Minute * 5)
sqlDB.SetMaxOpenConns(10)
sqlDB.SetMaxIdleConns(0)
```

### Loggnig Do's and Don'ts

- Establish clear logging objectives
- Use log levels correctly
- Structure your logs
- Write meaningful log entries
- Sample your logs
- Use canonical log lines
- Aggregate and centralize your logs
- Establish log retention policies
- Protect your logs
- Don't log sensitive data
- Don't ignore the performance cost of logging
- Don't use logs for monitoring