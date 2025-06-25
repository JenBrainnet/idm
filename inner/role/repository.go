package role

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

func (r *Repository) Save(role *Entity) (id int64, err error) {
	query := "insert into role (name) values ($1) returning id"
	err = r.db.QueryRowx(query, role.Name).Scan(&id)
	return id, err
}

func (r *Repository) FindById(id int64) (role Entity, err error) {
	query := "select * from role where id=$1"
	err = r.db.Get(&role, query, id)
	return role, err
}

func (r *Repository) FindAll() (roles []Entity, err error) {
	query := "select * from role"
	err = r.db.Select(&roles, query)
	return roles, err
}

func (r *Repository) FindAllByIds(ids []int64) (roles []Entity, err error) {
	if len(ids) == 0 {
		return []Entity{}, nil
	}
	query := "select * from role where id = ANY($1)"
	err = r.db.Select(&roles, query, pq.Array(ids))
	return roles, err
}

func (r *Repository) DeleteById(id int64) (err error) {
	query := "delete from role where id=$1"
	_, err = r.db.Exec(query, id)
	return err
}

func (r *Repository) DeleteAllByIds(ids []int64) (err error) {
	if len(ids) == 0 {
		return nil
	}
	query := "delete from role where id = ANY($1)"
	_, err = r.db.Exec(query, pq.Array(ids))
	return err
}
