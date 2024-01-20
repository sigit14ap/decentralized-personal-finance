package repositories

import (
	"github.com/sigit14ap/personal-finance/auth-service/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(user domain.User) error
	FindUserByUsername(username string) (domain.User, error)
}

type Repositories struct {
	Users UserRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users: NewUserRepository(db),
	}
}
