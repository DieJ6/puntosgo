package rabbit

import "github.com/streadway/amqp"

const (
	ExchangeName      = "puntos_exchange"
	QueueName         = "puntos_queue"
	ResultQueueName   = "puntos_result_queue"

	RkConsultaCompra      = "consulta_compra"
	RkFaltanteCompra      = "faltante_compra"
	RkInformacionCompra   = "informacion_compra"
	RkRegistrarCompra     = "registrar_compra"
)

// Configura el exchange y las colas
func Setup(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Exchange Direct
	if err := ch.ExchangeDeclare(
		ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Cola principal (donde escucha puntosgo)
	if _, err := ch.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Cola de resultados (para informacion_compra)
	if _, err := ch.QueueDeclare(
		ResultQueueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Binds de la cola principal
	for _, key := range []string{
		RkConsultaCompra,
		RkFaltanteCompra,
		RkRegistrarCompra,
	} {
		if err := ch.QueueBind(QueueName, key, ExchangeName, false, nil); err != nil {
			return err
		}
	}

	// Bind de resultados
	if err := ch.QueueBind(ResultQueueName, RkInformacionCompra, ExchangeName, false, nil); err != nil {
		return err
	}

	return nil
}
