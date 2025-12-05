package usecases

import (
	"encoding/json"
	"sort"

	"github.com/DieJ6/puntosgo/internal/category"
	"github.com/DieJ6/puntosgo/internal/equivalencia"
	"github.com/DieJ6/puntosgo/internal/movimiento"
	"github.com/DieJ6/puntosgo/internal/saldo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcesarCompraUC struct {
	CategorySrv category.Service
	EquivSrv    equivalencia.Service
	SaldoSrv    saldo.Service
	MvSrv       movimiento.Service
	Publisher   Publisher // interfaz, implementada por rabbit.Publisher
}

type ConsultaCompraInput struct {
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	Productos []struct {
		ID     string  `json:"id_producto"`
		Precio float64 `json:"precio"`
	} `json:"productos"`
}

// OJO: ResultadoCompra **NO** se declara acá,
// ya existe en devolver_resultado_compra.go
// y al estar en el mismo package usecases lo podemos usar directo.

func (uc *ProcesarCompraUC) Consume(body []byte) error {
	var input ConsultaCompraInput
	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	uid, _ := primitive.ObjectIDFromHex(input.UserID)

	// obtener saldo actual
	s, err := uc.SaldoSrv.GetSaldoActual(uid)
	if err != nil {
		return err
	}
	if s == nil {
		s, _ = uc.SaldoSrv.CrearSaldoInicial(uid)
	}

	puntosDisponibles := s.Monto
	puntosTotalesAplicados := 0
	faltanteTotal := 0.0

	// ordenar productos por prioridad
	type ProductoEval struct {
		ID        string
		Precio    float64
		Prioridad int
		Equiv     *equivalencia.Equivalencia
	}

	var evals []ProductoEval

	for _, p := range input.Productos {
		cat, err := uc.CategorySrv.FindByArticulo(p.ID)
		if err != nil || cat == nil {
			evals = append(evals, ProductoEval{
				ID:        p.ID,
				Precio:    p.Precio,
				Prioridad: 999,
				Equiv:     nil,
			})
			continue
		}
		eq, _ := uc.EquivSrv.GetByID(cat.ForKIdEquivalencia)
		evals = append(evals, ProductoEval{
			ID:        p.ID,
			Precio:    p.Precio,
			Prioridad: cat.Prioridad,
			Equiv:     eq,
		})
	}

	sort.Slice(evals, func(i, j int) bool {
		return evals[i].Prioridad < evals[j].Prioridad
	})

	for i := 0; i < len(evals); i++ {
		prod := evals[i]

		if prod.Equiv == nil {
			faltanteTotal += prod.Precio
			continue
		}

		valorPorPunto := float64(prod.Equiv.Pesos) / float64(prod.Equiv.Puntos)
		maxDescuento := float64(puntosDisponibles) * valorPorPunto

		if maxDescuento >= prod.Precio {
			puntosUsados := int(prod.Precio / valorPorPunto)
			puntosDisponibles -= puntosUsados
			puntosTotalesAplicados += puntosUsados
		} else {
			puntosUsados := puntosDisponibles
			puntosTotalesAplicados += puntosUsados
			puntosDisponibles = 0
			faltanteTotal += prod.Precio - maxDescuento
		}

		if puntosDisponibles <= 0 {
			for j := i + 1; j < len(evals); j++ {
				faltanteTotal += evals[j].Precio
			}
			break
		}
	}

	// actualizar saldo
	s.Monto = puntosDisponibles
	uc.SaldoSrv.ActualizarSaldo(s)

	// registrar movimiento de salida de puntos
	_, _ = uc.MvSrv.Registrar(&movimiento.Movimiento{
		Monto:         -puntosTotalesAplicados,
		ForKIdUsuario: uid,
	})

	// usamos la struct ResultadoCompra definida en devolver_resultado_compra.go
	result := ResultadoCompra{
		OrderID:         input.OrderID,
		PuntosAplicados: puntosTotalesAplicados,
		FaltantePesos:   faltanteTotal,
	}

	responseBody, _ := json.Marshal(result)
	return uc.Publisher.Publish("informacion_compra", responseBody)
}
