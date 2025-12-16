package saldo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Saldo struct {
    ID                primitive.ObjectID `bson:"_id" json:"id"`
    FechaCreacion     time.Time          `bson:"fechaCreacion" json:"fechaCreacion"`
    FechaModificacion time.Time          `bson:"fechaModificacion" json:"fechaModificacion"`
    Monto             int                `bson:"monto" json:"monto"`
    ForKIdUsuario     primitive.ObjectID `bson:"forK_id_usuario" json:"forK_id_usuario"`
}
