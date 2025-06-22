package tests

import (
	"github.com/jmoiron/sqlx"
	"idm/inner/employee"
	"idm/inner/role"
)

type Fixture struct {
	db        *sqlx.DB
	employees *employee.Repository
	roles     *role.Repository
}

func NewFixture(db *sqlx.DB) *Fixture {
	initSchema(db)
	return &Fixture{
		db:        db,
		employees: employee.NewRepository(db),
		roles:     role.NewRepository(db),
	}
}

func initSchema(db *sqlx.DB) {
	schema := `
	create table if not exists role (
    	id bigint primary key generated always as identity,
    	name text not null,
    	created_at timestamptz not null default now(),
    	updated_at timestamptz not null default now()
	);

	create table if not exists employee (
    	id bigint primary key generated always as identity,
    	name text not null,
    	created_at timestamptz not null default now(),
    	updated_at timestamptz not null default now()
	);`
	db.MustExec(schema)
}

func (f *Fixture) Employee(name string) int64 {
	entity := employee.Entity{
		Name: name,
	}
	newId, err := f.employees.Add(&entity)
	if err != nil {
		panic(err)
	}
	return newId
}

func (f *Fixture) Role(name string) int64 {
	entity := role.Entity{
		Name: name,
	}
	newId, err := f.roles.Add(&entity)
	if err != nil {
		panic(err)
	}
	return newId
}

func (f *Fixture) ClearDatabase() {
	f.db.MustExec("delete from employee")
	f.db.MustExec("delete from role")
}
