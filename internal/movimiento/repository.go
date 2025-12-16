package movimiento

import (
	"context"
	"time"
	"sort"

	"github.com/go-playground/validator/v10"
	"github.com/nmarsollier/commongo/db"
	"github.com/nmarsollier/commongo/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

type MovimientoRepository interface {
	Insert(m *Movimiento) (*Movimiento, error)
	FindByUsuario(uid primitive.ObjectID) ([]*Movimiento, error)
	FindByUsuarioAfter(uid primitive.ObjectID, after time.Time) ([]*Movimiento, error)
}

type repository struct {
	log        log.LogRusEntry
	collection db.Collection
}

func NewRepository(log log.LogRusEntry, c db.Collection) MovimientoRepository {
	return &repository{log, c}
}

func (r *repository) Insert(m *Movimiento) (*Movimiento, error) {
	if err := validate.Struct(m); err != nil {
		r.log.Error(err)
		return nil, err
	}

	m.ID = primitive.NewObjectID()
	m.FechaCreacion = time.Now()

	_, err := r.collection.InsertOne(context.Background(), m)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	return m, nil
}

func (r *repository) FindByUsuario(uid primitive.ObjectID) ([]*Movimiento, error) {
	var result []*Movimiento

	cur, err := r.collection.Find(context.Background(), bson.M{"forK_id_usuario": uid})
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var mv Movimiento
		if err := cur.Decode(&mv); err != nil {
			return nil, err
		}
		result = append(result, &mv)
	}

	// Ordenar por fechaCreacion DESC (más nuevo primero)
	sort.Slice(result, func(i, j int) bool {
		return result[i].FechaCreacion.After(result[j].FechaCreacion)
	})

	return result, nil
}

func (r *repository) FindByUsuarioAfter(uid primitive.ObjectID, after time.Time) ([]*Movimiento, error) {
	var result []*Movimiento

	filter := bson.M{
		"forK_id_usuario": uid,
		"fechaCreacion":   bson.M{"$gt": after},
	}

	cur, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var mv Movimiento
		if err := cur.Decode(&mv); err != nil {
			return nil, err
		}
		result = append(result, &mv)
	}

	return result, nil
}

