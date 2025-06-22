package tests

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/database"
	"testing"
)

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
