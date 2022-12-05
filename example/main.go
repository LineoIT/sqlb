package main

import (
	"fmt"
	"time"

	"github.com/LineoIT/sqlb"
)

func main() {
	builder := sqlb.FilterQuery{
		Stmt:    "select * from users",
		OrderBy: "ID",
		Limit:   90,
		Offset:  7,
	}
	builder.Where("id", "=", 1).
		Or("email", "=", "aaa@ajks.com").
		Where("age", "in", 30, 67, "80080").
		Or("role", "in", "admin", "driver").
		GroupBy("id", "age").
		Having("email", "=", "aaa@ajks.com").
		Or("item", "in", 0, 1).
		Build()
	fmt.Println(builder.Debug())

	// Insert and update
	type User struct {
		ID        int
		Email     string
		Phone     string
		Role      string
		Age       int
		UpdatedAt time.Time
	}

	user := User{
		Email: "abc@example.com",
		Phone: "+7939739473",
		Age:   29,
		ID:    3,
	}

	// insert
	q := sqlb.Insert("users").
		Value("email", user.Email).
		Value("phone", sqlb.Nullif(user.Phone, "''")).
		Value("age", sqlb.Coalesce(user.Age, "age")).
		Return("id").
		Build()
	fmt.Println(q.Debug())

	// update
	q1 := sqlb.Update("users").
		Set("email", user.Email).
		Set("phone", sqlb.Nullif(user.Phone, "''")).
		Set("age", sqlb.Coalesce(user.Age, "age")).
		Set("role", sqlb.Coalesce(sqlb.Nullif(user.Role, "''"), "role")).
		Where("id", "=", user.ID).
		Return("updated_at").
		Build()
	fmt.Println(q1.Debug())
}
