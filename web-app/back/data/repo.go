package data

import (
	"context"

	"github.com/Seifbarouni/private-git/web-app/back/db"

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
	GetRepo(id primitive.ObjectID) (*Repo, error)
	GetRepos() ([]Repo, error)
	GetReposByOwner(owner primitive.ObjectID) ([]Repo, error)
	UpdateRepo(repo *Repo) error
	DeleteRepo(id primitive.ObjectID) error
}

type RepoService struct{}

func InitRepoService() *RepoService {
	return &RepoService{}
}

var reposCol string = "repos"

func (rs *RepoService) CreateRepo(repo *Repo) error {
	repo.ID = primitive.NewObjectID()
	_, err := db.Collection(reposCol).InsertOne(context.TODO(), repo)
	return err
}

func (rs *RepoService) GetRepo(id primitive.ObjectID) (*Repo, error) {
	var repo Repo
	err := db.Collection(reposCol).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&repo)
	return &repo, err
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
	return repos, nil
}

func (rs *RepoService) GetReposByOwner(owner primitive.ObjectID) ([]Repo, error) {
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
	return repos, nil
}

func (rs *RepoService) UpdateRepo(repo *Repo) error {
	_, err := db.Collection(reposCol).UpdateOne(context.TODO(), bson.M{"_id": repo.ID}, bson.M{"$set": repo})
	return err
}

func (rs *RepoService) DeleteRepo(id primitive.ObjectID) error {
	_, err := db.Collection(reposCol).UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"status": "deleted"}})
	return err
}
