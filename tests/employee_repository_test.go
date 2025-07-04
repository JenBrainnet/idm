package tests

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/database"
	"idm/inner/employee"
	"testing"
)

func TestEmployeeRepository(t *testing.T) {
	a := assert.New(t)
	db := database.ConnectDb()
	fixture := NewFixture(db)
	defer func() {
		if r := recover(); r != nil {
			fixture.ClearDatabase()
		}
	}()

	t.Run("find an employee by id", func(t *testing.T) {
		defer fixture.ClearDatabase()
		newEmployeeId := fixture.Employee("Employee Name")

		got, err := fixture.employees.FindById(newEmployeeId)
		a.Nil(err)
		a.NotEmpty(got)
		a.NotEmpty(got.Name)
		a.NotEmpty(got.CreatedAt)
		a.NotEmpty(got.UpdatedAt)
		a.Equal("Employee Name", got.Name)
	})

	t.Run("find all employees", func(t *testing.T) {
		defer fixture.ClearDatabase()
		fixture.Employee("Ivan")
		fixture.Employee("Stepan")

		got, err := fixture.employees.FindAll()
		a.Nil(err)
		a.Equal(2, len(got))
	})

	t.Run("find employees by ids", func(t *testing.T) {
		defer fixture.ClearDatabase()
		fixture.Employee("Ivan")
		id1 := fixture.Employee("Bob")
		id2 := fixture.Employee("Alice")

		got, err := fixture.employees.FindAllByIds([]int64{id1, id2})
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(2, len(got))
	})

	t.Run("delete employee by id", func(t *testing.T) {
		defer fixture.ClearDatabase()
		id := fixture.Employee("Alice")

		err := fixture.employees.DeleteById(id)
		a.Nil(err)

		_, err = fixture.employees.FindById(id)
		a.NotNil(err)
	})

	t.Run("delete employees by ids", func(t *testing.T) {
		defer fixture.ClearDatabase()
		id1 := fixture.Employee("Bob")
		id2 := fixture.Employee("Alice")

		err := fixture.employees.DeleteAllByIds([]int64{id1, id2})
		a.Nil(err)

		got, _ := fixture.employees.FindAllByIds([]int64{id1, id2})
		a.Len(got, 0)
	})

	t.Run("create employee in transaction", func(t *testing.T) {
		defer fixture.ClearDatabase()

		tx, err := fixture.db.Beginx()
		a.NoError(err)

		exists, err := fixture.employees.FindByNameTx(tx, "Alice")
		a.NoError(err)
		a.False(exists)

		entity := employee.Entity{Name: "Alice"}
		id, err := fixture.employees.SaveTx(tx, entity)
		a.NoError(err)
		a.NotZero(id)

		exists, err = fixture.employees.FindByNameTx(tx, entity.Name)
		a.NoError(err)
		a.True(exists)

		err = tx.Commit()
		a.NoError(err)

		saved, err := fixture.employees.FindById(id)
		a.NoError(err)
		a.Equal(entity.Name, saved.Name)
	})
}
