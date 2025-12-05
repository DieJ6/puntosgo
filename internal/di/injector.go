package di

import (
	"fmt"
	"time"

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

// pequeño helper para reintentar la conexión a RabbitMQ
func dialRabbitWithRetry(url string, attempts int, delay time.Duration) *amqp.Connection {
	var conn *amqp.Connection
	var err error

	for i := 1; i <= attempts; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			return conn
		}

		fmt.Printf("RabbitMQ no disponible (intento %d/%d): %v\n", i, attempts, err)
		time.Sleep(delay)
	}

	// último intento, si falla, panic como antes
	conn, err = amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	return conn
}

func Initialize() *Injector {
	if injector != nil {
		return injector
	}

	// Logger "cero valor" (commongo/log no expone un constructor público simple)
	var logger log.LogRusEntry

	// Conexión a RabbitMQ con reintentos
	rabbitURL := env.Get().RabbitURL
	conn := dialRabbitWithRetry(rabbitURL, 10, 5*time.Second)

	// Declaramos exchange + queue según nuestra configuración
	if err := rabbit.Setup(conn); err != nil {
		panic(err)
	}

	// OJO: por ahora NO cableamos Mongo real.
	// Usamos un db.Collection nil tipado solo para que compile.
	var nilCollection db.Collection = nil

	// Repos (todavía apuntando a un collection nil)
	catRepo := category.NewRepository(logger, nilCollection)
	eqRepo := equivalencia.NewRepository(logger, nilCollection)
	mvRepo := movimiento.NewRepository(logger, nilCollection)
	sldRepo := saldo.NewRepository(logger, nilCollection)

	// Servicios
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
