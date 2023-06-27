# Go SQL builder

## Install

```bash
go get -u github.com/LineoIT/sqlb
```

## Basic usages



#### Select query


```go
package main

import (
	"github.com/LineoIT/sqlb"
)

func main() {

	// Example 1
	q := *sqlb.Select("*").
		From("users").
		Where(sqlb.In("age", []string{"76", "80"})).
		Or("email", "aaa@ajks.com").
		Where(sqlb.Between("salary", []float64{5000, 5900})).
		Or(sqlb.NotIn("role", []string{"admin", "driver"})).
		GroupBy("id", "age").
		Having(sqlb.NotEq("email", "aaa@ajks.com")).
		Or(sqlb.NotIn("item", []int{0, 1})).
		Limit(90).
		Offset(7).
		OrderBy("id").
		Build()
}
```


```bash
# output
select * from users where age in (76,80) or email aaa@ajks.com and salary between 5000 and 5900 or role not in (admin,driver) group by id,age having email <> aaa@ajks.com or item not in (0,760) order by id limit 90 offset 7;
```

```go
// Example 2
q := *sqlb.SQL("select * from users").
	Where(sqlb.In("age", []string{"76", "80"})).
	Or("email", "aaa@ajks.com").
	Where(sqlb.Between("salary", []float64{5000, 5900})).
	Or(sqlb.NotIn("role", []string{"admin", "driver"})).
	GroupBy("id", "age").
	Having(sqlb.NotEq("email", "aaa@ajks.com")).
	Or(sqlb.NotIn("item", []int{0, 1})).
	Limit(90).
	Offset(7).
	OrderBy("id").
	Build()

```

```bash
# output
select * from users where age in (76,80) or email aaa@ajks.com and salary between 5000 and 5900 or role not in (admin,driver) group by id,age having email <> aaa@ajks.com or item not in (0,760) order by id limit 90 offset 7;
```


### Insertion

```go
q := sqlb.Insert("users").Value("email", "mail@example.com").
	Value("age", 10).
	Return("id").
	Build()
```

```bash
# output
insert into users (email,age) values(mail@example.com,10) returning id
```

### Update

```go
q := sqlb.Update("users").Set("email", "mail@example.com").
	Set("age", 10).
	Where("id", 2).
	Or(sqlb.NotIn("item", []int{0, 1})).
	Return("updated_at").
	Build()
```

```bash
# output
update users set email=mail@example.com,age=10 where id=2 or item not in (0,1) returning updated_at;
```


### Delete

```go
q := sqlb.Delete("users").
	Where("id", 2).
	Or(sqlb.NotIn("item", []int{0, 1})).
	Build()
```

```bash
# output
delete from users where id=2 or item not in (0,1);
```


### Base Functions

* Select
* From
* SQL
* Build
* Stmt
* Error
* Args



### Clause Functions

* Where
* Or
* Having
* GroupBy
* OrderBy
* Offset
* Limit
* Sort
* Take
* Raw

### Expression functions

* Coalesce
* Nullif
* Equal
* NotEqual
* In
* NotIn
* Between
* IsNull
* IsNotNull
* Is
* IsNot
* Like
* Ilike
* Expression
* Greater
* GreaterOrEqual
* Less



### Helpers

* Debug
* CleanSQL