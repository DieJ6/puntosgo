package saldo

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nmarsollier/commongo/db"
	"github.com/nmarsollier/commongo/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

type SaldoRepository interface {
	Insert(s *Saldo) (*Saldo, error)
	Update(s *Saldo) (*Saldo, error)
	FindLatestByUsuario(uid primitive.ObjectID) (*Saldo, error)
}

type repository struct {
	log        log.LogRusEntry
	collection db.Collection
}

func NewRepository(log log.LogRusEntry, c db.Collection) SaldoRepository {
	return &repository{log, c}
}

func (r *repository) Insert(s *Saldo) (*Saldo, error) {
	if err := validate.Struct(s); err != nil {
		r.log.Error(err)
		return nil, err
	}

	s.ID = primitive.NewObjectID()
	now := time.Now()
	s.FechaCreacion = now
	s.FechaModificacion = now

	_, err := r.collection.InsertOne(context.Background(), s)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	return s, nil
}

func (r *repository) Update(s *Saldo) (*Saldo, error) {
	if err := validate.Struct(s); err != nil {
		r.log.Error(err)
		return nil, err
	}

	s.FechaModificacion = time.Now()

	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": s.ID},
		bson.M{"$set": s},
	)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	return s, nil
}

func (r *repository) FindLatestByUsuario(uid primitive.ObjectID) (*Saldo, error) {
	var s Saldo

	opts := db.FindOptions().
		SetSort(bson.D{{"fechaModificacion", -1}}).
		SetLimit(1)

	cur, err := r.collection.Find(context.Background(),
		bson.M{"forK_id_usuario": uid},
		opts,
	)

	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	if cur.Next(context.Background()) {
		if err := cur.Decode(&s); err != nil {
			return nil, err
		}
		return &s, nil
	}

	return nil, nil // No hay saldo aún
}
