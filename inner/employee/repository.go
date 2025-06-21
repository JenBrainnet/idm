package employee

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type EmployeeRepository struct {
	db *sqlx.DB
}

func NewEmployeeRepository(database *sqlx.DB) *EmployeeRepository {
	return &EmployeeRepository{db: database}
}

type EmployeeEntity struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *EmployeeRepository) Create(employee *EmployeeEntity) (id int64, err error) {
	query := "insert into employee (name) values ($1) returning id"
	err = r.db.QueryRowx(query, employee.Name).Scan(&id)
	return id, err
}

func (r *EmployeeRepository) FindById(id int64) (employee EmployeeEntity, err error) {
	query := "select * from employee where id = $1"
	err = r.db.Get(&employee, query, id)
	return employee, err
}

func (r *EmployeeRepository) FindAll() (employees []EmployeeEntity, err error) {
	query := "select * from employee"
	err = r.db.Select(&employees, query)
	return employees, err
}

func (r *EmployeeRepository) FindAllByIds(ids []int64) (employees []EmployeeEntity, err error) {
	if len(ids) == 0 {
		return []EmployeeEntity{}, nil
	}
	query := "select * from employee where id = ANY($1)"
	err = r.db.Select(&employees, query, pq.Array(ids))
	return employees, err
}

func (r *EmployeeRepository) DeleteById(id int64) (err error) {
	query := "delete from employee where id=$1"
	_, err = r.db.Exec(query, id)
	return err
}

func (r *EmployeeRepository) DeleteAllByIds(ids []int64) (err error) {
	if len(ids) == 0 {
		return nil
	}
	query := "delete from employee where id = ANY($1)"
	_, err = r.db.Exec(query, pq.Array(ids))
	return err
}
