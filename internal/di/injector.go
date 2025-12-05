package di

import (
    "github.com/nmarsollier/commongo/log"
    "github.com/nmarsollier/commongo/db"
    "github.com/streadway/amqp"

    "github.com/tuusuario/puntosgo/internal/env"

    "github.com/tuusuario/puntosgo/internal/category"
    "github.com/tuusuario/puntosgo/internal/equivalencia"
    "github.com/tuusuario/puntosgo/internal/movimiento"
    "github.com/tuusuario/puntosgo/internal/saldo"
    "github.com/tuusuario/puntosgo/internal/rabbit"
)

type Injector struct {
    Log       log.LogRusEntry
    DB        db.Database
    Rabbit    *amqp.Connection

    CategoryRepo category.CategoryRepository
    EquivRepo    equivalencia.EquivalenciaRepository
    MovRepo      movimiento.MovimientoRepository
    SaldoRepo    saldo.SaldoRepository

    CategorySrv category.Service
    EquivSrv    equivalencia.Service
    MovSrv      movimiento.Service
    SaldoSrv    saldo.Service
}

var injector *Injector

func Initialize() *Injector {
    if injector != nil {
        return injector
    }

    logger := log.New()

    mongo, err := db.NewMongoDatabase(env.Get().MongoURL)
    if err != nil {
        logger.Fatal(err)
    }

    conn, err := amqp.Dial(env.Get().RabbitURL)
    if err != nil {
        logger.Fatal(err)
    }

    // declaramos exchange + queue
    rabbit.Setup(conn)

    // repos
    catCol := mongo.Collection("categorias")
    eqCol  := mongo.Collection("equivalencias")
    mvCol  := mongo.Collection("movimientos")
    sldCol := mongo.Collection("saldos")

    catRepo := category.NewRepository(logger, catCol)
    eqRepo  := equivalencia.NewRepository(logger, eqCol)
    mvRepo  := movimiento.NewRepository(logger, mvCol)
    sldRepo := saldo.NewRepository(logger, sldCol)

    // services
    catSrv := category.NewService(catRepo)
    eqSrv  := equivalencia.NewService(eqRepo)
    mvSrv  := movimiento.NewService(mvRepo)
    sldSrv := saldo.NewService(sldRepo)

    injector = &Injector{
        Log:       logger,
        DB:        mongo,
        Rabbit:    conn,

        CategoryRepo: catRepo,
        EquivRepo:    eqRepo,
        MovRepo:      mvRepo,
        SaldoRepo:    sldRepo,

        CategorySrv: catSrv,
        EquivSrv:    eqSrv,
        MovSrv:      mvSrv,
        SaldoSrv:    sldSrv,
    }

    return injector
}

func Get() *Injector {
    return injector
}
