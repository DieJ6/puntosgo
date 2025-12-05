package rabbit

import (
	"github.com/streadway/amqp"
)

const (
	ExchangeName = "puntos_exchange"
	QueueName    = "puntos_queue"
)

// Configura el exchange y la queue
func Setup(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declarar exchange tipo Direct
	if err := ch.ExchangeDeclare(
		ExchangeName,
		"direct",
		true,  // durable
		false, // auto-delete
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Declarar queue
	_, err = ch.QueueDeclare(
		QueueName,
		true,  // durable
		false, 
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Bindings
	routingKeys := []string{
		"consulta_compra",
		"faltante_compra",
		"informacion_compra", // opcional, por si algún servicio escucha
	}

	for _, key := range routingKeys {
		if err := ch.QueueBind(
			QueueName,
			key,
			ExchangeName,
			false,
			nil,
		); err != nil {
			return err
		}
	}

	return nil
}
