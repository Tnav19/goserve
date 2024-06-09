package schema

import (
	"context"
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongod "go.mongodb.org/mongo-driver/mongo"
)

const CollectionName = "api_keys"

type Permission string

const (
	GeneralPermission Permission = "GENERAL"
)

type ApiKey struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Key         string             `bson:"key" validate:"required,max=1024"`
	Version     int                `bson:"version" validate:"required,min=1,max=100"`
	Permissions []Permission       `bson:"permissions" validate:"required"`
	Comments    []string           `bson:"comments" validate:"required,max=1000"`
	Status      bool               `bson:"status" validate:"-"`
	CreatedAt   time.Time          `bson:"createdAt" validate:"-"`
	UpdatedAt   time.Time          `bson:"updatedAt" validate:"-"`
}

func NewApiKey(key string, version int, permissions []Permission, comments []string) *ApiKey {
	currentTime := time.Now()
	return &ApiKey{
		Key:         key,
		Version:     version,
		Permissions: permissions,
		Comments:    comments,
		Status:      true,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}
}

func EnsureIndexes(db mongo.Database) {
	indexes := []mongod.IndexModel{
		{
			Keys: bson.D{
				{Key: "code", Value: 1},
				{Key: "status", Value: 1},
			},
		},
	}
	q := mongo.NewDatabaseQuery[ApiKey](db, CollectionName)
	q.CreateIndexes(context.Background(), indexes)
}
