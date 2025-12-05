package usecases

import (
	"errors"

	"github.com/DieJ6/puntosgo/internal/category"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgregarArticuloUC struct {
	CategorySrv category.Service
}

type AgregarArticuloInput struct {
	IDCategoria primitive.ObjectID `json:"id_categoria"`
	IDArticulo  string             `json:"id_articulo"`
}

func (uc *AgregarArticuloUC) Execute(input AgregarArticuloInput) error {

	if input.IDCategoria.IsZero() {
		return errors.New("id de categoría inválido")
	}
	if input.IDArticulo == "" {
		return errors.New("id de artículo inválido")
	}

	return uc.CategorySrv.AddArticulo(input.IDCategoria, input.IDArticulo)
}
