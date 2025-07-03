package employee

import (
	"errors"
	"fmt"
	"github.com/78bits/go-sqlmock-sqlx"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert" // импортируем библиотеку с ассерт-функциями
	"github.com/stretchr/testify/mock"   // импортируем пакет для создания моков
	"idm/inner/common"
	"idm/inner/validator"
	"testing"
	"time"
)

// объявляем структуру мок-репозитория
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) SaveTx(tx *sqlx.Tx, e Entity) (int64, error) {
	args := m.Called(tx, e)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) FindByNameTx(tx *sqlx.Tx, name string) (bool, error) {
	args := m.Called(tx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepo) BeginTransaction() (*sqlx.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sqlx.Tx), args.Error(1)
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

// реализуем интерфейс репозитория у мока
func (m *MockRepo) FindById(id int64) (employee Entity, err error) {

	// Общая конфигурация поведения мок-объекта
	args := m.Called(id)
	return args.Get(0).(Entity), args.Error(1)
}

func TestCreate(t *testing.T) {
	a := assert.New(t)

	t.Run("should return wrapped error when transaction begin fails", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		a.NoError(err)
		sqlxDB := sqlx.NewDb(db, "sqlmock")

		repo := &Repository{db: sqlxDB}
		svc := NewService(repo, validator.New())

		// создаём ошибку, которую должен вернуть Begin
		dbErr := errors.New("transaction begin error")
		want := fmt.Errorf("error creating transaction: %w", dbErr)

		// sqlmock должен сымитировать ошибку начала транзакции
		sqlMock.ExpectBegin().WillReturnError(dbErr)

		id, err := svc.Create(CreateRequest{Name: "test"})
		a.Equal(int64(0), id)
		a.NotNil(err)
		a.EqualError(err, want.Error())
	})

	t.Run("should return error when employee already exists", func(t *testing.T) {
		db, sqlMock, err := sqlmock.Newx()
		a.NoError(err)
		defer db.Close()

		sqlMock.ExpectBegin()

		tx, err := db.Beginx()
		a.NoError(err)

		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		entity := Entity{Name: "Alice"}
		want := common.AlreadyExistsError{
			Message: fmt.Sprintf("employee with name %s already exists", entity.Name),
		}

		repo.On("BeginTransaction").Return(tx, nil)
		repo.On("FindByNameTx", tx, entity.Name).Return(true, nil)

		id, err := svc.Create(CreateRequest{Name: entity.Name})

		a.Equal(int64(0), id)
		a.NotNil(err)
		a.Equal(err, want)
		a.True(repo.AssertNumberOfCalls(t, "BeginTransaction", 1))
		a.True(repo.AssertNumberOfCalls(t, "FindByNameTx", 1))
		a.NoError(sqlMock.ExpectationsWereMet())
	})

	t.Run("should return wrapped error when saving employee fails", func(t *testing.T) {
		db, _, err := sqlmock.Newx()
		a.NoError(err)
		defer db.Close()

		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		entity := Entity{Name: "Alice"}
		tx, _ := db.Beginx()
		dbErr := errors.New("save error")
		want := fmt.Errorf("error saving employee with name: %s %w", entity.Name, dbErr)

		repo.On("BeginTransaction").Return(tx, nil)
		repo.On("FindByNameTx", tx, entity.Name).Return(false, nil)
		repo.On("SaveTx", tx, entity).Return(int64(0), dbErr)

		_, err = svc.Create(CreateRequest{Name: entity.Name})
		a.NotNil(err)
		a.EqualError(err, want.Error())
		a.True(repo.AssertNumberOfCalls(t, "BeginTransaction", 1))
		a.True(repo.AssertNumberOfCalls(t, "FindByNameTx", 1))
		a.True(repo.AssertNumberOfCalls(t, "SaveTx", 1))
	})

	t.Run("should save employee", func(t *testing.T) {
		db, _, err := sqlmock.Newx()
		a.NoError(err)
		defer db.Close()

		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		entity := Entity{Name: "Alice"}
		tx, _ := db.Beginx()

		repo.On("BeginTransaction").Return(tx, nil)
		repo.On("FindByNameTx", tx, entity.Name).Return(false, nil)
		repo.On("SaveTx", tx, entity).Return(int64(1), nil)

		id, err := svc.Create(CreateRequest{Name: entity.Name})
		a.NoError(err)
		a.Equal(int64(1), id)
		a.True(repo.AssertNumberOfCalls(t, "BeginTransaction", 1))
		a.True(repo.AssertNumberOfCalls(t, "FindByNameTx", 1))
		a.True(repo.AssertNumberOfCalls(t, "SaveTx", 1))
	})
}

func TestFindById(t *testing.T) {
	a := assert.New(t)

	t.Run("should return found employee", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		entity := Entity{Id: 1, Name: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		want := entity.toResponse()

		repo.On("FindById", int64(1)).Return(entity, nil)

		got, err := svc.FindById(IdRequest{Id: 1})
		a.Nil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})

	t.Run("should return not found error", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		// создаём пустую структуру employee.Entity, которую сервис вернёт вместе с ошибкой
		entity := Entity{}

		// ошибка, которую вернёт репозиторий
		dbErr := errors.New("database error")

		// ошибка, которую должен будет вернуть сервис
		want := common.NotFoundError{
			Message: fmt.Sprintf("error finding employee with id 1: %v", dbErr),
		}

		repo.On("FindById", int64(1)).Return(entity, dbErr)

		response, err := svc.FindById(IdRequest{Id: 1})
		a.Empty(response)
		a.NotNil(err)
		a.Equal(want, err)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})
}

func TestFindAll(t *testing.T) {
	a := assert.New(t)

	t.Run("should return all employees", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		entities := []Entity{
			{Id: 1, Name: "First", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Id: 2, Name: "Second", CreatedAt: time.Now(), UpdatedAt: time.Now()},
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

	t.Run("should return not found error", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		dbErr := errors.New("database error")
		want := common.NotFoundError{
			Message: fmt.Sprintf("error retrieving all employees: %v", dbErr),
		}

		repo.On("FindAll").Return([]Entity{}, dbErr)

		got, err := svc.FindAll()
		a.Nil(got)
		a.NotNil(err)
		a.Equal(err, want)
		a.True(repo.AssertNumberOfCalls(t, "FindAll", 1))
	})
}

func TestFindAllByIds(t *testing.T) {
	a := assert.New(t)

	t.Run("should return employees by ids", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		ids := []int64{1, 2}
		entities := []Entity{
			{Id: 1, Name: "First", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Id: 2, Name: "Second", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		want := []Response{
			entities[0].toResponse(),
			entities[1].toResponse(),
		}
		repo.On("FindAllByIds", ids).Return(entities, nil)

		got, err := svc.FindAllByIds(IdsRequest{Ids: ids})
		a.Nil(err)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindAllByIds", 1))
	})

	t.Run("should return not found error", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		ids := []int64{1, 2}
		dbErr := errors.New("database error")
		want := common.NotFoundError{
			Message: fmt.Sprintf("error retrieving employees by ids %v: %v", ids, dbErr),
		}

		repo.On("FindAllByIds", ids).Return([]Entity{}, dbErr)

		got, err := svc.FindAllByIds(IdsRequest{Ids: ids})
		a.Nil(got)
		a.NotNil(err)
		a.Equal(err, want)
		a.True(repo.AssertNumberOfCalls(t, "FindAllByIds", 1))
	})
}

func TestDeleteById(t *testing.T) {
	a := assert.New(t)

	t.Run("should delete employee by id", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		repo.On("DeleteById", int64(1)).Return(nil)

		err := svc.DeleteById(IdRequest{Id: 1})
		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})

	t.Run("should return not found error", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		dbErr := errors.New("database error")
		want := common.NotFoundError{
			Message: fmt.Sprintf("error deleting employee with id %d: %v", 1, dbErr),
		}

		repo.On("DeleteById", int64(1)).Return(dbErr)

		err := svc.DeleteById(IdRequest{Id: 1})
		a.NotNil(err)
		a.Equal(err, want)
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})
}

func TestDeleteAllByIds(t *testing.T) {
	a := assert.New(t)

	t.Run("should delete all employees by ids", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		ids := []int64{1, 2}
		repo.On("DeleteAllByIds", ids).Return(nil)

		err := svc.DeleteAllByIds(IdsRequest{Ids: ids})
		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "DeleteAllByIds", 1))
	})

	t.Run("should return not found error", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo, validator.New())

		ids := []int64{1, 2}
		dbErr := errors.New("database error")
		want := common.NotFoundError{
			Message: fmt.Sprintf("error deleting employees by ids %v: %v", ids, dbErr),
		}

		repo.On("DeleteAllByIds", ids).Return(dbErr)

		err := svc.DeleteAllByIds(IdsRequest{Ids: ids})
		a.NotNil(err)
		a.Equal(err, want)
		a.True(repo.AssertNumberOfCalls(t, "DeleteAllByIds", 1))
	})
}
