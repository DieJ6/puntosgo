// internal/rabbit/consumer.go
package rabbit

import (
	"encoding/json"
	"fmt"

	"github.com/DieJ6/puntosgo/internal/usecases"
	"github.com/nmarsollier/commongo/log"
	"github.com/streadway/amqp"
)

const (
	RkConsultaCompra  = "consulta_compra"
	RkRegistrarCompra = "registrar_compra"
)

type Consumer struct {
	conn             *amqp.Connection
	log              log.LogRusEntry // lo dejamos, pero no lo usamos por ahora
	ProcesarCompraUC *usecases.ProcesarCompraUC
	RegistrarCompraUC *usecases.RegistrarCompraUC
}

func NewConsumer(
	conn *amqp.Connection,
	logger log.LogRusEntry,
	procesarUC *usecases.ProcesarCompraUC,
	registrarUC *usecases.RegistrarCompraUC,
) *Consumer {
	return &Consumer{
		conn:              conn,
		log:               logger,
		ProcesarCompraUC:  procesarUC,
		RegistrarCompraUC: registrarUC,
	}
}

func (c *Consumer) Start() {
	// ====== chequeos defensivos ======
	if c == nil {
		fmt.Println("rabbit consumer: instancia nil, no se inicia el consumer")
		return
	}
	if c.conn == nil {
		fmt.Println("rabbit consumer: conexión AMQP nil, no se inicia el consumer")
		return
	}
	if c.ProcesarCompraUC == nil {
		fmt.Println("rabbit consumer: ProcesarCompraUC es nil, no se inicia el consumer")
		return
	}
	if c.RegistrarCompraUC == nil {
		fmt.Println("rabbit consumer: RegistrarCompraUC es nil, no se inicia el consumer")
		return
	}

	// Recover general del consumer
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic en Rabbit Consumer.Start:", r)
		}
	}()

	ch, err := c.conn.Channel()
	if err != nil {
		fmt.Println("rabbit consumer: error al abrir canal:", err)
		return
	}
	defer ch.Close()

	// Declarar exchange
	if err := ch.ExchangeDeclare(
		ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		fmt.Println("rabbit consumer: error ExchangeDeclare:", err)
		return
	}

	// Declarar cola
	_, err = ch.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("rabbit consumer: error QueueDeclare:", err)
		return
	}

	// Bind cola ↔ routing key (consulta_compra)
	if err := ch.QueueBind(
		QueueName,
		RkConsultaCompra,
		ExchangeName,
		false,
		nil,
	); err != nil {
		fmt.Println("rabbit consumer: error QueueBind (consulta_compra):", err)
		return
	}

	// Bind cola ↔ routing key (registrar_compra)
	if err := ch.QueueBind(
		QueueName,
		RkRegistrarCompra,
		ExchangeName,
		false,
		nil,
	); err != nil {
		fmt.Println("rabbit consumer: error QueueBind (registrar_compra):", err)
		return
	}

	msgs, err := ch.Consume(
		QueueName,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("rabbit consumer: error Consume:", err)
		return
	}

	fmt.Println("Rabbit consumer de puntosgo iniciado y escuchando:", RkConsultaCompra, "y", RkRegistrarCompra)

	for msg := range msgs {
		// Aislamos el procesamiento de cada mensaje con un recover
		func(m amqp.Delivery) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("panic procesando mensaje:", r)
					_ = m.Nack(false, false)
				}
			}()

			switch m.RoutingKey {
			case RkConsultaCompra:
				if err := c.ProcesarCompraUC.Consume(m.Body); err != nil {
					fmt.Println("error en ProcesarCompraUC.Consume:", err)
					_ = m.Nack(false, false)
					return
				}
				_ = m.Ack(false)

			case RkRegistrarCompra:
				var input usecases.RegistrarCompraInput
				if err := json.Unmarshal(m.Body, &input); err != nil {
					fmt.Println("error json.Unmarshal RegistrarCompraInput:", err)
					_ = m.Nack(false, false)
					return
				}

				if err := c.RegistrarCompraUC.Ejecutar(input); err != nil {
					fmt.Println("error en RegistrarCompraUC.Ejecutar:", err)
					_ = m.Nack(false, false)
					return
				}

				_ = m.Ack(false)

			default:
				_ = m.Ack(false)
			}
		}(msg)
	}
}
