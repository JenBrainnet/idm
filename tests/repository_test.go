package tests

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/database"
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
}

func TestRoleRepository(t *testing.T) {
	a := assert.New(t)
	db := database.ConnectDb()
	fixture := NewFixture(db)
	defer func() {
		if r := recover(); r != nil {
			fixture.ClearDatabase()
		}
	}()

	t.Run("find a role by id", func(t *testing.T) {
		defer fixture.ClearDatabase()
		newRoleId := fixture.Role("Role Name")

		got, err := fixture.roles.FindById(newRoleId)
		a.Nil(err)
		a.NotEmpty(got)
		a.NotEmpty(got.Name)
		a.NotEmpty(got.CreatedAt)
		a.NotEmpty(got.UpdatedAt)
		a.Equal("Role Name", got.Name)
	})

	t.Run("find all roles", func(t *testing.T) {
		defer fixture.ClearDatabase()
		fixture.Role("Admin")
		fixture.Role("User")

		got, err := fixture.roles.FindAll()
		a.Nil(err)
		a.Equal(2, len(got))
	})

	t.Run("find roles by ids", func(t *testing.T) {
		defer fixture.ClearDatabase()
		fixture.Role("Manager")
		id1 := fixture.Role("Admin")
		id2 := fixture.Role("User")

		got, err := fixture.roles.FindAllByIds([]int64{id1, id2})
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(2, len(got))
	})

	t.Run("delete role by id", func(t *testing.T) {
		defer fixture.ClearDatabase()
		id := fixture.Role("Admin")

		err := fixture.roles.DeleteById(id)
		a.Nil(err)

		_, err = fixture.roles.FindById(id)
		a.NotNil(err)
	})

	t.Run("delete roles by ids", func(t *testing.T) {
		defer fixture.ClearDatabase()
		id1 := fixture.Role("Admin")
		id2 := fixture.Role("User")

		err := fixture.roles.DeleteAllByIds([]int64{id1, id2})
		a.Nil(err)

		got, _ := fixture.roles.FindAllByIds([]int64{id1, id2})
		a.Len(got, 0)
	})

}
