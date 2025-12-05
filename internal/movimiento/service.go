package movimiento

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Registrar(m *Movimiento) (*Movimiento, error)
	GetByUsuario(uid primitive.ObjectID) ([]*Movimiento, error)
}

type service struct {
	repo MovimientoRepository
}

func NewService(r MovimientoRepository) Service {
	return &service{repo: r}
}

func (s *service) Registrar(m *Movimiento) (*Movimiento, error) {
	return s.repo.Insert(m)
}

func (s *service) GetByUsuario(uid primitive.ObjectID) ([]*Movimiento, error) {
	return s.repo.FindByUsuario(uid)
}
