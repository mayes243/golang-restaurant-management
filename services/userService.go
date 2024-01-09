package services

import (
	"context"
	"errors"
	userRepo "golang-restaurant-management/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsers(ctx context.Context, limit int, page int) (map[string]interface{}, error) {
	skip := (page - 1) * limit

	pipeline := bson.A{
		bson.D{
			{Key: "$facet", Value: bson.D{
				{Key: "metadata", Value: bson.A{
					bson.D{{Key: "$count", Value: "totalUsers"}},
					bson.D{
						{
							Key: "$addFields", Value: bson.D{
								{Key: "currentPage", Value: page},
								{Key: "limit", Value: limit},
								{Key: "totalPages", Value: bson.D{{Key: "$ceil", Value: bson.A{"$totalUsers"}}}},
							},
						},
					},
				}},
				{Key: "users", Value: bson.A{
					bson.D{{Key: "$skip", Value: skip}},
					bson.D{{Key: "$limit", Value: limit}},
					bson.D{
						{Key: "$project", Value: bson.D{
							{Key: "refresh_token", Value: 0},
							{Key: "password", Value: 0},
							{Key: "token", Value: 0},
						}},
					},
				}},
			}},
		},
	}

	result, err := userRepo.ExecuteAggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var allUsers []bson.M
	if err := result.All(ctx, &allUsers); err != nil {
		return nil, err
	}

	metadataArray, ok := allUsers[0]["metadata"].(bson.A)
	if !ok {
		return nil, errors.New("metadata is not in the expected format")
	}

	response := map[string]interface{}{
		"totalUsers":  metadataArray[0].(bson.M)["totalUsers"],
		"currentPage": metadataArray[0].(bson.M)["currentPage"],
		"limit":       metadataArray[0].(bson.M)["limit"],
		"totalPages":  metadataArray[0].(bson.M)["totalPages"],
		"users":       allUsers[0]["users"],
	}

	return response, nil
}

func GetUserByID(ctx context.Context, userID string) (map[string]interface{}, error) {

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userID}}}},
		bson.D{{Key: "$project", Value: bson.D{{Key: "password", Value: 0}}}},
	}

	result, err := userRepo.ExecuteAggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var user bson.M
	if result.Next(ctx) {
		if err := result.Decode(&user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func UserExists(ctx context.Context, data interface{}) (bool, error) {
	count, err := userRepo.CountDocuments(ctx, data)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateUser(ctx context.Context, user interface{}) (map[string]interface{}, error) {
	resultInsertionNumber, err := userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Retrieve the inserted user data
	createdUser, err := GetUserByID(ctx, resultInsertionNumber.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
