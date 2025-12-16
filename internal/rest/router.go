package rest

import (
    "github.com/go-chi/chi/v5"
    "github.com/DieJ6/puntosgo/internal/di"
)

func Router(r chi.Router, inj *di.Injector) {
	// auth middleware global (carga usuario)
	r.Use(RequireAuth(inj))

	// Category (ADMIN)
	categoryHandlers := CategoryHandlers{Inj: inj}
	r.With(RequireAdmin).Post("/v1/puntosPorCompra/categoria", categoryHandlers.CreateCategory)
	r.With(RequireAdmin).Post("/v1/puntosPorCompra/categoria/addArticle", categoryHandlers.AddArticle)
	r.With(RequireAdmin).Put("/v1/puntosPorCompra/categoria/delArticle", categoryHandlers.RemoveArticle)

	// Points (USER)
	pointsHandlers := PointsHandlers{Inj: inj}
	r.With(RequireUser).Get("/v1/puntosPorCompra/misPuntos", pointsHandlers.GetPoints)

	// Movements (USER)
	movHandlers := MovHandlers{Inj: inj}
	r.With(RequireUser).Get("/v1/puntosPorCompra/misMovimientos", movHandlers.GetMovements)

	// Consulta compra (USER)
	compraHandlers := CompraHandlers{Inj: inj}
	r.With(RequireUser).Post("/v1/puntosPorCompra/consultaCompra", compraHandlers.TriggerConsultaCompra)
}
