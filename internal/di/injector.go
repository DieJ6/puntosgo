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

	"go.mongodb.org/mongo-driver/mongo"
)

type Injector struct {
	Log    log.LogRusEntry
	Rabbit *amqp.Connection

	// opcional: útil para debug
	MongoDB *mongo.Database

	CategoryRepo category.CategoryRepository
	EquivRepo    equivalencia.EquivalenciaRepository
	MvRepo       movimiento.MovimientoRepository
	SaldoRepo    saldo.SaldoRepository

	CategorySrv category.Service
	EquivSrv    equivalencia.Service
	MvSrv       movimiento.Service
	SaldoSrv    saldo.Service

	AuthURL string
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

	// Logger "cero valor"
	var logger log.LogRusEntry

	cfg := env.Get()

	// =====================
	// RabbitMQ
	// =====================
	conn := dialRabbitWithRetry(cfg.RabbitURL, 10, 5*time.Second)

	// Declaramos exchange + queue según nuestra configuración
	if err := rabbit.Setup(conn); err != nil {
		panic(err)
	}

	// =====================
	// MongoDB (REAL)
	// =====================
	const mongoDBName = "puntos"

	mongoDB, err := db.NewDatabase(cfg.MongoURL, mongoDBName)
	if err != nil {
		panic(err)
	}

	onError := func(e error) {
		// no usamos logger acá porque es "cero valor"; imprimimos para debug
		fmt.Println("mongo error:", e)
	}

	// Colecciones (reales)
	// Si querés índices, pasalos al final (indexes ...string)
	catCollection, err := db.NewCollection(logger, mongoDB, "categories", onError, "prioridad")
	if err != nil {
		panic(err)
	}

	eqCollection, err := db.NewCollection(logger, mongoDB, "equivalencias", onError)
	if err != nil {
		panic(err)
	}

	mvCollection, err := db.NewCollection(logger, mongoDB, "movimientos", onError, "ForKIdUsuario")
	if err != nil {
		panic(err)
	}

	sldCollection, err := db.NewCollection(logger, mongoDB, "saldos", onError, "ForKIdUsuario")
	if err != nil {
		panic(err)
	}

	// Repos
	catRepo := category.NewRepository(logger, catCollection)
	eqRepo := equivalencia.NewRepository(logger, eqCollection)
	mvRepo := movimiento.NewRepository(logger, mvCollection)
	sldRepo := saldo.NewRepository(logger, sldCollection)

	// Servicios
	catSrv := category.NewService(catRepo)
	eqSrv := equivalencia.NewService(eqRepo)
	mvSrv := movimiento.NewService(mvRepo)
	sldSrv := saldo.NewService(sldRepo)

	injector = &Injector{
		Log:     logger,
		Rabbit:  conn,
		MongoDB: mongoDB,

		CategoryRepo: catRepo,
		EquivRepo:    eqRepo,
		MvRepo:       mvRepo,
		SaldoRepo:    sldRepo,

		CategorySrv: catSrv,
		EquivSrv:    eqSrv,
		MvSrv:       mvSrv,
		SaldoSrv:    sldSrv,

		AuthURL: cfg.AuthURL,
	}

	return injector
}

func Get() *Injector {
	return injector
}
