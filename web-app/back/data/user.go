package data

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       int    `bson:"id"`
	FullName string `bson:"full_name"`
	Email    string `bson:"email"`
	UserName string `bson:"user_name"`
	Password string `bson:"password"`
	Status   string `bson:"status"`
}

type UserServiceInterface interface {
	CreateUser(user *User) error
	GetUser(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int) error
}

type UserService struct {
	collection *mongo.Collection
}

func InitUserService(col *mongo.Collection) *UserService {
	return &UserService{collection: col}
}

func (us *UserService) CreateUser(user *User) error {
	var userCheck User
	err := us.collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&userCheck)
	if err == nil {
		return errors.New("user already exists")
	}
	_, err = us.collection.InsertOne(context.TODO(), user)

	return err
}

func (us *UserService) GetUser(id int) (*User, error) {
	var user User
	err := us.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) GetUserByEmail(email string) (*User, error) {
	var user User
	err := us.collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) UpdateUser(user *User) error {
	_, err := us.collection.ReplaceOne(context.TODO(), bson.M{"id": user.ID}, user)
	return err
}

func (us *UserService) DeleteUser(id int) error {
	_, err := us.collection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": bson.M{"status": "deleted"}})
	return err
}
