package usecases

import (
	"errors"
	"time"

	"github.com/tuusuario/puntosgo/internal/category"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CrearCategoriaUC struct {
	CategorySrv category.Service
}

type CrearCategoriaInput struct {
	Nombre             string             `json:"nombre"`
	ForKIdEquivalencia primitive.ObjectID `json:"forK_id_equivalencia"`
	Prioridad          int                `json:"prioridad"`
	Articulos          []string           `json:"articulos"`
}

func (uc *CrearCategoriaUC) Execute(input CrearCategoriaInput) (*category.Category, error) {

	if input.Nombre == "" {
		return nil, errors.New("nombre requerido")
	}
	if input.ForKIdEquivalencia.IsZero() {
		return nil, errors.New("id equivalencia requerido")
	}
	if input.Prioridad < 1 {
		return nil, errors.New("prioridad inválida")
	}

	cat := &category.Category{
		Nombre:             input.Nombre,
		ForKIdEquivalencia: input.ForKIdEquivalencia,
		Prioridad:          input.Prioridad,
		Articulos:          input.Articulos,
		FechaCreacion:      time.Now(),
	}

	return uc.CategorySrv.Create(cat)
}
