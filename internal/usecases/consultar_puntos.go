package usecases

import (
	"time"

	"github.com/DieJ6/puntosgo/internal/movimiento"
	"github.com/DieJ6/puntosgo/internal/saldo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsultarPuntosUC struct {
	SaldoSrv saldo.Service
	MvSrv    movimiento.Service
}

type ConsultarPuntosOutput struct {
	Puntos            int       `json:"puntos"`
	FechaModificacion time.Time `json:"fechaModificacion"`
}

func (uc *ConsultarPuntosUC) Execute(userID primitive.ObjectID) (*ConsultarPuntosOutput, error) {

	s, err := uc.SaldoSrv.GetSaldoActual(userID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		// usuario sin saldo → iniciar en 0
		s, err = uc.SaldoSrv.CrearSaldoInicial(userID)
		if err != nil {
			return nil, err
		}
	}

	movs, err := uc.MvSrv.GetByUsuario(userID)
	if err != nil {
		return nil, err
	}

	total := s.Monto

	for _, mv := range movs {
		if mv.FechaCreacion.After(s.FechaModificacion) {
			total += mv.Monto
		}
	}

	return &ConsultarPuntosOutput{
		Puntos:            total,
		FechaModificacion: s.FechaModificacion,
	}, nil
}
