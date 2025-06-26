package employee

import (
	"errors"
	"fmt"
	"github.com/78bits/go-sqlmock-sqlx"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert" // импортируем библиотеку с ассерт-функциями
	"github.com/stretchr/testify/mock"   // импортируем пакет для создания моков
	"testing"
	"time"
)

// объявляем структуру мок-репозитория
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) SaveTx(tx *sqlx.Tx, e *Entity) (int64, error) {
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

func TestSave(t *testing.T) {
	a := assert.New(t)

	t.Run("should add employee", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		entity := &Entity{Name: "Alice"}
		repo.On("Save", entity).Return(int64(11), nil)

		id, err := svc.Save(entity)
		a.Nil(err)
		a.Equal(int64(11), id)
		a.True(repo.AssertNumberOfCalls(t, "Save", 1))
	})

	t.Run("should return wrapped error", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		entity := &Entity{Name: "Alice"}
		dbErr := errors.New("database error")
		want := fmt.Errorf("error adding employee: %w", dbErr)

		repo.On("Save", entity).Return(int64(0), dbErr)

		id, err := svc.Save(entity)
		a.Equal(int64(0), id)
		a.EqualError(err, want.Error())
		a.True(repo.AssertNumberOfCalls(t, "Save", 1))
	})
}

func TestSaveIfNameUnique(t *testing.T) {
	a := assert.New(t)

	t.Run("should return wrapped error when transaction begin fails", func(t *testing.T) {
		repo := new(MockRepo)
		svc := NewService(repo)

		entity := &Entity{Name: "Alice"}
		dbErr := errors.New("tx error")
		want := fmt.Errorf("error creating transaction: %w", dbErr)

		repo.On("BeginTransaction").Return((*sqlx.Tx)(nil), dbErr)

		_, err := svc.SaveIfNameUnique(entity)
		a.NotNil(err)
		a.EqualError(err, want.Error())
		a.True(repo.AssertNumberOfCalls(t, "BeginTransaction", 1))
	})

	t.Run("should return wrapped error when checking name fails", func(t *testing.T) {
		db, _, err := sqlmock.Newx()
		a.NoError(err)
		defer db.Close()

		repo := new(MockRepo)
		svc := NewService(repo)

		entity := &Entity{Name: "Alice"}
		tx, _ := db.Beginx()
		dbErr := errors.New("check error")
		want := fmt.Errorf("error checking name: %w", dbErr)

		repo.On("BeginTransaction").Return(tx, nil)
		repo.On("FindByNameTx", tx, entity.Name).Return(false, dbErr)

		_, err = svc.SaveIfNameUnique(entity)
		a.NotNil(err)
		a.EqualError(err, want.Error())
		a.True(repo.AssertNumberOfCalls(t, "BeginTransaction", 1))
		a.True(repo.AssertNumberOfCalls(t, "FindByNameTx", 1))
	})

	t.Run("should return error when employee already exists", func(t *testing.T) {
		db, _, err := sqlmock.Newx()
		a.NoError(err)
		defer db.Close()

		repo := new(MockRepo)
		svc := NewService(repo)

		entity := &Entity{Name: "Alice"}
		tx, _ := db.Beginx()
		want := fmt.Errorf("employee with name %s already exists", entity.Name)

		repo.On("BeginTransaction").Return(tx, nil)
		repo.On("FindByNameTx", tx, entity.Name).Return(true, nil)

		_, err = svc.SaveIfNameUnique(entity)
		a.NotNil(err)
		a.EqualError(err, want.Error())
		a.True(repo.AssertNumberOfCalls(t, "BeginTransaction", 1))
		a.True(repo.AssertNumberOfCalls(t, "FindByNameTx", 1))
	})

	t.Run("should return wrapped error when saving employee fails", func(t *testing.T) {
		db, _, err := sqlmock.Newx()
		a.NoError(err)
		defer db.Close()

		repo := new(MockRepo)
		svc := NewService(repo)

		entity := &Entity{Name: "Alice"}
		tx, _ := db.Beginx()
		dbErr := errors.New("save error")
		want := fmt.Errorf("error saving employee: %w", dbErr)

		repo.On("BeginTransaction").Return(tx, nil)
		repo.On("FindByNameTx", tx, entity.Name).Return(false, nil)
		repo.On("SaveTx", tx, entity).Return(int64(0), dbErr)

		_, err = svc.SaveIfNameUnique(entity)
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
		svc := NewService(repo)

		entity := &Entity{Name: "Alice"}
		tx, _ := db.Beginx()

		repo.On("BeginTransaction").Return(tx, nil)
		repo.On("FindByNameTx", tx, entity.Name).Return(false, nil)
		repo.On("SaveTx", tx, entity).Return(int64(1), nil)

		id, err := svc.SaveIfNameUnique(entity)
		a.NoError(err)
		a.Equal(int64(1), id)
		a.True(repo.AssertNumberOfCalls(t, "BeginTransaction", 1))
		a.True(repo.AssertNumberOfCalls(t, "FindByNameTx", 1))
		a.True(repo.AssertNumberOfCalls(t, "SaveTx", 1))
	})
}

func TestFindById(t *testing.T) {

	// создаём экземпляр объекта с ассерт-функциями
	a := assert.New(t)

	t.Run("should return found employee", func(t *testing.T) {

		// создаём экземпляр мок-объекта
		repo := new(MockRepo)

		// создаём экземпляр сервиса, который собираемся тестировать. Передаём в его конструктор мок вместо реального репозитория
		svc := NewService(repo)

		// создаём Entity, которую должен вернуть репозиторий
		entity := Entity{
			Id:        1,
			Name:      "John Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// создаём Response, который ожидаем получить от сервиса
		want := entity.toResponse()

		// конфигурируем поведение мок-репозитория (при вызове метода FindById с аргументом 1 вернуть Entity, созданную нами выше)
		repo.On("FindById", int64(1)).Return(entity, nil)

		// вызываем сервис с аргументом id = 1
		got, err := svc.FindById(1)

		// проверяем, что сервис не вернул ошибку
		a.Nil(err)
		// проверяем, что сервис вернул нам тот employee.Response, который мы ожилали получить
		a.Equal(want, got)
		// проверяем, что сервис вызвал репозиторий ровно 1 раз
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})

	t.Run("should return wrapped error", func(t *testing.T) {

		// Создаём для теста новый экземпляр мока репозитория.
		// Мы собираемся проверить счётчик вызовов, поэтому хотим, чтобы счётчик содержал количество вызовов к репозиторию,
		// выполненных в рамках одного нашего теста.
		// Ели сделать мок общим для нескольких тестов, то он посчитает вызовы, которые сделали все тесты
		repo := new(MockRepo)

		// создаём новый экземпляр сервиса (чтобы передать ему новый мок репозитория)
		svc := NewService(repo)

		// создаём пустую структуру employee.Entity, которую сервис вернёт вместе с ошибкой
		entity := Entity{}

		// ошибка, которую вернёт репозиторий
		dbErr := errors.New("database error")

		// ошибка, которую должен будет вернуть сервис
		want := fmt.Errorf("error finding employee with id 1: %w", dbErr)

		repo.On("FindById", int64(1)).Return(entity, dbErr)

		response, err := svc.FindById(1)

		// проверяем результаты теста
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
		want := fmt.Errorf("error retrieving all employees: %w", dbErr)

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

	t.Run("should return employees by ids", func(t *testing.T) {
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
		want := fmt.Errorf("error retrieving employees by ids %v: %w", ids, dbErr)

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

	t.Run("should delete employee by id", func(t *testing.T) {
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
		want := fmt.Errorf("error deleting employee with id %d: %w", 1, dbErr)

		repo.On("DeleteById", int64(1)).Return(dbErr)

		err := svc.DeleteById(1)
		a.NotNil(err)
		a.EqualError(err, want.Error())
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})
}

func TestDeleteAllByIds(t *testing.T) {
	a := assert.New(t)

	t.Run("should delete all employees by ids", func(t *testing.T) {
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
		want := fmt.Errorf("error deleting employee by ids %v: %w", ids, dbErr)

		repo.On("DeleteAllByIds", ids).Return(dbErr)

		err := svc.DeleteAllByIds(ids)
		a.NotNil(err)
		a.EqualError(err, want.Error())
		a.True(repo.AssertNumberOfCalls(t, "DeleteAllByIds", 1))
	})
}
