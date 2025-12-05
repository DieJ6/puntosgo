package usecases

import (
	"encoding/json"

	"github.com/tuusuario/puntosgo/internal/rabbit"
)

type DevolverCompraUC struct {
	Publisher rabbit.Publisher
}

func (uc *DevolverCompraUC) Execute(orderID string, puntos int, faltante float64) error {
	body, _ := json.Marshal(map[string]interface{}{
		"order_id":         orderID,
		"puntos_aplicados": puntos,
		"faltante_pesos":   faltante,
	})

	return uc.Publisher.Publish("informacion_compra", body)
}
