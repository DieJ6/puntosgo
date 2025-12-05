package usecases

import (
	"errors"

	"github.com/tuusuario/puntosgo/internal/category"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuitarArticuloUC struct {
	CategorySrv category.Service
}

type QuitarArticuloInput struct {
	IDCategoria primitive.ObjectID `json:"id_categoria"`
	IDArticulo  string             `json:"id_articulo"`
}

func (uc *QuitarArticuloUC) Execute(input QuitarArticuloInput) error {

	if input.IDCategoria.IsZero() {
		return errors.New("id de categoría inválido")
	}
	if input.IDArticulo == "" {
		return errors.New("id de artículo inválido")
	}

	return uc.CategorySrv.RemoveArticulo(input.IDCategoria, input.IDArticulo)
}
