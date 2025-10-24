package mongo

import (
	domain "7-solutions-test-golang/core/domains"
	util "7-solutions-test-golang/utils"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoUserRepository is a secondary adapter that persists users to MongoDB.
type MongoUserRepository struct {
	collection *mongo.Collection
}

func createUniqueIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	return nil
}

func NewMongoUserRepository(collection *mongo.Collection) domain.UserRepository {
	createUniqueIndex(collection)
	return &MongoUserRepository{collection}
}

// Create implements domain.UserRepository.
func (m *MongoUserRepository) Create(user domain.User) error {
	_, err := m.collection.InsertOne(context.TODO(), user)
	if mongo.IsDuplicateKeyError(err) {
		return util.ErrorDuplicateKey
	}
	return err
}

// GetByEmail implements domain.UserRepository.
func (m *MongoUserRepository) GetOne(key string, value any) (*domain.User, error) {
	var user domain.User
	var checkValue = value
	if key == "_id" {
		objectID, err := primitive.ObjectIDFromHex(value.(string))
		if err != nil {
			return nil, mongo.ErrNoDocuments
		}
		checkValue = objectID
	}

	err := m.collection.FindOne(context.TODO(), bson.D{{Key: key, Value: checkValue}}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// GetAll implements domain.UserRepository.
func (m *MongoUserRepository) GetAll() ([]domain.User, error) {
	cursor, err := m.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	if cursor.RemainingBatchLength() <= 0 {
		return []domain.User{}, nil
	}

	// Iterate over the cursor and decode documents
	var users []domain.User
	for cursor.Next(context.TODO()) {
		var user domain.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, err
}

// GetCount implements domain.UserRepository.
func (m *MongoUserRepository) GetCount() (int64, error) {
	return m.collection.CountDocuments(context.TODO(), bson.D{})
}

// Update implements domain.UserRepository.
func (m *MongoUserRepository) Update(user domain.User) error {
	_, err := m.collection.UpdateByID(context.TODO(), user.ID, bson.M{"$set": user})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return util.ErrorDuplicateKey
		}
		return err
	}
	return nil
}

// Delete implements domain.UserRepository.
func (m *MongoUserRepository) Delete(id string) error {
	oldData, err := m.GetOne("_id", id)
	if err != nil {
		return err
	}
	result := m.collection.FindOneAndDelete(context.TODO(), bson.D{{Key: "_id", Value: oldData.ID}})
	return result.Err()
}
