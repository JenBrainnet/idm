package role

import (
	"fmt"
	"idm/inner/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	Save(e Entity) (int64, error)
	FindById(id int64) (Entity, error)
	FindAll() ([]Entity, error)
	FindAllByIds(ids []int64) ([]Entity, error)
	DeleteById(id int64) error
	DeleteAllByIds(ids []int64) error
}

type Validator interface {
	Validate(request any) error
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{
		repo:      repo,
		validator: validator,
	}
}

func (svc *Service) Create(request CreateRequest) (int64, error) {
	err := svc.validator.Validate(request)
	if err != nil {
		// возвращаем кастомную ошибку в случае, если запрос не прошёл валидацию
		return 0, common.RequestValidationError{Message: err.Error()}
	}
	id, err := svc.repo.Save(request.ToEntity())
	if err != nil {
		return 0, fmt.Errorf("error saving role: %v", err)
	}
	return id, nil
}

func (svc *Service) FindById(request IdRequest) (Response, error) {
	err := svc.validator.Validate(request)
	if err != nil {
		return Response{}, common.RequestValidationError{Message: err.Error()}
	}
	entity, err := svc.repo.FindById(request.Id)
	if err != nil {
		return Response{}, common.NotFoundError{
			Message: fmt.Sprintf("error finding role with id %d: %v", request.Id, err),
		}
	}
	return entity.toResponse(), nil
}

func (svc *Service) FindAll() ([]Response, error) {
	entities, err := svc.repo.FindAll()
	if err != nil {
		return nil, common.NotFoundError{
			Message: fmt.Sprintf("error retrieving all roles: %v", err),
		}
	}
	var responses []Response
	for _, entity := range entities {
		responses = append(responses, entity.toResponse())
	}
	return responses, nil
}

func (svc *Service) FindAllByIds(request IdsRequest) ([]Response, error) {
	err := svc.validator.Validate(request)
	if err != nil {
		return nil, common.RequestValidationError{Message: err.Error()}
	}
	entities, err := svc.repo.FindAllByIds(request.Ids)
	if err != nil {
		return nil, common.NotFoundError{
			Message: fmt.Sprintf("error retrieving roles by ids %v: %v", request.Ids, err),
		}
	}
	var responses []Response
	for _, entity := range entities {
		responses = append(responses, entity.toResponse())
	}
	return responses, nil
}

func (svc *Service) DeleteById(request IdRequest) error {
	err := svc.validator.Validate(request)
	if err != nil {
		return common.RequestValidationError{Message: err.Error()}
	}
	err = svc.repo.DeleteById(request.Id)
	if err != nil {
		return common.NotFoundError{
			Message: fmt.Sprintf("error deleting role with id %d: %v", request.Id, err),
		}
	}
	return nil
}

func (svc *Service) DeleteAllByIds(request IdsRequest) error {
	err := svc.validator.Validate(request)
	if err != nil {
		return common.RequestValidationError{Message: err.Error()}
	}
	err = svc.repo.DeleteAllByIds(request.Ids)
	if err != nil {
		return common.NotFoundError{
			Message: fmt.Sprintf("error deleting roles by ids %v: %v", request.Ids, err),
		}
	}
	return nil
}
