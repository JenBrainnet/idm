package role

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type StubRepo struct {
	FindByIdResult Entity
	FindByIdError  error
}

func (s *StubRepo) Save(e *Entity) (int64, error) {
	panic("implement me")
}

func (s *StubRepo) FindById(id int64) (Entity, error) {
	return s.FindByIdResult, s.FindByIdError
}

func (s *StubRepo) FindAll() ([]Entity, error) {
	panic("implement me")
}

func (s *StubRepo) FindAllByIds(ids []int64) ([]Entity, error) {
	panic("implement me")
}

func (s *StubRepo) DeleteById(id int64) error {
	panic("implement me")
}

func (s *StubRepo) DeleteAllByIds(ids []int64) error {
	panic("implement me")
}

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Save(e *Entity) (int64, error) {
	args := m.Called(e)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) FindAll() ([]Entity, error) {
	args := m.Called()
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *MockRepo) FindAllByIds(ids []int64) ([]Entity, error) {
	args := m.Called(ids)
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *MockRepo) DeleteById(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) DeleteAllByIds(ids []int64) error {
	args := m.Called(ids)
	return args.Error(0)
}

func (m *MockRepo) FindById(id int64) (role Entity, err error) {
	args := m.Called(id)
	return args.Get(0).(Entity), args.Error(1)
}

func TestSave(t *testing.T) {
	a := assert.New(t)

	t.Run("should add role", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		entity := &Entity{Name: "Programmer"}
		repo.On("Save", entity).Return(int64(11), nil)

		id, err := svc.Save(entity)
		a.Nil(err)
		a.Equal(int64(11), id)
		a.True(repo.AssertNumberOfCalls(t, "Save", 1))
	})

	t.Run("should return wrapped error", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		entity := &Entity{Name: "Admin"}
		dbErr := errors.New("database error")
		want := fmt.Errorf("error adding role: %w", dbErr)

		repo.On("Save", entity).Return(int64(0), dbErr)

		id, err := svc.Save(entity)
		a.Equal(int64(0), id)
		a.EqualError(err, want.Error())
		a.True(repo.AssertNumberOfCalls(t, "Save", 1))
	})
}

func TestFindById(t *testing.T) {
	a := assert.New(t)

	t.Run("should return found role", func(t *testing.T) {
		stubRepo := &StubRepo{
			FindByIdResult: Entity{
				Id:        1,
				Name:      "Admin",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			FindByIdError: nil,
		}
		svc := NewService(stubRepo)

		want := stubRepo.FindByIdResult.toResponse()

		got, err := svc.FindById(1)
		a.Nil(err)
		a.Equal(want, got)
	})

	t.Run("should return wrapped error", func(t *testing.T) {
		stubRepo := &StubRepo{
			FindByIdResult: Entity{},
			FindByIdError:  errors.New("database error"),
		}
		svc := NewService(stubRepo)

		want := fmt.Errorf("error finding role with id 1: %w", stubRepo.FindByIdError)

		got, err := svc.FindById(1)
		a.Empty(got)
		a.NotNil(err)
		a.Equal(want, err)
	})
}

func TestFindAll(t *testing.T) {
	a := assert.New(t)

	t.Run("should return all roles", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		entities := []Entity{
			{
				Id:        1,
				Name:      "First",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				Id:        2,
				Name:      "Second",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		want := []Response{
			entities[0].toResponse(),
			entities[1].toResponse(),
		}
		repo.On("FindAll").Return(entities, nil)

		got, err := svc.FindAll()
		a.Nil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindAll", 1))
	})

	t.Run("should return wrapped error", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		dbErr := errors.New("database error")
		want := fmt.Errorf("error retrieving all roles: %w", dbErr)

		repo.On("FindAll").Return([]Entity{}, dbErr)

		got, err := svc.FindAll()
		a.Nil(got)
		a.NotNil(err)
		a.EqualError(err, want.Error())
		a.True(repo.AssertNumberOfCalls(t, "FindAll", 1))
	})
}

func TestFindAllByIds(t *testing.T) {
	a := assert.New(t)

	t.Run("should return roles by ids", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		ids := []int64{1, 2}
		entities := []Entity{
			{
				Id:        1,
				Name:      "First",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				Id:        2,
				Name:      "Second",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		want := []Response{
			entities[0].toResponse(),
			entities[1].toResponse(),
		}
		repo.On("FindAllByIds", ids).Return(entities, nil)

		got, err := svc.FindAllByIds(ids)
		a.Nil(err)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindAllByIds", 1))
	})

	t.Run("should return wrapped error", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		ids := []int64{1, 2}
		dbErr := errors.New("database error")
		want := fmt.Errorf("error retrieving roles by ids %v: %w", ids, dbErr)

		repo.On("FindAllByIds", ids).Return([]Entity{}, dbErr)

		got, err := svc.FindAllByIds(ids)
		a.Nil(got)
		a.NotNil(err)
		a.EqualError(err, want.Error())
		a.True(repo.AssertNumberOfCalls(t, "FindAllByIds", 1))
	})
}

func TestDeleteById(t *testing.T) {
	a := assert.New(t)

	t.Run("should delete role by id", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		repo.On("DeleteById", int64(1)).Return(nil)

		err := svc.DeleteById(1)
		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})

	t.Run("should return wrapped error", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		dbErr := errors.New("database error")
		want := fmt.Errorf("error deleting role with id %d: %w", 1, dbErr)

		repo.On("DeleteById", int64(1)).Return(dbErr)

		err := svc.DeleteById(1)
		a.NotNil(err)
		a.EqualError(err, want.Error())
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})
}

func TestDeleteAllByIds(t *testing.T) {
	a := assert.New(t)

	t.Run("should delete all roles by ids", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		ids := []int64{1, 2}
		repo.On("DeleteAllByIds", ids).Return(nil)

		err := svc.DeleteAllByIds(ids)
		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "DeleteAllByIds", 1))
	})

	t.Run("should return wrapped error", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		ids := []int64{1, 2}
		dbErr := errors.New("database error")
		want := fmt.Errorf("error deleting role by ids %v: %w", ids, dbErr)

		repo.On("DeleteAllByIds", ids).Return(dbErr)

		err := svc.DeleteAllByIds(ids)
		a.NotNil(err)
		a.EqualError(err, want.Error())
		a.True(repo.AssertNumberOfCalls(t, "DeleteAllByIds", 1))
	})
}
