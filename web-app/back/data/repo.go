package data

import (
	"context"
	"errors"

	"github.com/Seifbarouni/private-git/web-app/back/db"
	"github.com/Seifbarouni/private-git/web-app/back/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repo struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Owner       primitive.ObjectID `bson:"owner" json:"owner"`
	Status      string             `bson:"status" json:"status"`
}

type RepoServiceInterface interface {
	CreateRepo(repo *Repo) error
	GetRepo(id string, userId string) (*Repo, error)
	GetRepos() ([]Repo, error)
	GetReposByOwner(owner string) ([]Repo, error)
	UpdateRepo(repo *Repo) error
	DeleteRepo(id string) error
}

type RepoService struct{}

func InitRepoService() *RepoService {
	return &RepoService{}
}

var reposCol string = "repos"
var userService UserServiceInterface = InitUserService()

func (rs *RepoService) CreateRepo(repo *Repo) error {
	user, err := userService.GetUser(repo.Owner)
	if err != nil {
		return err
	}

	err = utils.AddUserToRepo(user.UserName, user.SSHKey, repo.Name, "RW+")
	if err != nil {
		return err
	}

	repo.ID = primitive.NewObjectID()
	_, err = db.Collection(reposCol).InsertOne(context.TODO(), repo)
	return err
}

func (rs *RepoService) GetRepo(id string, userId string) (*Repo, error) {
	var repo Repo
	err := db.Collection(reposCol).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&repo)
	if err != nil {
		return nil, err
	}

	if repo.Owner.Hex() != userId {
		return nil, errors.New("unauthorized")
	}

	return &repo, nil
}

func (rs *RepoService) GetRepos() ([]Repo, error) {
	var repos []Repo
	cursor, err := db.Collection(reposCol).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var repo Repo
		cursor.Decode(&repo)
		repos = append(repos, repo)
	}
	if repos == nil {
		repos = []Repo{}
	}

	return repos, nil
}

func (rs *RepoService) GetReposByOwner(owner string) ([]Repo, error) {
	var repos []Repo
	cursor, err := db.Collection(reposCol).Find(context.TODO(), bson.M{"owner": owner})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var repo Repo
		cursor.Decode(&repo)
		repos = append(repos, repo)
	}
	if repos == nil {
		repos = []Repo{}
	}

	return repos, nil
}

func (rs *RepoService) UpdateRepo(repo *Repo) error {
	_, err := db.Collection(reposCol).UpdateOne(context.TODO(), bson.M{"_id": repo.ID}, bson.M{"$set": repo})
	return err
}

func (rs *RepoService) DeleteRepo(id string) error {
	_, err := db.Collection(reposCol).UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"status": "deleted"}})
	return err
}
