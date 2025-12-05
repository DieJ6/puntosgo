package category

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Create(cat *Category) (*Category, error)
	AddArticulo(catID primitive.ObjectID, productID string) error
	RemoveArticulo(catID primitive.ObjectID, productID string) error
	GetByID(id primitive.ObjectID) (*Category, error)
	FindByArticulo(productID string) (*Category, error)
}

type service struct {
	repo CategoryRepository
}

func NewService(r CategoryRepository) Service {
	return &service{repo: r}
}

func (s *service) Create(cat *Category) (*Category, error) {
	return s.repo.Insert(cat)
}

func (s *service) AddArticulo(catID primitive.ObjectID, productID string) error {
	if productID == "" {
		return errors.New("id de producto inválido")
	}
	return s.repo.AddArticulo(catID, productID)
}

func (s *service) RemoveArticulo(catID primitive.ObjectID, productID string) error {
	if productID == "" {
		return errors.New("id de producto inválido")
	}
	return s.repo.RemoveArticulo(catID, productID)
}

func (s *service) GetByID(id primitive.ObjectID) (*Category, error) {
	return s.repo.FindByID(id)
}

func (s *service) FindByArticulo(productID string) (*Category, error) {
	return s.repo.FindByArticuloID(productID)
}
