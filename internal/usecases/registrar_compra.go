package usecases

import (
	"errors"

	"github.com/DieJ6/puntosgo/internal/movimiento"
	"github.com/DieJ6/puntosgo/internal/saldo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		return nil
	}

	// 1) Obtener saldo actual (si no existe, crearlo)
	s, err := uc.SaldoSrv.GetSaldoActual(uid)
	if err != nil {
		// Sólo tratamos como "no existe" el ErrNoDocuments
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
		s = nil
	}
	if s == nil {
		s, err = uc.SaldoSrv.CrearSaldoInicial(uid)
		if err != nil {
			return err
		}
	}

	// 2) Registrar movimiento (primero)
	if _, err := uc.MvSrv.Registrar(&movimiento.Movimiento{
		Monto:         puntos,
		ForKIdUsuario: uid,
	}); err != nil {
		return err
	}

	// 3) Actualizar saldo (después)
	s.Monto += puntos
	_, err = uc.SaldoSrv.ActualizarSaldo(s)
	return err
}
