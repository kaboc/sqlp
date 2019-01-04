# sqlp

sqlp is a Go package extending database/sql to make it a little easier to use by adding some features that may come in handy for you.
The key features are:

* Bulk inserting data in structs into a table
* Reading rows into structs, maps or slices
* Easier binding of values to unnamed/named placeholders

### Table of Contents

* [Installation](#installation)
* [Usage](#usage)
    * [Example Table](#example-table)
    * [Getting Started](#getting-started)
    * [Insert](#insert)
        * [Note](#note)
    * [Scan](#scan)
        * [Into struct](#into-struct)
        * [Into map](#into-map)
        * [Into slice](#into-slice)
    * [Select](#select)
        * [Into slice of structs](#into-slice-of-structs)
        * [Into slice of maps](#into-slice-of-maps)
        * [Into slice of slices](#into-slice-of-slices)
* [Placeholders](#placeholders)
    * [Unnamed Placeholder](#unnamed-placeholder)
        * [For different types of placeholder](#for-different-types-of-placeholder)
    * [Named Placeholder](#named-placeholder)
* [License](#license)

## Installation

```
go get github.com/kaboc/sqlp
```

## Usage

Basic usage is mostly the same as for database/sql, so only major differences are described in this section.

### Example Table

MySQL

```sql
CREATE TABLE user (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(32) NOT NULL,
  age tinyint(3) unsigned NOT NULL,
  recorded_at datetime DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

PostgreSQL

```sql
create table "user" (
  id serial not null primary key,
  name varchar(32) not null,
  age smallint not null,
  recorded_at timestamp
);
```

### Getting Started

Import `sqlp` and other necessary packages including a database driver like `go-sql-driver/mysql`.

```go
import (
    "github.com/go-sql-driver/mysql"
    "github.com/kaboc/sqlp"
)
```

```go
db, err := sqlp.Open("mysql", "user:pw@tcp(host:3306)/dbname")
```

Use `sqlp.Init()` instead if there is a connection already opened by database/sql's `Open()`.

```go
sqlDB, err := sql.Open("mysql", "user:pw@tcp(host:3306)/dbname")
db := sqlp.Init(sqlDB)
```

### Insert

You can bulk insert multiple rows easily using a slice of structs containing sets of data to be inserted.

```go
type tUser struct {
    Name       string
    Age        int
    RecordedAt mysql.NullTime `col:"recorded_at"`
}

now := mysql.NullTime{Time: time.Now(), Valid: true}
data := []tUser{
    {Name: "User1", Age: 22, RecordedAt: now},
    {Name: "User2", Age: 27, RecordedAt: mysql.NullTime{}},
    {Name: "User3", Age: 31, RecordedAt: now},
}

res, err := db.Insert("user", data)
if err != nil {
    log.Fatal(err)
}

cnt, _ := res.RowsAffected()
fmt.Printf("%d rows were affected", cnt) // 3 rows were affected
```

Struct fields need to be capitalized so that sqlp can access them.

A tag is necessary only when the field name is not the same as the column name.
In the above example, `col:"age"` can be omitted since columns are case insensitive in MySQL by default and `Age` and `age` are not distinguished.

Values are processed via [placeholders](#placeholders) internally and escaped to be safe. There is no need to worry about SQL injection.

#### Note

* If the table name is a reserved keyword, it has to be enclosed with back quotes, double quotes, etc. depending on the DBMS. Below is an example for PostgreSQL.

    ```go
    res, err := db.Insert(`"user"`, data)
    ```

### Scan

#### Into struct

```go
type tUser struct {
    Name       string
    Age        int
    RecordedAt mysql.NullTime
}

rows, err := db.Query(`SELECT name, age, recorded_at FROM user`)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var u tUser
    err = rows.ScanToStruct(&u)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s: %d yo [%s, %t]\n", u.Name, u.Age, u.RecordedAt.Time, u.RecordedAt.Valid)
}

// User1: 22 yo [2018-06-24 01:23:45 +0000 UTC, true]
// User2: 27 yo [0001-01-01 00:00:00 +0000 UTC, false]
// User3: 31 yo [2018-06-24 01:23:45 +0000 UTC, true]
```

Columns are mapped to corresponding struct fields.

Here, unlike in the above Insert example, the `RecordedAt` field does not have the `` `col:"recorded_at` `` tag.
This is because `RecordedAt` is regarded as identical to `recorded_at` by case-insensitive comparison after underscores are removed.

#### Into map

```go
for rows.Next() {
    u, err := rows.ScanToMap()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s: %s yo [%s]\n", u["name"], u["age"], u["recorded_at"])
}

// User1: 22 yo [2018-06-24T01:23:45+00:00]
// User2: 27 yo []
// User3: 31 yo [2018-06-24T01:23:45+00:00]
```

#### Into slice

```go
for rows.Next() {
    u, err := rows.ScanToSlice()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s: %s yo, [%s]\n", u[0], u[1], u[2])
}

// User1: 22 yo, [2018-06-24T01:23:45+00:00]
// User2: 27 yo, []
// User3: 31 yo, [2018-06-24T01:23:45+00:00]
```

### Select

#### Into slice of structs

```go
type tUser struct {
    Name       string
    Age        int
    RecordedAt mysql.NullTime
}

var u []tUser
err := db.SelectToStruct(&u, `SELECT name, age, recorded_at FROM user`)
fmt.Println(u)
```

This saves you from making a query and then scanning each row.
It is convenient, but be careful not to use up huge amounts of memory by fetching too many rows into a slice at a time.

#### Into slice of maps

```go
u, err := db.SelectToMap(`SELECT name, age, recorded_at FROM user`)
fmt.Println(u)
```

#### Into slice of slices

```go
u, err := db.SelectToSlice(`SELECT name, age, recorded_at FROM user`)
fmt.Println(u)
```

## Placeholders

sqlp provides both named and unnamed placeholders.

### Unnamed Placeholder

This is quite similar to database/sql's placeholder, with only several differences:

* Only `?` is used regardless of the type of DBMS or the database driver. `$1` or other types are not available.
* `WHERE name IN (?, ?)` can be replaced with `WHERE name IN ?[2]`.
* Binding values are passed as literals, variables, slices, or combinations of these.

Example:

```go
q := `SELECT name, age, recorded_at FROM user
      WHERE name LIKE ? AND age IN ?[2]`
```

This is internally converted to the next statement:

```sql
SELECT name, age, recorded_at FROM user
WHERE name LIKE ? AND age IN (?,?)
```

The following three ways of binding values are all acceptable.

```go
u, err := db.SelectToMap(q, "User%", 22, 31)
```

```go
b1 := "User%"
b2 := []interface{}{22, 31}
u, err := db.SelectToMap(q, b1, b2)
```

```go
b := []interface{}{"User%", 22, 31}
u, err := db.SelectToMap(q, b)
//u, err := db.SelectToMap(q, b...) // This works fine too.
```

#### For different types of placeholder

If the DBMS or the database driver that you use is not compatible with the `?` type of placeholder, you will need to instruct sqlp to use another one.

For example, PostgreSQL uses `$1` instead of `?`.

```sql
SELECT name, age, recorded_at FROM user
WHERE name LIKE $1 AND age IN ($2,$3)
```

This type of placeholder is defined as the constant `placeholder.Dollar`.
If you specify it by `placeholder.SetType()`, sqlp converts `?` to `$1` internally, so you can use `?` in your query.

```go
placeholder.SetType(placeholder.Dollar)
q := "SELECT * FROM user WHERE name LIKE ? AND age IN ?[2]"
u, err := db.SelectToMap(q, "User%", 22, 31)
```

Another way is to define a conversion function on your own.
You should be able to make do with this even if definition of your required type is missing in sqlp.

```go
placeholder.SetConvertFunc(func(query *string) {
    cnt := strings.Count(*query, "?")
    for i := 1; i <= cnt; i++ {
        *query = strings.Replace(*query, "?", "$"+strconv.Itoa(i), 1)
    }
})
```

### Named Placeholder

This is radically different from database/sql's named placeholder.
Here is an example similar to the previous one.

```go
q := `SELECT name, age, recorded_at FROM user
      WHERE name LIKE :like AND age IN :age[2]`
```

`:like` and `:age[2]` are the named placeholders.
They are internally converted to unnamed ones as below:

```sql
SELECT name, age, recorded_at FROM user
WHERE name LIKE ? AND age IN (?,?)
```
Values are passed only in the form of a single map.
`:XXXX[N]` requires a slice of interface{} with N numbers of elements.

```go
b := map[string]interface{}{
    "like": "User%",
    "age":  []interface{}{22, 31},
}
u, err := db.SelectToMap(q, b)
```

It is the same here as for unnamed placeholder that `placeholer.SetType()` or `placeholder.SetConvertFunc()` is necessary if your DBMS or database driver does not support the `?` type.

## License

[MIT](./LICENSE)
