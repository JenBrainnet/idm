package role

type CreateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=55"`
}

func (r *CreateRequest) ToEntity() Entity {
	return Entity{Name: r.Name}
}

type IdRequest struct {
	Id int64 `json:"id" validate:"required,gt=0"`
}

type IdsRequest struct {
	Ids []int64 `json:"ids" validate:"required,min=1,dive,gt=0"`
}
