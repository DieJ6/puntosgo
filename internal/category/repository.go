package category

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nmarsollier/commongo/db"
	"github.com/nmarsollier/commongo/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

type CategoryRepository interface {
	Insert(cat *Category) (*Category, error)
	Update(cat *Category) (*Category, error)
	FindByID(id primitive.ObjectID) (*Category, error)
	FindByArticuloID(productID string) (*Category, error)
	AddArticulo(catID primitive.ObjectID, productID string) error
	RemoveArticulo(catID primitive.ObjectID, productID string) error
}

type categoryRepository struct {
	log        log.LogRusEntry
	collection db.Collection
}

func NewRepository(
	log log.LogRusEntry,
	collection db.Collection,
) CategoryRepository {
	return &categoryRepository{log, collection}
}

func (r *categoryRepository) Insert(cat *Category) (*Category, error) {
	if err := validate.Struct(cat); err != nil {
		r.log.Error(err)
		return nil, err
	}

	cat.ID = primitive.NewObjectID()
	cat.FechaCreacion = time.Now()

	_, err := r.collection.InsertOne(context.Background(), cat)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return cat, nil
}

func (r *categoryRepository) Update(cat *Category) (*Category, error) {
	if err := validate.Struct(cat); err != nil {
		r.log.Error(err)
		return nil, err
	}

	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": cat.ID},
		bson.M{"$set": cat},
	)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return cat, nil
}

func (r *categoryRepository) FindByID(id primitive.ObjectID) (*Category, error) {
	var cat Category
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}, &cat)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) FindByArticuloID(productID string) (*Category, error) {
	var cat Category
	err := r.collection.FindOne(context.Background(),
		bson.M{"articulos": productID},
		&cat,
	)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) AddArticulo(catID primitive.ObjectID, productID string) error {

	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": catID},
		bson.M{"$addToSet": bson.M{"articulos": productID}},
	)
	if err != nil {
		r.log.Error(err)
	}
	return err
}

func (r *categoryRepository) RemoveArticulo(catID primitive.ObjectID, productID string) error {

	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": catID},
		bson.M{"$pull": bson.M{"articulos": productID}},
	)
	if err != nil {
		r.log.Error(err)
	}
	return err
}
