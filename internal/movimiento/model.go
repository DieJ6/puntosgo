package movimiento

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movimiento struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FechaCreacion time.Time          `bson:"fechaCreacion" json:"fechaCreacion"`
	Monto         int                `bson:"monto" json:"monto" validate:"required"`
	ForKIdUsuario primitive.ObjectID `bson:"forK_id_usuario" json:"forK_id_usuario" validate:"required"`
}
