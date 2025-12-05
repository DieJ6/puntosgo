package equivalencia

import "go.mongodb.org/mongo-driver/bson/primitive"

type Service interface {
	GetByID(id primitive.ObjectID) (*Equivalencia, error)
}

type service struct {
	repo EquivalenciaRepository
}

func NewService(r EquivalenciaRepository) Service {
	return &service{repo: r}
}

func (s *service) GetByID(id primitive.ObjectID) (*Equivalencia, error) {
	return s.repo.FindByID(id)
}
