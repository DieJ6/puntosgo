package usecases

import "encoding/json"

// Estructura que mandamos por Rabbit
type ResultadoCompra struct {
	OrderID         string  `json:"order_id"`
	PuntosAplicados int     `json:"puntos_aplicados"`
	FaltantePesos   float64 `json:"faltante_pesos"`
}

// Use case que solo se encarga de serializar y publicar el resultado
type DevolverResultadoCompraUC struct {
	Publisher Publisher
}

func (uc *DevolverResultadoCompraUC) Execute(res ResultadoCompra) error {
	body, err := json.Marshal(res)
	if err != nil {
		return err
	}

	// Ya no usamos rabbit.Publisher, sino la interface Publisher
	return uc.Publisher.Publish("informacion_compra", body)
}
