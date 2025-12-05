package usecases

import (
	"encoding/json"
	"sort"

	"github.com/streadway/amqp"

	"github.com/tuusuario/puntosgo/internal/category"
	"github.com/tuusuario/puntosgo/internal/equivalencia"
	"github.com/tuusuario/puntosgo/internal/saldo"
	"github.com/tuusuario/puntosgo/internal/movimiento"
	"github.com/tuusuario/puntosgo/internal/rabbit"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcesarCompraUC struct {
	CategorySrv  category.Service
	EquivSrv     equivalencia.Service
	SaldoSrv     saldo.Service
	MvSrv        movimiento.Service
	Publisher    rabbit.Publisher
}

type ConsultaCompraInput struct {
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	Productos []struct {
		ID    string  `json:"id_producto"`
		Precio float64 `json:"precio"`
	} `json:"productos"`
}

type ResultadoCompra struct {
	OrderID        string  `json:"order_id"`
	PuntosAplicados int     `json:"puntos_aplicados"`
	FaltantePesos   float64 `json:"faltante_pesos"`
}

func (uc *ProcesarCompraUC) Consume(msg amqp.Delivery) error {
	var input ConsultaCompraInput
	if err := json.Unmarshal(msg.Body, &input); err != nil {
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
			// si no hay categoría → prioridad infinita
			evals = append(evals, ProductoEval{p.ID, p.Precio, 999, nil})
			continue
		}
		eq, _ := uc.EquivSrv.GetByID(cat.ForKIdEquivalencia)
		evals = append(evals, ProductoEval{p.ID, p.Precio, cat.Prioridad, eq})
	}

	// ordenar por prioridad ascendente
	sort.Slice(evals, func(i, j int) bool {
		return evals[i].Prioridad < evals[j].Prioridad
	})

	// aplicar puntos por producto
	for _, prod := range evals {
		if prod.Equiv == nil {
			faltanteTotal += prod.Precio
			continue
		}

		// relación: tantos puntos → tantos pesos
		valorPorPunto := float64(prod.Equiv.Pesos) / float64(prod.Equiv.Puntos)

		maxDescuento := float64(puntosDisponibles) * valorPorPunto

		if maxDescuento >= prod.Precio {
			// cubrir totalidad del producto
			puntosUsados := int(prod.Precio / valorPorPunto)
			puntosDisponibles -= puntosUsados
			puntosTotalesAplicados += puntosUsados
		} else {
			// cubrir parcialmente
			puntosUsados := puntosDisponibles
			puntosTotalesAplicados += puntosUsados
			puntosDisponibles = 0

			faltanteTotal += prod.Precio - maxDescuento
		}

		if puntosDisponibles <= 0 {
			// todos los productos restantes se suman al faltante
			for _, rest := range evals {
				if rest.ID == prod.ID {
					continue
				}
				faltanteTotal += rest.Precio
			}
			break
		}
	}

	// actualizar saldo
	s.Monto = puntosDisponibles
	uc.SaldoSrv.ActualizarSaldo(s)

	// registrar movimiento
	_, _ = uc.MvSrv.Registrar(&movimiento.Movimiento{
		Monto:         -puntosTotalesAplicados,
		ForKIdUsuario: uid,
	})

	// enviar respuesta
	result := ResultadoCompra{
		OrderID:        input.OrderID,
		PuntosAplicados: puntosTotalesAplicados,
		FaltantePesos:   faltanteTotal,
	}

	body, _ := json.Marshal(result)
	return uc.Publisher.Publish("informacion_compra", body)
}
