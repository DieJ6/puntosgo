package usecases

import (
	"github.com/tuusuario/puntosgo/internal/movimiento"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsultarMovimientosUC struct {
	MvSrv movimiento.Service
}

func (uc *ConsultarMovimientosUC) Execute(userID primitive.ObjectID) ([]*movimiento.Movimiento, error) {
	return uc.MvSrv.GetByUsuario(userID)
}
