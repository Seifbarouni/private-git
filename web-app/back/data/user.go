package data

import (
	"context"
	"errors"

	"github.com/Seifbarouni/private-git/web-app/back/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	FullName string             `bson:"full_name" json:"full_name"`
	Email    string             `bson:"email" json:"email"`
	UserName string             `bson:"user_name" json:"user_name"`
	Password string             `bson:"password" json:"password"`
	Status   string             `bson:"status" json:"status"`
}

type UserServiceInterface interface {
	CreateUser(user *User) error
	GetUser(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int) error
}

type UserService struct {
}

func InitUserService() *UserService {
	return &UserService{}
}

var usersCol string = "users"

func (us *UserService) CreateUser(user *User) error {
	var userCheck User
	u := db.Collection(usersCol)
	err := u.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&userCheck)
	if err == nil {
		return errors.New("user already exists")
	}
	user.ID = primitive.NewObjectID()
	_, err = u.InsertOne(context.TODO(), user)

	return err
}

func (us *UserService) GetUser(id int) (*User, error) {
	var user User
	err := db.Collection(usersCol).FindOne(context.TODO(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) GetUserByEmail(email string) (*User, error) {
	var user User
	err := db.Collection(usersCol).FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) UpdateUser(user *User) error {
	_, err := db.Collection(usersCol).ReplaceOne(context.TODO(), bson.M{"id": user.ID}, user)
	return err
}

func (us *UserService) DeleteUser(id int) error {
	_, err := db.Collection(usersCol).UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": bson.M{"status": "deleted"}})
	return err
}
