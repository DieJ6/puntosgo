package equivalencia

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/nmarsollier/commongo/db"
	"github.com/nmarsollier/commongo/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

type EquivalenciaRepository interface {
	FindByID(id primitive.ObjectID) (*Equivalencia, error)
	Insert(eq *Equivalencia) (*Equivalencia, error)
}

type repository struct {
	log        log.LogRusEntry
	collection db.Collection
}

func NewRepository(log log.LogRusEntry, c db.Collection) EquivalenciaRepository {
	return &repository{log, c}
}

func (r *repository) FindByID(id primitive.ObjectID) (*Equivalencia, error) {
	var e Equivalencia
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}, &e)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	return &e, nil
}

func (r *repository) Insert(e *Equivalencia) (*Equivalencia, error) {
	if err := validate.Struct(e); err != nil {
		r.log.Error(err)
		return nil, err
	}

	e.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(context.Background(), e)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	return e, nil
}
