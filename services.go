package honuadb

import (
	"context"

	"github.com/JonasBordewick/honua-db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (hdb *HonuaDB) AddServices(identity string, services []*models.Service) error {
	exist, err := hdb.HasServices(identity)
	if err != nil {
		return err
	}

	if exist {
		err = hdb.DeleteServices(identity)
		if err != nil {
			return err
		}
	}

	_, err = hdb.mongoDB.Collection("services").InsertOne(context.Background(), &models.Services{
		Id: identity,
		Services: services,
	})

	return err
}

func (hdb *HonuaDB) GetServices(id string) ([]*models.Service, error) {
	filter := bson.M{"_id": id}
	var result *models.Services
	err := hdb.mongoDB.Collection("services").FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Services, nil
}

func (hdb *HonuaDB) DeleteServices(id string) error {
	filter := bson.M{"_id": id}
	_, err := hdb.mongoDB.Collection("services").DeleteOne(context.Background(), filter)
	return err
}

func (hdb *HonuaDB) HasServices(id string) (bool, error) {
	filter := bson.M{"_id": id}
	var result *models.Services
	err := hdb.mongoDB.Collection("services").FindOne(context.Background(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
