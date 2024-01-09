package repository

import (
	"context"
	"golang-restaurant-management/database"

	"go.mongodb.org/mongo-driver/mongo"
)

var UserCollection = database.OpenCollection(database.Client, "user")

func ExecuteAggregate(ctx context.Context, pipeline interface{}) (*mongo.Cursor, error) {
	return UserCollection.Aggregate(ctx, pipeline)
}

func CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	count, err := UserCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// create a user
func CreateUser(ctx context.Context, user interface{}) (*mongo.InsertOneResult, error) {
	resultInsertionNumber, err := UserCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return resultInsertionNumber, nil
}
