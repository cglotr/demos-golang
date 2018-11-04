# crud-mysql

## Development

1. Create `goblog` database

```
CREATE DATABASE goblog;
```

2. Create MySQL `test` user with password `T3$t!992`

```
CREATE USER 'test'@'localhost' IDENTIFIED BY 'T3$t!992';
```

3. Grant privileges to `test` user

```
GRANT ALL PRIVILEGES ON goblog.* TO 'test'@'localhost' IDENTIFIED BY 'T3$t!992';
```

4. Create `employee` table

```
CREATE TABLE employee (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR (50) NOT NULL,
    city VARCHAR (50),
    PRIMARY KEY (id)
);
```

4. Run the server with `go run main.go`
