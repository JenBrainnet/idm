package employee

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	repo Repo
}

type Repo interface {
	Save(e *Entity) (int64, error)
	SaveTx(tx *sqlx.Tx, e *Entity) (int64, error)
	FindById(id int64) (Entity, error)
	FindAll() ([]Entity, error)
	FindAllByIds(ids []int64) ([]Entity, error)
	DeleteById(id int64) error
	DeleteAllByIds(ids []int64) error
	FindByNameTx(tx *sqlx.Tx, name string) (bool, error)
	BeginTransaction() (*sqlx.Tx, error)
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (svc *Service) Save(e *Entity) (int64, error) {
	id, err := svc.repo.Save(e)
	if err != nil {
		return 0, fmt.Errorf("error adding employee: %w", err)
	}
	return id, nil
}

func (svc *Service) SaveIfNameUnique(e *Entity) (id int64, err error) {
	tx, err := svc.repo.BeginTransaction()

	if err != nil {
		return 0, fmt.Errorf("error creating transaction: %w", err)
	}
	// отложенная функция завершения транзакции
	defer func() {
		if tx == nil {
			return // если транзакция не началась — ничего не делаем
		}
		// проверяем, не было ли паники
		if r := recover(); r != nil {
			err = fmt.Errorf("creating employee panic: %v", r)
			// если была паника, то откатываем транзакцию
			errTx := tx.Rollback()
			if errTx != nil {
				err = fmt.Errorf("creating employee: rolling back transaction errors: %w, %w", err, errTx)
			}
		} else if err != nil {
			// если произошла другая ошибка (не паника), то откатываем транзакцию
			errTx := tx.Rollback()
			if errTx != nil {
				err = fmt.Errorf("creating employee: rolling back transaction errors: %w, %w", err, errTx)
			}
		} else {
			// если ошибок нет, то коммитим транзакцию
			errTx := tx.Commit()
			if errTx != nil {
				err = fmt.Errorf("creating employee: commiting transaction error: %w", errTx)
			}
		}
	}()

	// выполняем несколько запросов в базе данных
	exists, err := svc.repo.FindByNameTx(tx, e.Name)
	if err != nil {
		return 0, fmt.Errorf("error checking name: %w", err)
	}
	if exists {
		return 0, fmt.Errorf("employee with name %s already exists", e.Name)
	}

	id, err = svc.repo.SaveTx(tx, e)
	if err != nil {
		return 0, fmt.Errorf("error saving employee: %w", err)
	}
	return id, nil
}

func (svc *Service) FindById(id int64) (Response, error) {
	entity, err := svc.repo.FindById(id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding employee with id %d: %w", id, err)
	}
	return entity.toResponse(), nil
}

func (svc *Service) FindAll() ([]Response, error) {
	entities, err := svc.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error retrieving all employees: %w", err)
	}

	var responses []Response
	for _, entity := range entities {
		responses = append(responses, entity.toResponse())
	}
	return responses, nil
}

func (svc *Service) FindAllByIds(ids []int64) ([]Response, error) {
	entities, err := svc.repo.FindAllByIds(ids)
	if err != nil {
		return nil, fmt.Errorf("error retrieving employees by ids %v: %w", ids, err)
	}

	var responses []Response
	for _, entity := range entities {
		responses = append(responses, entity.toResponse())
	}
	return responses, nil
}

func (svc *Service) DeleteById(id int64) error {
	err := svc.repo.DeleteById(id)
	if err != nil {
		return fmt.Errorf("error deleting employee with id %d: %w", id, err)
	}
	return nil
}

func (svc *Service) DeleteAllByIds(ids []int64) error {
	err := svc.repo.DeleteAllByIds(ids)
	if err != nil {
		return fmt.Errorf("error deleting employee by ids %v: %w", ids, err)
	}
	return nil
}
