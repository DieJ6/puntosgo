package usecases

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/DieJ6/puntosgo/internal/movimiento"
)

type ConsultarMovimientosUC struct {
	MvSrv movimiento.Service
}

func (uc *ConsultarMovimientosUC) Execute(uid primitive.ObjectID) ([]*movimiento.Movimiento, error) {
	movs, err := uc.MvSrv.GetByUsuario(uid)
	if err != nil {
		return nil, err
	}
	if movs == nil {
		return []*movimiento.Movimiento{}, nil
	}
	return movs, nil
}
