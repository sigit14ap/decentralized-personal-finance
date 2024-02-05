package repositories

import (
	"github.com/sigit14ap/github.com/sigit14ap/decentralized-personal-finance-auth-service/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	CreateUser(user domain.User) (domain.User, error)
	FindUserByEmail(username string) (domain.User, error)
}

type Repositories struct {
	Users Users
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users: NewUserRepository(db),
	}
}
