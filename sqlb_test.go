package sqlb

import (
	"strings"
	"testing"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func check(t *testing.T, name, want, got string) {
	t.Helper()
	if want != got {
		t.Errorf("%s\n  want: %s\n   got: %s", name, want, got)
	}
}

func checkArgs(t *testing.T, name string, want, got []interface{}) {
	t.Helper()
	if len(want) != len(got) {
		t.Errorf("%s args length: want %d got %d", name, len(want), len(got))
		return
	}
	for i := range want {
		if want[i] != got[i] {
			t.Errorf("%s args[%d]: want %v got %v", name, i, want[i], got[i])
		}
	}
}

// ── SELECT ────────────────────────────────────────────────────────────────────

func TestSelectBasic(t *testing.T) {
	q := Select("id", "name").From("users").Build()
	check(t, "SelectBasic", "select id,name from users;", q.Stmt())
}

func TestSelectStar(t *testing.T) {
	q := Select().From("users").Build()
	check(t, "SelectStar", "select * from users;", q.Stmt())
}

func TestSelectDistinct(t *testing.T) {
	q := Select("email").Distinct().From("users").Build()
	check(t, "SelectDistinct", "select distinct email from users;", q.Stmt())
}

func TestSelectLimitOffset(t *testing.T) {
	q := Select("*").From("users").Limit(10).Offset(20).Build()
	check(t, "LimitOffset", "select * from users limit 10 offset 20;", q.Stmt())
}

func TestSelectTake(t *testing.T) {
	q := Select("*").From("users").Take(5, 15).Build()
	check(t, "Take", "select * from users limit 5 offset 15;", q.Stmt())
}

func TestSelectOrderBy(t *testing.T) {
	q := Select("*").From("users").OrderBy("created_at").Sort(DESC).Build()
	check(t, "OrderBy", "select * from users order by created_at DESC;", q.Stmt())
}

func TestOrderByRaw(t *testing.T) {
	q := Select("*").From("users").OrderByRaw("lower(name)").Build()
	check(t, "OrderByRaw", "select * from users order by lower(name);", q.Stmt())
}

// ── WHERE ─────────────────────────────────────────────────────────────────────

func TestWhereEqual(t *testing.T) {
	q := Select("*").From("users").Where("id", 1).Build()
	check(t, "WhereEqual stmt", "select * from users where id=$1;", q.Stmt())
	checkArgs(t, "WhereEqual args", []interface{}{1}, q.Args())
}

func TestWhereAnd(t *testing.T) {
	q := Select("*").From("users").Where("active", true).Where("role", "admin").Build()
	check(t, "WhereAnd", "select * from users where active=$1 and role=$2;", q.Stmt())
}

func TestWhereOr(t *testing.T) {
	q := Select("*").From("users").Where("active", true).Or("role", "admin").Build()
	check(t, "WhereOr", "select * from users where active=$1 or role=$2;", q.Stmt())
}

func TestWhereIn(t *testing.T) {
	q := Select("*").From("users").Where(In("id", []int{1, 2, 3})).Build()
	check(t, "WhereIn stmt", "select * from users where id in ($1,$2,$3);", q.Stmt())
	checkArgs(t, "WhereIn args", []interface{}{1, 2, 3}, q.Args())
}

func TestWhereNotIn(t *testing.T) {
	q := Select("*").From("users").Where(NotIn("role", []string{"admin", "root"})).Build()
	check(t, "WhereNotIn", "select * from users where role not in ($1,$2);", q.Stmt())
}

func TestWhereBetween(t *testing.T) {
	q := Select("*").From("orders").Where(Between("total", []float64{100, 500})).Build()
	check(t, "WhereBetween", "select * from orders where total between $1 and $2;", q.Stmt())
	checkArgs(t, "WhereBetween args", []interface{}{float64(100), float64(500)}, q.Args())
}

func TestWhereBetweenNoDoubleSpace(t *testing.T) {
	q := Select("*").From("users").
		Where("active", true).
		Where(Between("age", []int{18, 65})).
		Build()
	stmt := q.Stmt()
	if strings.Contains(stmt, "  ") {
		t.Errorf("WhereBetween: double space in stmt: %s", stmt)
	}
	check(t, "WhereBetweenNoDoubleSpace",
		"select * from users where active=$1 and age between $2 and $3;", stmt)
}

func TestWhereIsNull(t *testing.T) {
	q := Select("*").From("users").Where(IsNull("deleted_at")).Build()
	check(t, "WhereIsNull", "select * from users where deleted_at is null;", q.Stmt())
}

func TestWhereIsNotNull(t *testing.T) {
	q := Select("*").From("users").Where(IsNotNull("email")).Build()
	check(t, "WhereIsNotNull", "select * from users where email is not null;", q.Stmt())
}

func TestWhereLike(t *testing.T) {
	q := Select("*").From("users").Where(Like("name", "John%")).Build()
	check(t, "WhereLike", "select * from users where name like $1;", q.Stmt())
}

func TestWhereIlike(t *testing.T) {
	q := Select("*").From("users").Where(Ilike("name", "%john%")).Build()
	check(t, "WhereIlike", "select * from users where name ilike $1;", q.Stmt())
}

func TestWhereGreater(t *testing.T) {
	q := Select("*").From("products").Where(Greater("price", 100)).Build()
	check(t, "WhereGreater", "select * from products where price > $1;", q.Stmt())
}

func TestWhereLessOrEqual(t *testing.T) {
	q := Select("*").From("products").Where(LessOrEqual("stock", 10)).Build()
	check(t, "WhereLessOrEqual", "select * from products where stock <= $1;", q.Stmt())
}

func TestWhereNotEqual(t *testing.T) {
	q := Select("*").From("users").Where(NotEqual("status", "banned")).Build()
	check(t, "WhereNotEqual", "select * from users where status <> $1;", q.Stmt())
}

func TestWhereScope(t *testing.T) {
	q := Select("*").From("users").
		Where("active", true).
		Scope().Or("role", "admin").Or("role", "mod").EndScope().
		Build()
	check(t, "WhereScope",
		"select * from users where active=$1 or (role=$2 or role=$3);", q.Stmt())
}

func TestWhereRaw(t *testing.T) {
	q := Select("*").From("users").WhereRaw("id = $1", 42).Build()
	check(t, "WhereRaw stmt", "select * from users where id = $1;", q.Stmt())
	checkArgs(t, "WhereRaw args", []interface{}{42}, q.Args())
}

func TestWhereRawAfterWhere(t *testing.T) {
	q := Select("*").From("users").
		Where("active", true).
		WhereRaw("name ilike $1", "%john%").
		Build()
	check(t, "WhereRawAfterWhere",
		"select * from users where active=$1 and name ilike $2;", q.Stmt())
	checkArgs(t, "WhereRawAfterWhere args", []interface{}{true, "%john%"}, q.Args())
}

func TestOrRaw(t *testing.T) {
	q := Select("*").From("users").
		Where("role", "admin").
		OrRaw("email ilike $1", "%@company.com").
		Build()
	check(t, "OrRaw",
		"select * from users where role=$1 or email ilike $2;", q.Stmt())
}

// ── GROUP BY / HAVING ─────────────────────────────────────────────────────────

func TestGroupByHaving(t *testing.T) {
	q := Select("role", "count(*)").From("users").
		GroupBy("role").
		Having(Greater("count(*)", 5)).
		Build()
	check(t, "GroupByHaving",
		"select role,count(*) from users group by role having count(*) > $1;", q.Stmt())
}

// ── CONTAINS / STARTWITH / ENDWITH ───────────────────────────────────────────

func TestContains(t *testing.T) {
	q := Select("*").From("users").Where("active", true).Contains("name", "john").Build()
	check(t, "Contains",
		"select * from users where active=$1 and name ilike $2;", q.Stmt())
	checkArgs(t, "Contains args", []interface{}{true, "%john%"}, q.Args())
}

func TestStartWith(t *testing.T) {
	q := Select("*").From("users").Where("active", true).StartWith("email", "admin").Build()
	check(t, "StartWith",
		"select * from users where active=$1 and email ilike $2;", q.Stmt())
}

func TestEndWith(t *testing.T) {
	q := Select("*").From("users").Where("active", true).EndWith("email", "@acme.com").Build()
	check(t, "EndWith",
		"select * from users where active=$1 and email ilike $2;", q.Stmt())
}

// ── JOINS ─────────────────────────────────────────────────────────────────────

func TestInnerJoin(t *testing.T) {
	q := Select("u.id", "o.total").From("users u").
		InnerJoin("orders o").On("u.id", "o.user_id").
		Where("u.active", true).
		Build()
	check(t, "InnerJoin",
		"select u.id,o.total from users u inner join orders o on u.id=o.user_id where u.active=$1;",
		q.Stmt())
}

func TestLeftJoin(t *testing.T) {
	q := Select("u.id", "p.bio").From("users u").
		LeftJoin("profiles p").On("u.id", "p.user_id").
		Build()
	check(t, "LeftJoin",
		"select u.id,p.bio from users u left join profiles p on u.id=p.user_id;", q.Stmt())
}

func TestFullJoin(t *testing.T) {
	q := Select("a.id", "b.id").From("table_a a").
		FullJoin("table_b b").On("a.id", "b.ref_id").
		Build()
	check(t, "FullJoin",
		"select a.id,b.id from table_a a full join table_b b on a.id=b.ref_id;", q.Stmt())
}

// ── INSERT ────────────────────────────────────────────────────────────────────

func TestInsertBasic(t *testing.T) {
	q := Insert("users").
		Value("email", "alice@example.com").
		Value("age", 30).
		Return("id").
		Build()
	check(t, "InsertBasic",
		"insert into users (email,age) values($1,$2) returning id", q.Stmt())
	checkArgs(t, "InsertBasic args", []interface{}{"alice@example.com", 30}, q.Values())
}

func TestInsertOnConflict(t *testing.T) {
	q := Insert("users").
		Value("email", "alice@example.com").
		Value("age", 30).
		OnConflict("email").Update().
		Set("age", 31).
		Return("id").
		Build()
	stmt := q.Stmt()
	if !strings.Contains(stmt, "on conflict(email) do update set") {
		t.Errorf("InsertOnConflict: missing ON CONFLICT clause: %s", stmt)
	}
}

func TestInsertOnConflictNothing(t *testing.T) {
	q := Insert("products").
		Value("sku", "ABC-123").
		Value("price", 9.99).
		OnConflict("sku").Nothing().
		Build()
	stmt := q.Stmt()
	if !strings.Contains(stmt, "on conflict(sku) do nothing") {
		t.Errorf("OnConflictNothing: %s", stmt)
	}
}

// ── BATCH INSERT ──────────────────────────────────────────────────────────────

func TestBatchInsert(t *testing.T) {
	q := BatchInsert("users", "email", "age").
		Values("alice@example.com", 30).
		Values("bob@example.com", 25).
		Build()
	check(t, "BatchInsert",
		"insert into users(email,age) values($1, $2),($3, $4);", q.Stmt())
}

// ── UPDATE ────────────────────────────────────────────────────────────────────

func TestUpdateBasic(t *testing.T) {
	q := Update("users").
		Set("name", "Alice").
		Set("age", 31).
		Where("id", 1).
		Return("updated_at").
		Build()
	check(t, "UpdateBasic",
		"update users set name=$1,age=$2 where id=$3 returning updated_at;", q.Stmt())
	checkArgs(t, "UpdateBasic args", []interface{}{"Alice", 31, 1}, q.Values())
}

func TestUpdateOrNotIn(t *testing.T) {
	q := Update("users").
		Set("email", "mail@example.com").
		Set("age", 10).
		Where("id", 2).
		Or(NotIn("item", []int{0, 1})).
		Return("updated_at").
		Build()
	check(t, "UpdateOrNotIn",
		"update users set email=$1,age=$2 where id=$3 or item not in ($4,$5) returning updated_at;",
		q.Stmt())
}

// ── DELETE ────────────────────────────────────────────────────────────────────

func TestDeleteBasic(t *testing.T) {
	q := Delete("users").Where("id", 5).Build()
	check(t, "DeleteBasic", "delete from users where id=$1", q.Stmt())
}

func TestDeleteOrNotIn(t *testing.T) {
	q := Delete("users").
		Where("id", 2).
		Or(NotIn("item", []int{0, 1})).
		Build()
	check(t, "DeleteOrNotIn",
		"delete from users where id=$1 or item not in ($2,$3)", q.Stmt())
}

// ── WHERE EXISTS / SUBQUERY ───────────────────────────────────────────────────

func TestWhereExists(t *testing.T) {
	sub := Select("1").From("orders").Where("orders.user_id", 0) // placeholder
	// In practice the correlated column would come from the outer query
	q := Select("*").From("users").
		WhereExists(Select("1").From("orders").
			Where("orders.user_id", 99).Build()).
		Build()
	_ = sub
	stmt := q.Stmt()
	if !strings.Contains(stmt, "exists (select 1 from orders") {
		t.Errorf("WhereExists: missing exists clause: %s", stmt)
	}
}

func TestWhereNotExists(t *testing.T) {
	sub := Select("1").From("orders").Where("orders.user_id", 99).Build()
	q := Select("*").From("users").WhereNotExists(sub).Build()
	stmt := q.Stmt()
	if !strings.Contains(stmt, "not exists") {
		t.Errorf("WhereNotExists: %s", stmt)
	}
}

func TestWhereInSubquery(t *testing.T) {
	sub := Select("user_id").From("orders").Where(Greater("total", 1000)).Build()
	q := Select("*").From("users").WhereInSubquery("id", sub).Build()
	stmt := q.Stmt()
	if !strings.Contains(stmt, "id in (select user_id") {
		t.Errorf("WhereInSubquery: %s", stmt)
	}
}

// ── CTE ───────────────────────────────────────────────────────────────────────

func TestCTEBasic(t *testing.T) {
	active := Select("id", "name").From("users").Where("active", true).Build()
	q := With("active_users", active).
		Query(Select("*").From("active_users").Build()).
		Build()
	stmt := q.Stmt()
	if !strings.HasPrefix(stmt, "with active_users as (") {
		t.Errorf("CTEBasic: missing WITH prefix: %s", stmt)
	}
	if !strings.Contains(stmt, "select * from active_users") {
		t.Errorf("CTEBasic: missing main query: %s", stmt)
	}
	// arg from the sub-query must be present
	checkArgs(t, "CTEBasic args", []interface{}{true}, q.Args())
}

func TestCTEMultiple(t *testing.T) {
	recent := Select("id", "user_id").From("orders").Where("status", "paid").Build()
	active := Select("id", "name").From("users").Where("active", true).Build()

	q := NewCTE().
		Add("recent_orders", recent).
		Add("active_users", active).
		Query(Select("u.name", "o.id").From("active_users u").
			InnerJoin("recent_orders o").On("u.id", "o.user_id").Build()).
		Build()

	stmt := q.Stmt()
	if !strings.HasPrefix(stmt, "with recent_orders as (") {
		t.Errorf("CTEMultiple: wrong prefix: %s", stmt)
	}
	if !strings.Contains(stmt, "active_users as (") {
		t.Errorf("CTEMultiple: missing second CTE: %s", stmt)
	}
	// args: "paid"($1), true($2) — both shifted properly
	checkArgs(t, "CTEMultiple args", []interface{}{"paid", true}, q.Args())
}

func TestCTEArgRenumbering(t *testing.T) {
	// CTE uses $1 internally; main query also starts at $1.
	// After combination: CTE keeps $1, main query becomes $2.
	sub := Select("id").From("logs").Where("level", "error").Build()
	main := Select("*").From("users").Where("id", 1).Build()

	q := With("err_logs", sub).Query(main).Build()
	stmt := q.Stmt()
	// sub-query arg "error" is $1, main query arg 1 is shifted to $2
	if !strings.Contains(stmt, "where level=$1") {
		t.Errorf("CTEArgRenumbering: sub arg not $1: %s", stmt)
	}
	if !strings.Contains(stmt, "where id=$2") {
		t.Errorf("CTEArgRenumbering: main arg not $2: %s", stmt)
	}
	checkArgs(t, "CTEArgRenumbering args", []interface{}{"error", 1}, q.Args())
}

func TestCTERecursive(t *testing.T) {
	base := SQL("select 1 as n")
	q := WithRecursive("nums", base).
		AddRaw("", "select n+1 from nums where n < 10"). // illustrative raw
		Query(Select("*").From("nums").Build()).
		Build()
	stmt := q.Stmt()
	if !strings.HasPrefix(stmt, "with recursive") {
		t.Errorf("CTERecursive: missing WITH RECURSIVE: %s", stmt)
	}
}

// ── SECURITY: identifier injection ───────────────────────────────────────────

func TestSecurityTableNameInjection(t *testing.T) {
	q := Select("*").From("users; DROP TABLE users--")
	if q.Error() == nil {
		t.Error("expected error for injected table name")
	}
}

func TestSecurityColumnNameInjection(t *testing.T) {
	q := Select("*").From("users").Where("id; DROP TABLE users--", 1)
	if q.Error() == nil {
		t.Error("expected error for injected column name")
	}
}

func TestSecurityOrderByInjection(t *testing.T) {
	q := Select("*").From("users").OrderBy("id; DROP TABLE users--")
	if q.Error() == nil {
		t.Error("expected error for injected ORDER BY column")
	}
}

func TestSecurityGroupByInjection(t *testing.T) {
	q := Select("*").From("users").GroupBy("role; DROP TABLE users--")
	if q.Error() == nil {
		t.Error("expected error for injected GROUP BY column")
	}
}

func TestSecurityJoinInjection(t *testing.T) {
	q := Select("*").From("users").Join("orders; DROP TABLE orders--")
	if q.Error() == nil {
		t.Error("expected error for injected JOIN table")
	}
}

func TestSecurityOnInjection(t *testing.T) {
	q := Select("*").From("users").Join("orders").On("users.id; DROP TABLE users", "orders.user_id")
	if q.Error() == nil {
		t.Error("expected error for injected ON key")
	}
}

func TestSecurityInsertTableInjection(t *testing.T) {
	q := Insert("users; DROP TABLE users--")
	if q.Error() == nil {
		t.Error("expected error for injected INSERT table")
	}
}

func TestSecurityInsertColumnInjection(t *testing.T) {
	q := Insert("users").Value("email; DROP TABLE users--", "x")
	if q.Error() == nil {
		t.Error("expected error for injected INSERT column")
	}
}

func TestSecurityUpdateTableInjection(t *testing.T) {
	q := Update("users; DROP TABLE users--")
	if q.Error() == nil {
		t.Error("expected error for injected UPDATE table")
	}
}

func TestSecurityDeleteTableInjection(t *testing.T) {
	q := Delete("users; DROP TABLE users--")
	if q.Error() == nil {
		t.Error("expected error for injected DELETE table")
	}
}

// ── VALID IDENTIFIERS (regression: must not reject legitimate names) ──────────

func TestValidQualifiedIdentifiers(t *testing.T) {
	q := Select("u.id", "o.total").From("public.users u").
		InnerJoin("public.orders o").On("u.id", "o.user_id").
		Where("u.active", true).
		OrderBy("u.created_at").
		GroupBy("u.id").
		Build()
	if q.Error() != nil {
		t.Errorf("ValidQualifiedIdentifiers: unexpected error: %v", q.Error())
	}
}

// ── COALESCE / NULLIF (ValueFunc) ────────────────────────────────────────────

func TestCoalesce(t *testing.T) {
	q := Insert("users").
		Value("name", Coalesce("Alice", "unknown")).
		Build()
	stmt := q.Stmt()
	if !strings.Contains(stmt, "coalesce(") {
		t.Errorf("Coalesce: %s", stmt)
	}
}

// ── renumberArgs ──────────────────────────────────────────────────────────────

func TestRenumberArgs(t *testing.T) {
	got := renumberArgs("where id=$1 and role=$2", 3)
	check(t, "renumberArgs", "where id=$4 and role=$5", got)
}

func TestRenumberArgsZeroOffset(t *testing.T) {
	s := "where id=$1"
	got := renumberArgs(s, 0)
	check(t, "renumberArgsZero", s, got)
}

// ── SQL helper ────────────────────────────────────────────────────────────────

func TestSQLHelper(t *testing.T) {
	q := SQL("select 1").Where("id", 5).Build()
	check(t, "SQLHelper", "select 1 where id=$1;", q.Stmt())
}

// ── existing regression tests (updated expected values) ───────────────────────

func TestFilterQuery(t *testing.T) {
	q := Select("*").
		From("users").
		Where(In("age", []string{"76", "80"})).
		Or("email", "aaa@ajks.com").
		Where(Between("salary", []float64{5000, 5900})).
		Or(NotIn("role", []string{"admin", "driver"})).
		GroupBy("id", "age").
		Having(NotEqual("email", "aaa@ajks.com")).
		Or(NotIn("item", []int{0, 1})).
		Limit(90).
		Offset(7).
		OrderBy("id").
		Build()

	want := "select * from users where age in (76,80) or email=aaa@ajks.com and salary between 5000 and 5900 or role not in (admin,driver) group by id,age having email <> aaa@ajks.com or item not in (0,1) order by id limit 90 offset 7;"
	if want != q.Debug() {
		t.Errorf("TestFilterQuery:\n  want: %v\n   got: %v", want, q.Debug())
	}
}

func TestInsertQuery(t *testing.T) {
	q := Insert("users").
		Value("email", "mail@example.com").
		Value("age", 10).
		Return("id").
		Build()
	want := "insert into users (email,age) values(mail@example.com,10) returning id"
	if want != q.Debug() {
		t.Errorf("TestInsertQuery:\n  want: %v\n   got: %v", want, q.Debug())
	}
}

func TestUpdateQuery(t *testing.T) {
	q := Update("users").
		Set("email", "mail@example.com").
		Set("age", 10).
		Where("id", 2).
		Or(NotIn("item", []int{0, 1})).
		Return("updated_at").
		Build()
	want := "update users set email=mail@example.com,age=10 where id=2 or item not in (0,1) returning updated_at;"
	if want != q.Debug() {
		t.Errorf("TestUpdateQuery:\n  want: %v\n   got: %v", want, q.Debug())
	}
}

func TestDeleteQuery(t *testing.T) {
	q := Delete("users").
		Where("id", 2).
		Or(NotIn("item", []int{0, 1})).
		Build()
	want := "delete from users where id=2 or item not in (0,1)"
	if want != q.Debug() {
		t.Errorf("TestDeleteQuery:\n  want: %v\n   got: %v", want, q.Debug())
	}
}
