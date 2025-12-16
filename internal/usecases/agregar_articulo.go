package usecases

import (
	"errors"

	"github.com/DieJ6/puntosgo/internal/category"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrArticuloYaAsignado = errors.New("Artículo ya asignado a otra categoría")

type AgregarArticuloUC struct {
	CategorySrv category.Service
}

type AgregarArticuloInput struct {
	IDCategoria primitive.ObjectID
	IDArticulo  string
}

func (uc *AgregarArticuloUC) Execute(input AgregarArticuloInput) error {

	if input.IDArticulo == "" {
		return errors.New("Id de artículo requerido")
	}

	// 1️⃣ Buscar si el artículo ya pertenece a alguna categoría
	catExistente, err := uc.CategorySrv.FindByArticulo(input.IDArticulo)
	if err == nil && catExistente != nil {

		// Si es la MISMA categoría, no hacemos nada (idempotente)
		if catExistente.ID == input.IDCategoria {
			return nil
		}

		// Si es OTRA categoría → conflicto
		return ErrArticuloYaAsignado
	}

	// 2️⃣ Agregar a la categoría solicitada
	return uc.CategorySrv.AddArticulo(input.IDCategoria, input.IDArticulo)
}
