package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	defaultDatabaseName        = "gamedb"
	databaseCollectionAccounts = "accounts"
)

type Mongo struct {
	client *mongo.Client
}

func NewMongo(client *mongo.Client) *Mongo {
	return &Mongo{
		client: client,
	}
}

func (m *Mongo) SaveAccount(sa *StorageAccount) error {
	collection := m.client.Database(defaultDatabaseName).Collection(databaseCollectionAccounts)

	ctx := context.Background()
	res, err := collection.InsertOne(ctx, sa)
	if err != nil {
		return err
	}

	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("can't convert inserted id to object id")
	}

	sa.ID = objectID

	return nil
}

func (m *Mongo) UpdateAccount(sa *StorageAccount) error {
	collection := m.client.Database(defaultDatabaseName).Collection(databaseCollectionAccounts)

	filter := bson.D{{"_id", sa.ID}}

	ctx := context.Background()
	_, err := collection.ReplaceOne(ctx, filter, sa)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) GetAccount(id string) (StorageAccount, bool) {
	collection := m.client.Database(defaultDatabaseName).Collection(databaseCollectionAccounts)

	var res StorageAccount

	ctx := context.Background()
	filter := bson.D{{"_id", id}}

	err := collection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		fmt.Println("can't do find one:", err)
		return StorageAccount{}, false
	}

	return res, true
}

func (m *Mongo) GetAccountByLogin(login string) (StorageAccount, bool) {
	collection := m.client.Database(defaultDatabaseName).Collection(databaseCollectionAccounts)

	var res StorageAccount

	ctx := context.Background()
	filter := bson.D{{"login", login}}

	err := collection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		fmt.Println("can't do find one:", err)
		return StorageAccount{}, false
	}

	return res, true
}
