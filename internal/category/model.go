package category

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID                 primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Nombre             string               `bson:"nombre" json:"nombre" validate:"required,min=1,max=100"`
	ForKIdEquivalencia primitive.ObjectID   `bson:"forK_id_equivalencia" json:"forK_id_equivalencia" validate:"required"`
	Prioridad          int                  `bson:"prioridad" json:"prioridad" validate:"required,min=1"`
	Articulos          []string             `bson:"articulos" json:"articulos" validate:"dive,min=1"`
	FechaCreacion      time.Time            `bson:"fechaCreacion" json:"fechaCreacion"`
}
