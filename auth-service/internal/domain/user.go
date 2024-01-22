package domain

type User struct {
	Username string `bson:"username" validate:"required"`
	Password string `bson:"password" validate:"required"`
}
