package tests

import (
	"github.com/jmoiron/sqlx"
	"idm/inner/employee"
	"idm/inner/role"
)

type Fixture struct {
	db        *sqlx.DB
	employees *employee.EmployeeRepository
	roles     *role.RoleRepository
}

func NewFixture(db *sqlx.DB) *Fixture {
	initSchema(db)
	return &Fixture{
		db:        db,
		employees: employee.NewEmployeeRepository(db),
		roles:     role.NewRoleRepository(db),
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
	entity := employee.EmployeeEntity{
		Name: name,
	}
	newId, err := f.employees.Create(&entity)
	if err != nil {
		panic(err)
	}
	return newId
}

func (f *Fixture) Role(name string) int64 {
	entity := role.RoleEntity{
		Name: name,
	}
	newId, err := f.roles.Create(&entity)
	if err != nil {
		panic(err)
	}
	return newId
}

func (f *Fixture) ClearDatabase() {
	f.db.MustExec("delete from employee")
	f.db.MustExec("delete from role")
}
