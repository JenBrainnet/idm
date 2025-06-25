package role

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
		return 0, fmt.Errorf("error adding role: %w", err)
	}
	return id, nil
}

func (svc *Service) FindById(id int64) (Response, error) {
	entity, err := svc.repo.FindById(id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding role with id %d: %w", id, err)
	}
	return entity.toResponse(), nil
}

func (svc *Service) FindAll() ([]Response, error) {
	entities, err := svc.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error retrieving all roles: %w", err)
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
		return nil, fmt.Errorf("error retrieving roles by ids %v: %w", ids, err)
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
		return fmt.Errorf("error deleting role with id %d: %w", id, err)
	}
	return nil
}

func (svc *Service) DeleteAllByIds(ids []int64) error {
	err := svc.repo.DeleteAllByIds(ids)
	if err != nil {
		return fmt.Errorf("error deleting role by ids %v: %w", ids, err)
	}
	return nil
}
