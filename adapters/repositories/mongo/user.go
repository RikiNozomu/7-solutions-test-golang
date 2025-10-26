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

type MongoUserRepository struct {
	collection *mongo.Collection // MongoDB collection for users
}

// createUniqueIndex ensures the "email" field is unique across documents.
func createUniqueIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}}, // Index on "email"
		Options: options.Index().SetUnique(true),  // Enforce uniqueness
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}

// NewMongoUserRepository initializes the repository and sets up indexes.
func NewMongoUserRepository(collection *mongo.Collection) domain.UserRepository {
	createUniqueIndex(collection)
	return &MongoUserRepository{collection}
}

// Create inserts a new user document into MongoDB.
func (m *MongoUserRepository) Create(user domain.User) error {
	_, err := m.collection.InsertOne(context.TODO(), user)
	if mongo.IsDuplicateKeyError(err) {
		return util.ErrorDuplicateKey
	}
	return err
}

// GetOne retrieves a user by key ("_id" or "email").
func (m *MongoUserRepository) GetOne(key string, value any) (*domain.User, error) {
	var user domain.User
	checkValue := value

	// Convert string ID to ObjectID if key is "_id"
	if key == "_id" {
		objectID, err := primitive.ObjectIDFromHex(value.(string))
		if err != nil {
			return nil, mongo.ErrNoDocuments
		}
		checkValue = objectID
	}

	// Query MongoDB and decode result
	err := m.collection.FindOne(context.TODO(), bson.D{{Key: key, Value: checkValue}}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all user documents from MongoDB.
func (m *MongoUserRepository) GetAll() ([]domain.User, error) {
	cursor, err := m.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	if cursor.RemainingBatchLength() <= 0 {
		return []domain.User{}, nil
	}

	// Decode each document into a User struct
	var users []domain.User
	for cursor.Next(context.TODO()) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, err
}

// GetCount returns the total number of user documents.
func (m *MongoUserRepository) GetCount() (int64, error) {
	return m.collection.CountDocuments(context.TODO(), bson.D{})
}

// Update modifies an existing user document by ID.
func (m *MongoUserRepository) Update(user domain.User) error {
	_, err := m.collection.UpdateByID(context.TODO(), user.ID, bson.M{"$set": user})
	if mongo.IsDuplicateKeyError(err) {
		return util.ErrorDuplicateKey
	}
	return err
}

// Delete removes a user document by ID.
func (m *MongoUserRepository) Delete(id string) error {
	// First, retrieve the user to get the ObjectID
	oldData, err := m.GetOne("_id", id)
	if err != nil {
		return err
	}

	// Delete the document by ObjectID
	result := m.collection.FindOneAndDelete(context.TODO(), bson.D{{Key: "_id", Value: oldData.ID}})
	return result.Err()
}
