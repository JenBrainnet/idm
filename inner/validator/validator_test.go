package validator_test

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/employee"
	"idm/inner/role"
	"idm/inner/validator"
	"strings"
	"testing"
)

func TestValidatorEmployeeCreateRequest(t *testing.T) {
	a := assert.New(t)
	v := validator.New()

	t.Run("should pass with valid name", func(t *testing.T) {
		request := employee.CreateRequest{Name: "John"}
		err := v.Validate(request)
		a.NoError(err)
	})

	t.Run("should fail for empty name", func(t *testing.T) {
		request := employee.CreateRequest{Name: ""}
		err := v.Validate(request)
		a.Error(err)
	})

	t.Run("should fail for short name", func(t *testing.T) {
		request := employee.CreateRequest{Name: "J"}
		err := v.Validate(request)
		a.Error(err)
	})

	t.Run("should fail for long name", func(t *testing.T) {
		request := employee.CreateRequest{Name: strings.Repeat("J", 156)}
		err := v.Validate(request)
		a.Error(err)
	})
}

func TestValidatorEmployeeIdRequest(t *testing.T) {
	a := assert.New(t)
	v := validator.New()

	t.Run("should pass for valid id", func(t *testing.T) {
		err := v.Validate(employee.IdRequest{Id: 1})
		a.NoError(err)
	})

	t.Run("should fail for zero id", func(t *testing.T) {
		err := v.Validate(employee.IdRequest{Id: 0})
		a.Error(err)
	})
}

func TestValidatorEmployeeIdsRequest(t *testing.T) {
	a := assert.New(t)
	v := validator.New()

	t.Run("should pass for valid ids", func(t *testing.T) {
		err := v.Validate(employee.IdsRequest{Ids: []int64{1, 2}})
		a.NoError(err)
	})

	t.Run("should fail for empty ids", func(t *testing.T) {
		err := v.Validate(employee.IdsRequest{Ids: []int64{}})
		a.Error(err)
	})
}

func TestValidatorRoleCreateRequest(t *testing.T) {
	a := assert.New(t)
	v := validator.New()

	t.Run("should pass with valid name", func(t *testing.T) {
		request := role.CreateRequest{Name: "Admin"}
		err := v.Validate(request)
		a.NoError(err)
	})

	t.Run("should fail for empty name", func(t *testing.T) {
		request := role.CreateRequest{Name: ""}
		err := v.Validate(request)
		a.Error(err)
	})

	t.Run("should fail for short name", func(t *testing.T) {
		request := role.CreateRequest{Name: "A"}
		err := v.Validate(request)
		a.Error(err)
	})

	t.Run("should fail for long name", func(t *testing.T) {
		request := role.CreateRequest{Name: strings.Repeat("A", 56)}
		err := v.Validate(request)
		a.Error(err)
	})
}

func TestValidatorRoleIdRequest(t *testing.T) {
	a := assert.New(t)
	v := validator.New()

	t.Run("should pass for valid id", func(t *testing.T) {
		err := v.Validate(role.IdRequest{Id: 1})
		a.NoError(err)
	})

	t.Run("should fail for zero id", func(t *testing.T) {
		err := v.Validate(role.IdRequest{Id: 0})
		a.Error(err)
	})
}

func TestValidatorRoleIdsRequest(t *testing.T) {
	a := assert.New(t)
	v := validator.New()

	t.Run("should pass for valid ids", func(t *testing.T) {
		err := v.Validate(role.IdsRequest{Ids: []int64{1, 2}})
		a.NoError(err)
	})

	t.Run("should fail for empty ids", func(t *testing.T) {
		err := v.Validate(role.IdsRequest{Ids: []int64{}})
		a.Error(err)
	})
}
