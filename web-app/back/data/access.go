package data

import (
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
	RevokeAccess(userId string, repoId string) error
}

type AccessService struct{}

func InitAccessService() *AccessService {
	return &AccessService{}
}

func (ac *AccessService) GrantAccess(access *Access) error {
	return nil
}

func (ac *AccessService) GetAccessesByRepoId(repoId string) ([]Access, error) {
	return []Access{}, nil
}

func (ac *AccessService) GetAccessesByUserId(repoId string) ([]Access, error) {
	return []Access{}, nil
}

func (ac *AccessService) RevokeAccess(userId string, repoId string) error {
	return nil
}
