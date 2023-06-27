package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/LineoIT/sqlb"
)

func DoesNotContainSpecialCharacters(str string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_]*[a-zA-Z0-9][<!=>]?$`) // Combined pattern

	return regex.MatchString(str)
}

func main() {
	q := sqlb.Update("users").Set("email", "mail@example.com").
		Set("age", 10).
		Where("id", 2).
		Return("updated_at").
		Or(sqlb.NotIn("role", []string{"admin", "driver"})).
		GroupBy("id", "age").
		Having(sqlb.NotEq("email", "aaa@ajks.com")).
		Or(sqlb.NotIn("item", []int{0, 1})).
		Limit(90).
		Offset(7).
		OrderBy("id").
		Build()

	fmt.Println(q.Stmt())
	fmt.Println()
	fmt.Println(q.Debug())
}
func selectExample() {

	builder := sqlb.Select("*").
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

	fmt.Println(builder.Stmt())
	fmt.Println()
	fmt.Println(builder.Debug())

}

func selectExample2() {

	builder := *sqlb.SQL("select * from users").
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

	fmt.Println(builder.Stmt())
	fmt.Println()
	fmt.Println(builder.Debug())

}

func selectExample3() {

	builder := *sqlb.SQL("select * from users").
		Where("age in", []string{"76", "80"}).
		Or("email", "aaa@ajks.com").
		Where("phone", "sdfsj").
		Where(sqlb.Like("age", "sdfsj")).
		Or(sqlb.Ilike("age", "'%acc%'")).
		Where(sqlb.Expression("age", "<>", 20)).
		Build()

	if err := builder.Error(); err != nil {
		panic(err)
	}

	fmt.Println(builder.Stmt())
	fmt.Println()
	fmt.Println(builder.Debug())

}

func insertUpdateExamples() {
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
		Where("id", user.ID).
		Return("updated_at").
		Build()
	fmt.Println(q1.Debug())
}
