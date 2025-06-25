package employee

import "fmt"

type Service struct {
	repo Repo
}

type Repo interface {
	Save(e *Entity) (int64, error)
	FindById(id int64) (Entity, error)
	FindAll() ([]Entity, error)
	FindAllByIds(ids []int64) ([]Entity, error)
	DeleteById(id int64) error
	DeleteAllByIds(ids []int64) error
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
