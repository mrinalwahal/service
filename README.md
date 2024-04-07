# Microservice Boilerplate

### Connection Pooling
Only open a connection to your database w/ `gorm.Open()` just once in your code and pass it everywhere from a global variable. Gorm's underlying `sql.DB()` interface will automatically use connection pooling for every transaction. To use connection pooling, it is important to configure the following values:

```
sqlDB.SetConnMaxIdleTime(time.Minute * 5)
sqlDB.SetConnMaxLifetime(time.Minute * 5)
sqlDB.SetMaxOpenConns(10)
sqlDB.SetMaxIdleConns(0)
```