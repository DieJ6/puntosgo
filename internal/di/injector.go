package di

import (
	"github.com/streadway/amqp"

	"github.com/nmarsollier/commongo/db"
	"github.com/nmarsollier/commongo/log"

	"github.com/DieJ6/puntosgo/internal/category"
	"github.com/DieJ6/puntosgo/internal/env"
	"github.com/DieJ6/puntosgo/internal/equivalencia"
	"github.com/DieJ6/puntosgo/internal/movimiento"
	"github.com/DieJ6/puntosgo/internal/rabbit"
	"github.com/DieJ6/puntosgo/internal/saldo"
)

type Injector struct {
	Log    log.LogRusEntry
	Rabbit *amqp.Connection

	CategoryRepo category.CategoryRepository
	EquivRepo    equivalencia.EquivalenciaRepository
	MvRepo       movimiento.MovimientoRepository
	SaldoRepo    saldo.SaldoRepository

	CategorySrv category.Service
	EquivSrv    equivalencia.Service
	MvSrv       movimiento.Service
	SaldoSrv    saldo.Service
}

var injector *Injector

func Initialize() *Injector {
	if injector != nil {
		return injector
	}

	// Logger "cero valor": commongo/log no expone log.New(),
	// pero el tipo LogRusEntry se puede usar en blanco.
	var logger log.LogRusEntry

	// Conexión a RabbitMQ
	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		panic(err)
	}

	// Declaramos exchange + queue según nuestra config de Rabbit
	if err := rabbit.Setup(conn); err != nil {
		panic(err)
	}

	// Por ahora NO cableamos Mongo real:
	// usamos nil como db.Collection tipado, solo para que compile.
	var nilCollection db.Collection = nil

	catRepo := category.NewRepository(logger, nilCollection)
	eqRepo := equivalencia.NewRepository(logger, nilCollection)
	mvRepo := movimiento.NewRepository(logger, nilCollection)
	sldRepo := saldo.NewRepository(logger, nilCollection)

	catSrv := category.NewService(catRepo)
	eqSrv := equivalencia.NewService(eqRepo)
	mvSrv := movimiento.NewService(mvRepo)
	sldSrv := saldo.NewService(sldRepo)

	injector = &Injector{
		Log:    logger,
		Rabbit: conn,

		CategoryRepo: catRepo,
		EquivRepo:    eqRepo,
		MvRepo:       mvRepo,
		SaldoRepo:    sldRepo,

		CategorySrv: catSrv,
		EquivSrv:    eqSrv,
		MvSrv:       mvSrv,
		SaldoSrv:    sldSrv,
	}

	return injector
}

func Get() *Injector {
	return injector
}
