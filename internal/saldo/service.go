package saldo

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	CrearSaldoInicial(uid primitive.ObjectID) (*Saldo, error)
	GetSaldoActual(uid primitive.ObjectID) (*Saldo, error)
	ActualizarSaldo(s *Saldo) (*Saldo, error)
}

type service struct {
	repo SaldoRepository
}

func NewService(r SaldoRepository) Service {
	return &service{repo: r}
}

func (s *service) CrearSaldoInicial(uid primitive.ObjectID) (*Saldo, error) {
	if uid.IsZero() {
		return nil, errors.New("id de usuario inválido")
	}

	sld := &Saldo{
		ForKIdUsuario: uid,
		Monto:         0,
	}

	return s.repo.Insert(sld)
}

func (s *service) GetSaldoActual(uid primitive.ObjectID) (*Saldo, error) {
	if uid.IsZero() {
		return nil, errors.New("id de usuario inválido")
	}

	return s.repo.FindLatestByUsuario(uid)
}

func (s *service) ActualizarSaldo(sld *Saldo) (*Saldo, error) {
	return s.repo.Update(sld)
}
