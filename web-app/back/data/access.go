package data

import (
	"context"

	"github.com/Seifbarouni/private-git/web-app/back/db"
	"github.com/Seifbarouni/private-git/web-app/back/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Access struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserId    primitive.ObjectID `bson:"user_id" json:"user_id"`
	RepoId    primitive.ObjectID `bson:"repo_id" json:"repo_id"`
	GrantType string             `bson:"grant_type" json:"grant_type"`
}

type AccessServiceInterface interface {
	GrantAccess(access *Access) error
	GetAccessesByRepoId(repoId string) ([]Access, error)
	GetAccessesByUserId(userId string) ([]Access, error)
	RevokeAccess(userId string, repo *Repo) error
}

type AccessService struct{}

func InitAccessService() *AccessService {
	return &AccessService{}
}

var accessCol string = "accesses"

func (ac *AccessService) GrantAccess(access *Access) error {
	user, err := userService.GetUser(access.UserId)
	if err != nil {
		return err
	}

	err = utils.AddUserToRepo(user.UserName, user.SSHKey, access.RepoId.Hex(), access.GrantType)

	if err != nil {
		return err
	}

	access.ID = primitive.NewObjectID()
	_, err = db.Collection(accessCol).InsertOne(context.TODO(), access)
	return err
}

func (ac *AccessService) GetAccessesByRepoId(repoId string) ([]Access, error) {
	var accesses []Access
	cursor, err := db.Collection(accessCol).Find(context.TODO(), bson.M{"repo_id": repoId})
	if err != nil {
		return accesses, err
	}

	if err = cursor.All(context.TODO(), &accesses); err != nil {
		return accesses, err
	}

	return accesses, nil
}

func (ac *AccessService) GetAccessesByUserId(repoId string) ([]Access, error) {
	var accesses []Access
	cursor, err := db.Collection(accessCol).Find(context.TODO(), bson.M{"user_id": repoId})
	if err != nil {
		return accesses, err
	}

	if err = cursor.All(context.TODO(), &accesses); err != nil {
		return accesses, err
	}

	return accesses, nil
}

func (ac *AccessService) RevokeAccess(userId string, repo *Repo) error {
	primUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	user, err := userService.GetUser(primUserId)
	if err != nil {
		return err
	}

	err = utils.RemoveUserFromRepo(user.UserName, repo.Name)
	if err != nil {
		return err
	}

	_, err = db.Collection(accessCol).DeleteOne(context.TODO(), bson.M{"user_id": userId, "repo_id": repo.ID.Hex()})
	return err
}
