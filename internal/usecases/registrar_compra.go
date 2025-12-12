package usecases

import (
	"github.com/DieJ6/puntosgo/internal/movimiento"
	"github.com/DieJ6/puntosgo/internal/saldo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegistrarCompraUC struct {
	SaldoSrv saldo.Service
	MvSrv    movimiento.Service
}

// Input para registrar puntos de una compra
type RegistrarCompraInput struct {
	UserID string  `json:"user_id"`
	Monto  float64 `json:"monto"` // precio total de la compra
}

// Regla ejemplo: 10 pesos = 1 punto
const PESOS_POR_PUNTO = 10.0

func (uc *RegistrarCompraUC) Ejecutar(input RegistrarCompraInput) error {
	uid, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return err
	}

	// cálculo de puntos por compra
	puntos := int(input.Monto / PESOS_POR_PUNTO)

	// obtener o crear saldo
	s, err := uc.SaldoSrv.GetSaldoActual(uid)
	if err != nil {
		return err
	}
	if s == nil {
		s, err = uc.SaldoSrv.CrearSaldoInicial(uid)
		if err != nil {
			return err
		}
	}

	// sumar puntos
	s.Monto += puntos
	_, err = uc.SaldoSrv.ActualizarSaldo(s)
	if err != nil {
		return err
	}

	// registrar movimiento +puntos
	_, err = uc.MvSrv.Registrar(&movimiento.Movimiento{
		Monto:         puntos,
		ForKIdUsuario: uid,
	})

	return err
}
