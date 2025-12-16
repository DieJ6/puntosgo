package usecases

import (
	"errors"

	"github.com/DieJ6/puntosgo/internal/movimiento"
	"github.com/DieJ6/puntosgo/internal/saldo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegistrarCompraUC struct {
	SaldoSrv saldo.Service
	MvSrv    movimiento.Service
}

type RegistrarCompraInput struct {
	UserID string  `json:"user_id"`
	Monto  float64 `json:"monto"`
}

const PESOS_POR_PUNTO = 10.0

func (uc *RegistrarCompraUC) Ejecutar(input RegistrarCompraInput) error {
	if uc == nil || uc.SaldoSrv == nil || uc.MvSrv == nil {
		return errors.New("usecase no inicializado")
	}

	uid, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return errors.New("user_id inválido")
	}

	puntos := int(input.Monto / PESOS_POR_PUNTO)
	if puntos <= 0 {
		// si querés, podés permitir 0 y no registrar nada
		return nil
	}

	// obtener saldo actual
	s, err := uc.SaldoSrv.GetSaldoActual(uid)
	if err != nil {
		// si no encontrás saldo, lo creás; cualquier otro error se propaga
		s = nil
	}

	if s == nil {
		s, err = uc.SaldoSrv.CrearSaldoInicial(uid)
		if err != nil {
			return err
		}
	}

	s.Monto += puntos
	if _, err := uc.SaldoSrv.ActualizarSaldo(s); err != nil {
		return err
	}

	_, err = uc.MvSrv.Registrar(&movimiento.Movimiento{
		Monto:         puntos,
		ForKIdUsuario: uid,
	})
	return err
}
