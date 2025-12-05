package rest

import (
	"encoding/json"
	"net/http"

	"github.com/DieJ6/puntosgo/internal/di"
	"github.com/DieJ6/puntosgo/internal/rabbit"
)

type CompraHandlers struct {
	Inj *di.Injector
}

func (h CompraHandlers) TriggerConsultaCompra(w http.ResponseWriter, r *http.Request) {

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "body inválido", 400)
		return
	}

	msg, _ := json.Marshal(body)

	pub := rabbit.NewPublisher(h.Inj.Rabbit, h.Inj.Log)

	err := pub.Publish("consulta_compra", msg)
	if err != nil {
		http.Error(w, "error enviando evento", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Consulta de compra enviada",
	})
}
