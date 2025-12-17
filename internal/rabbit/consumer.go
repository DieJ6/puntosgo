package rabbit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/nmarsollier/commongo/log"
	"github.com/streadway/amqp"

	"github.com/DieJ6/puntosgo/internal/usecases"
)

type Consumer struct {
	conn             *amqp.Connection
	log              log.LogRusEntry // lo dejamos pero NO lo usamos (puede ser nil)
	ProcesarCompraUC  *usecases.ProcesarCompraUC
	RegistrarCompraUC *usecases.RegistrarCompraUC
}

func NewConsumer(
	conn *amqp.Connection,
	logger log.LogRusEntry,
	procesarUC *usecases.ProcesarCompraUC,
	registrarUC *usecases.RegistrarCompraUC,
) *Consumer {
	return &Consumer{
		conn:             conn,
		log:              logger,
		ProcesarCompraUC:  procesarUC,
		RegistrarCompraUC: registrarUC,
	}
}

func (c *Consumer) Start() {
	if c == nil || c.conn == nil {
		fmt.Println("rabbit consumer: instancia o conexión nil, no se inicia")
		return
	}
	if c.ProcesarCompraUC == nil || c.RegistrarCompraUC == nil {
		fmt.Println("rabbit consumer: UCs nil, no se inicia")
		return
	}

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

	// Exchange + Queue (por si rabbit.Setup no lo dejó listo)
	if err := ch.ExchangeDeclare(ExchangeName, "direct", true, false, false, false, nil); err != nil {
		fmt.Println("rabbit consumer: error ExchangeDeclare:", err)
		return
	}
	if _, err := ch.QueueDeclare(QueueName, true, false, false, false, nil); err != nil {
		fmt.Println("rabbit consumer: error QueueDeclare:", err)
		return
	}

	// Bind ambas routing keys (UNA SOLA VEZ)
	for _, rk := range []string{RkConsultaCompra, RkRegistrarCompra} {
		if err := ch.QueueBind(QueueName, rk, ExchangeName, false, nil); err != nil {
			fmt.Println("rabbit consumer: error QueueBind:", rk, err)
			return
		}
	}

	msgs, err := ch.Consume(QueueName, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println("rabbit consumer: error Consume:", err)
		return
	}

	fmt.Println("Rabbit consumer de puntosgo iniciado y escuchando: consulta_compra y registrar_compra")

	for msg := range msgs {
		func(m amqp.Delivery) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("panic procesando mensaje:", r)
					_ = m.Nack(false, false) // descartar (no requeue)
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
				input, err := decodeRegistrarCompraInput(m.Body)
				if err != nil {
					fmt.Println("error decode registrar_compra:", err)
					_ = m.Nack(false, false)
					return
				}

				// Validación mínima defensiva
				if strings.TrimSpace(input.UserID) == "" {
					fmt.Println("registrar_compra: user_id vacío")
					_ = m.Nack(false, false)
					return
				}
				if input.Monto <= 0 {
					fmt.Println("registrar_compra: monto inválido")
					_ = m.Nack(false, false)
					return
				}

				if err := c.RegistrarCompraUC.Ejecutar(*input); err != nil {
					fmt.Println("error RegistrarCompraUC.Ejecutar:", err)
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

// Acepta body como JSON objeto o como JSON string que contiene un objeto.
// También tolera que el string interno tenga saltos de línea reales (caso del error '\n' in string literal).
func decodeRegistrarCompraInput(body []byte) (*usecases.RegistrarCompraInput, error) {
	b := bytes.TrimSpace(body)
	if len(b) == 0 {
		return nil, errors.New("body vacío")
	}

	// 1) Intento directo: { "user_id": "...", "monto": 1200 }
	var in usecases.RegistrarCompraInput
	if err := json.Unmarshal(b, &in); err == nil {
		return &in, nil
	}

	// 2) Fallback: el body viene como string: "{\"user_id\":\"...\",\"monto\":1200}"
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}

	// Si dentro del string vinieron saltos de línea reales, los limpiamos
	// (esto es lo que te rompía con: invalid character '\n' in string literal)
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimSpace(s)

	if s == "" {
		return nil, errors.New("payload string vacío")
	}

	if err := json.Unmarshal([]byte(s), &in); err != nil {
		return nil, err
	}
	return &in, nil
}
