package role

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type RoleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(database *sqlx.DB) *RoleRepository {
	return &RoleRepository{db: database}
}

type RoleEntity struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *RoleRepository) Create(role *RoleEntity) (id int64, err error) {
	query := "insert into role (name) values ($1) returning id"
	err = r.db.QueryRowx(query, role.Name).Scan(&id)
	return id, err
}

func (r *RoleRepository) FindById(id int64) (role RoleEntity, err error) {
	query := "select * from role where id=$1"
	err = r.db.Get(&role, query, id)
	return role, err
}

func (r *RoleRepository) FindAll() (roles []RoleEntity, err error) {
	query := "select * from role"
	err = r.db.Select(&roles, query)
	return roles, err
}

func (r *RoleRepository) FindAllByIds(ids []int64) (roles []RoleEntity, err error) {
	if len(ids) == 0 {
		return []RoleEntity{}, nil
	}
	query := "select * from role where id = ANY($1)"
	err = r.db.Select(&roles, query, pq.Array(ids))
	return roles, err
}

func (r *RoleRepository) DeleteById(id int64) (err error) {
	query := "delete from role where id=$1"
	_, err = r.db.Exec(query, id)
	return err
}

func (r *RoleRepository) DeleteAllByIds(ids []int64) (err error) {
	if len(ids) == 0 {
		return nil
	}
	query := "delete from role where id = ANY($1)"
	_, err = r.db.Exec(query, pq.Array(ids))
	return err
}
