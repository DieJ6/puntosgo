package mongo

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
    Client *mongo.Client
    DB     *mongo.Database
}

func NewDatabase(url string, dbName string) (*Database, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
    if err != nil {
        return nil, err
    }

    return &Database{
        Client: client,
        DB:     client.Database(dbName),
    }, nil
}

func (d *Database) Collection(name string) *mongo.Collection {
    return d.DB.Collection(name)
}
