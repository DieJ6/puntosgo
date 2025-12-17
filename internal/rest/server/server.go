package server

import (
    "fmt"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    httpSwagger "github.com/swaggo/http-swagger"

    "github.com/DieJ6/puntosgo/internal/rest"
    "github.com/DieJ6/puntosgo/internal/env"
    "github.com/DieJ6/puntosgo/internal/di"
    "github.com/DieJ6/puntosgo/internal/rabbit"
    "github.com/DieJ6/puntosgo/internal/usecases"
)

func Start() {

    inj := di.Initialize()

    // ——— Consumer RabbitMQ (consulta_compra) ———
    procUC := &usecases.ProcesarCompraUC{
        CategorySrv: inj.CategorySrv,
        EquivSrv:    inj.EquivSrv,
        SaldoSrv:    inj.SaldoSrv,
        MvSrv:       inj.MvSrv,
        Publisher:   rabbit.NewPublisher(inj.Rabbit, inj.Log),
    }

    regUC := &usecases.RegistrarCompraUC{
        SaldoSrv: inj.SaldoSrv,
        MvSrv:    inj.MvSrv,
    }

    consumer := rabbit.NewConsumer(inj.Rabbit, inj.Log, procUC, regUC)
    go consumer.Start()

    //resConsumer := rabbit.NewResultConsumer(inj.Rabbit, inj.Log)
    //go resConsumer.Start()

    // ——— REST API ———
    r := chi.NewRouter()

    rest.Router(r, inj)

    r.Get("/swagger/*", httpSwagger.WrapHandler)

    port := env.Get().Port
    fmt.Println("PuntosGo escuchando en puerto", port)

    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
