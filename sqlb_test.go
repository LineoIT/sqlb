package sqlb

import (
	"testing"
)

func TestFilterQuery(t *testing.T) {
	q := FilterQuery{
		Stmt:    "select * from users",
		OrderBy: "id",
		Limit:   90,
		Offset:  7,
	}
	q.Where("id", "=", 1).
		Or("email", "=", "aaa@ajks.com").
		Where("age", "in", 30, 67, "80080").
		Or("role", "in", "admin", "driver").
		GroupBy("id", "age").
		Having("email", "=", "aaa@ajks.com").
		Or("item", "in", 0, 1).
		Build()
	result := "select * from users where id = 1 or email = aaa@ajks.com and age in (30,67,80080) or role in (admin,driver) group by id,age having email = aaa@ajks.com or item in (0,10) order by id limit 90 offset 7;"
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
		Where("id", "=", 2).
		Return("updated_at").
		Build()
	result := "update users set email=mail@example.com,age=10 where id = 2 returning updated_at"
	if result != q.Debug() {
		t.Errorf("TestInsertQuery:\texpected: %v, \n\tgot:%v", result, q.Debug())
	}
}
