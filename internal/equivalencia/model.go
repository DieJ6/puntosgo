package equivalencia

import "go.mongodb.org/mongo-driver/bson/primitive"

type Equivalencia struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Puntos int                `bson:"puntos" json:"puntos" validate:"required,min=1"`
	Pesos  int                `bson:"pesos" json:"pesos" validate:"required,min=1"`
}
