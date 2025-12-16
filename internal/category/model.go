package category

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre             string             `bson:"nombre" json:"nombre"`
	ForKIdEquivalencia primitive.ObjectID `bson:"forK_id_equivalencia" json:"forK_id_equivalencia"`
	Prioridad          int                `bson:"prioridad" json:"prioridad"`
	Articulos          []string           `bson:"articulos" json:"articulos"`
	FechaCreacion      time.Time          `bson:"fechaCreacion" json:"fechaCreacion"`
}
