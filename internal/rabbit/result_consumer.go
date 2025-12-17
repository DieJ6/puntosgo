package rabbit

import (
	"encoding/json"
	"fmt"

	"github.com/nmarsollier/commongo/log"
	"github.com/streadway/amqp"
)

type ResultConsumer struct {
	conn *amqp.Connection
	log  log.LogRusEntry
}

func NewResultConsumer(conn *amqp.Connection, logger log.LogRusEntry) *ResultConsumer {
	return &ResultConsumer{conn: conn, log: logger}
}

func (c *ResultConsumer) Start() {
	if c == nil || c.conn == nil {
		fmt.Println("result consumer: instancia o conexión nil, no se inicia")
		return
	}

	ch, err := c.conn.Channel()
	if err != nil {
		fmt.Println("result consumer: error al abrir canal:", err)
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume(ResultQueueName, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println("result consumer: error Consume:", err)
		return
	}

	fmt.Println("ResultConsumer iniciado y escuchando informacion_compra en puntos_result_queue")

	for msg := range msgs {
		func(m amqp.Delivery) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("panic procesando resultado:", r)
					_ = m.Nack(false, false)
				}
			}()

			// Log “bonito”
			var pretty map[string]any
			if err := json.Unmarshal(m.Body, &pretty); err == nil {
				b, _ := json.MarshalIndent(pretty, "", "  ")
				fmt.Println("RESULTADO COMPRA (informacion_compra):\n" + string(b))
			} else {
				fmt.Println("RESULTADO COMPRA (raw):", string(m.Body))
			}

			_ = m.Ack(false)
		}(msg)
	}
}
