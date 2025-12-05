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
	return &repository{log: log, collection: c}
}

func (r *repository) Insert(s *Saldo) (*Saldo, error) {
	if err := validate.Struct(s); err != nil {
		r.log.Error(err)
		return nil, err
	}

	now := time.Now()
	s.ID = primitive.NewObjectID()
	s.FechaCreacion = now
	s.FechaModificacion = now

	if _, err := r.collection.InsertOne(context.Background(), s); err != nil {
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
		nil, // firma de commongo/db
	)

	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return s, nil
}

func (r *repository) FindLatestByUsuario(uid primitive.ObjectID) (*Saldo, error) {
	filter := bson.M{"forK_id_usuario": uid}

	cur, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	var latest *Saldo
	var latestTime time.Time

	for cur.Next(context.Background()) {
		s := &Saldo{}
		if err := cur.Decode(s); err != nil {
			r.log.Error(err)
			return nil, err
		}

		refDate := s.FechaModificacion
		if refDate.IsZero() {
			refDate = s.FechaCreacion
		}

		if latest == nil || refDate.After(latestTime) {
			latest = s
			latestTime = refDate
		}
	}

	// NO existe cur.Err() en commongo/db → simplemente lo omitimos.

	return latest, nil
}
