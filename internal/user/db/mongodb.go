package db

import (
	"context"
	"fmt"
	"github.com/OlegKapat/Rest-api-mongo/internal/apperror"
	"github.com/OlegKapat/Rest-api-mongo/internal/user"
	"github.com/OlegKapat/Rest-api-mongo/pkg/logging"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("Create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("Faild create user with error %/v", err)
	}
	d.logger.Debug("Converte insertedId to objectId")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("Faild to convert objectId to hex brobably oid %/s", oid)
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, r error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to objectid hex = %s", id)
	}
	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			//Err entity not found
			return u, fmt.Errorf("ErrEntityNotFound")
		}
		return u, apperror.ErrNotFound
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user from db by id:%/s due to error %v", id, err)
	}
	return u, nil
}
func (d *db) FindAll(ctx context.Context) (u []user.User, r error) {

	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return u, fmt.Errorf("failed to find all users  due to error %v", err)
	}
	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents from cursor err %v", err)
	}
	return u, nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert userId to objectId ID=%s err:%/v", id, err)
	}
	filter := bson.M{"_id": objectID}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("Failed to execute query. error %v", err)
	}
	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}
	d.logger.Tracef("deleted %d documents", result.DeletedCount)
	return nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert userId to objectId ID=%s err:%/v", user.ID, err)
	}
	filter := bson.M{"_id": objectID}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user, error %v", err)
	}
	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("FAILD TO UNMARSHAL user bytes,error %v", err)
	}
	delete(updateUserObj, "_id")
	update := bson.M{
		"$set": updateUserObj,
	}
	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("faild to execute update user query,error %v", err)
	}
	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}
	d.logger.Tracef("matched %d documents and Modified %d documents", result.MatchedCount, result.ModifiedCount)
	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
