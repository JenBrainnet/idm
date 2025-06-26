package employee

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(database *sqlx.DB) *Repository {
	return &Repository{db: database}
}

func (r *Repository) Save(employee *Entity) (id int64, err error) {
	query := "insert into employee (name) values ($1) returning id"
	err = r.db.QueryRowx(query, employee.Name).Scan(&id)
	return id, err
}

func (r *Repository) FindById(id int64) (employee Entity, err error) {
	query := "select * from employee where id = $1"
	err = r.db.Get(&employee, query, id)
	return employee, err
}

func (r *Repository) FindAll() (employees []Entity, err error) {
	query := "select * from employee"
	err = r.db.Select(&employees, query)
	return employees, err
}

func (r *Repository) FindAllByIds(ids []int64) (employees []Entity, err error) {
	if len(ids) == 0 {
		return []Entity{}, nil
	}
	query := "select * from employee where id = ANY($1)"
	err = r.db.Select(&employees, query, pq.Array(ids))
	return employees, err
}

func (r *Repository) DeleteById(id int64) (err error) {
	query := "delete from employee where id=$1"
	_, err = r.db.Exec(query, id)
	return err
}

func (r *Repository) DeleteAllByIds(ids []int64) (err error) {
	if len(ids) == 0 {
		return nil
	}
	query := "delete from employee where id = ANY($1)"
	_, err = r.db.Exec(query, pq.Array(ids))
	return err
}

func (r *Repository) BeginTransaction() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *Repository) FindByNameTx(tx *sqlx.Tx, name string) (exists bool, err error) {
	query := "select exists(select 1 from employee where name = $1)"
	err = tx.Get(&exists, query, name)
	return exists, err
}

func (r *Repository) SaveTx(tx *sqlx.Tx, e *Entity) (id int64, err error) {
	query := "insert into employee (name) values ($1) returning id"
	err = tx.QueryRowx(query, e.Name).Scan(&id)
	return id, err
}
