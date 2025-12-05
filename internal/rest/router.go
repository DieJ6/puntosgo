package rest

import (
    "github.com/go-chi/chi/v5"
    "github.com/tuusuario/puntosgo/internal/di"
)

func Router(r chi.Router, inj *di.Injector) {

    // Category
    categoryHandlers := CategoryHandlers{Inj: inj}
    r.Post("/v1/puntosPorCompra/categoria", categoryHandlers.CreateCategory)
    r.Post("/v1/puntosPorCompra/categoria/addArticle", categoryHandlers.AddArticle)
    r.Put("/v1/puntosPorCompra/categoria/delArticle", categoryHandlers.RemoveArticle)

    // Points
    pointsHandlers := PointsHandlers{Inj: inj}
    r.Get("/v1/puntosPorCompra/misPuntos", pointsHandlers.GetPoints)

    // Movements
    movHandlers := MovHandlers{Inj: inj}
    r.Get("/v1/puntosPorCompra/misMovimientos", movHandlers.GetMovements)

    // Async request trigger
    compraHandlers := CompraHandlers{Inj: inj}
    r.Post("/v1/puntosPorCompra/consultaCompra", compraHandlers.TriggerConsultaCompra)
}
