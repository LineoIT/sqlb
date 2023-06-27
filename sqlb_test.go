package sqlb

import (
	"testing"
)

func TestFilterQuery(t *testing.T) {
	q := Select("*").
		From("users").
		Where(In("age", []string{"76", "80"})).
		Or("email", "aaa@ajks.com").
		Where(Between("salary", []float64{5000, 5900})).
		Or(NotIn("role", []string{"admin", "driver"})).
		GroupBy("id", "age").
		Having(NotEq("email", "aaa@ajks.com")).
		Or(NotIn("item", []int{0, 1})).
		Limit(90).
		Offset(7).
		OrderBy("id").
		Build()
	result := "select * from users where age in (76,80) or email=aaa@ajks.com and salary between 5000 and 5900 or role not in (admin,driver) group by id,age having email <> aaa@ajks.com or item not in (0,760) order by id limit 90 offset 7;"
	if result != q.Debug() {
		t.Errorf("TestFilterQuery:\texpected: %v, \n\tgot:%v", result, q.Debug())
	}
}

func TestInsertQuery(t *testing.T) {
	q := Insert("users").Value("email", "mail@example.com").
		Value("age", 10).
		Return("id").
		Build()
	result := "insert into users (email,age) values(mail@example.com,10) returning id"
	if result != q.Debug() {
		t.Errorf("TestInsertQuery:\texpected: %v, \n\tgot:%v", result, q.Debug())
	}
}

func TestUpdateQuery(t *testing.T) {
	q := Update("users").Set("email", "mail@example.com").
		Set("age", 10).
		Where("id", 2).
		Or(NotIn("item", []int{0, 1})).
		Return("updated_at").
		Build()
	result := "update users set email=mail@example.com,age=10 where id=2 or item not in (0,1) returning updated_at;"
	if result != q.Debug() {
		t.Errorf("TestUpdateQuery:\texpected: %v, \n\tgot:%v", result, q.Debug())
	}
}

func TestDeleteQuery(t *testing.T) {
	q := Delete("users").
		Where("id", 2).
		Or(NotIn("item", []int{0, 1})).
		Build()
	result := "delete from users where id=2 or item not in (0,1)"
	if result != q.Debug() {
		t.Errorf("TestDeleteQuery:\texpected: %v, \n\tgot:%v", result, q.Debug())
	}
}
