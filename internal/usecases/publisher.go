package usecases

// Publisher es la abstracción que usan los casos de uso
// para publicar mensajes en Rabbit (o donde sea).
type Publisher interface {
	Publish(routingKey string, body []byte) error
}
