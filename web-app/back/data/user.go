package data

import (
	"context"
	"errors"

	"github.com/Seifbarouni/private-git/web-app/back/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID   `bson:"_id" json:"id"`
	FullName string               `bson:"full_name" json:"full_name" validate:"required"`
	Email    string               `bson:"email" json:"email" validate:"email,required"`
	UserName string               `bson:"user_name" json:"user_name" validate:"required"`
	Password string               `bson:"password" json:"password" validate:"required"`
	Status   string               `bson:"status" json:"status"`
	SSHKey   string               `bson:"ssh_key" json:"ssh_key"`
	Repos    []primitive.ObjectID `bson:"repos" json:"repos"`
}

type SSHKey struct {
	Key string `json:"ssh_key" validate:"required"`
}

type UserServiceInterface interface {
	CreateUser(user *User) error
	GetUser(id primitive.ObjectID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id primitive.ObjectID) error
	AddPublicKey(id primitive.ObjectID, key string) error
	AddRepoToUser(repoId primitive.ObjectID, userId primitive.ObjectID) error
}

type UserService struct{}

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
	user.Status = "active"
	user.SSHKey = ""
	user.Repos = []primitive.ObjectID{}
	_, err = u.InsertOne(context.TODO(), user)

	return err
}

func (us *UserService) GetUser(id primitive.ObjectID) (*User, error) {
	var user User
	err := db.Collection(usersCol).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
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

func (us *UserService) AddPublicKey(id primitive.ObjectID, key string) error {
	_, err := db.Collection(usersCol).UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"ssh_key": key}})
	return err
}

func (us *UserService) UpdateUser(user *User) error {
	_, err := db.Collection(usersCol).ReplaceOne(context.TODO(), bson.M{"_id": user.ID}, user)
	return err
}

func (us *UserService) DeleteUser(id primitive.ObjectID) error {
	_, err := db.Collection(usersCol).UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"status": "deleted"}})
	return err
}

func (us *UserService) AddRepoToUser(repoId primitive.ObjectID, userId primitive.ObjectID) error {
	_, err := db.Collection(usersCol).UpdateOne(context.TODO(), bson.M{"_id": userId}, bson.M{"$addToSet": bson.M{"repos": repoId}})
	return err
}
