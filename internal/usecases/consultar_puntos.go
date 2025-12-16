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

type ConsultarPuntosResult struct {
	Puntos            int       `json:"puntos"`
	FechaModificacion time.Time `json:"fechaModificacion"`
}

func (uc *ConsultarPuntosUC) Execute(uid primitive.ObjectID) (*ConsultarPuntosResult, error) {
	// 1) saldo más reciente
	s, err := uc.SaldoSrv.GetSaldoActual(uid)
	if err != nil {
		return nil, err
	}

	base := 0
	ref := time.Time{}

	if s != nil {
		base = s.Monto
		ref = s.FechaModificacion
		if ref.IsZero() {
			ref = s.FechaCreacion
		}
	}

	// 2) movimientos posteriores al saldo
	movs, err := uc.MvSrv.GetByUsuarioAfter(uid, ref)
	if err != nil {
		return nil, err
	}

	// 3) saldo + movimientos
	total := base
	for _, m := range movs {
		total += m.Monto
	}

	if ref.IsZero() {
		ref = time.Now()
	}

	return &ConsultarPuntosResult{
		Puntos:            total,
		FechaModificacion: ref,
	}, nil
}
