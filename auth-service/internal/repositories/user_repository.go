package repositories

import (
	"context"
	"errors"

	"github.com/sigit14ap/personal-finance/auth-service/internal/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {

	collection := db.Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("unable to create user collection index, %v", err)
	}

	return &UserRepository{
		db: collection,
	}
}

func (repo *UserRepository) CreateUser(user domain.User) (domain.User, error) {
	_, err := repo.db.InsertOne(context.Background(), user)
	return user, err
}

func (repo *UserRepository) FindUserByEmail(email string) (domain.User, error) {
	var user domain.User

	filter := bson.M{"email": email}
	err := repo.db.FindOne(context.Background(), filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return domain.User{}, errors.New("user not found")
	}

	return user, err
}

func (repo *UserRepository) FindUserById(id string) (domain.User, error) {
	var user domain.User

	filter := bson.M{"_id": id}
	err := repo.db.FindOne(context.Background(), filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return domain.User{}, errors.New("user not found")
	}

	return user, err
}
