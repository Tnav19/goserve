package model

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongod "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CollectionName = "blogs"

type Blog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title" validate:"required,max=500"`
	Description string             `bson:"description" validate:"required,max=2000"`
	Text        *string            `bson:"text,omitempty"`
	DraftText   string             `bson:"draftText" validate:"required"`
	Tags        []string           `bson:"tags"`
	Author      primitive.ObjectID `bson:"author" validate:"required"`
	ImgURL      *string            `bson:"imgUrl,omitempty"`
	Slug        string             `bson:"slug" validate:"required,unique,min=3,max=200"`
	Score       float64            `bson:"score" validate:"min=0,max=1"`
	IsSubmitted bool               `bson:"isSubmitted" validate:"required"`
	IsDraft     bool               `bson:"isDraft" validate:"required"`
	IsPublished bool               `bson:"isPublished" validate:"required"`
	Status      bool               `bson:"status"`
	PublishedAt *time.Time         `bson:"publishedAt,omitempty"`
	CreatedBy   primitive.ObjectID `bson:"createdBy" validate:"required"`
	UpdatedBy   primitive.ObjectID `bson:"updatedBy" validate:"required"`
	CreatedAt   time.Time          `bson:"createdAt" validate:"required"`
	UpdatedAt   time.Time          `bson:"updatedAt" validate:"required"`
}

func NewBlog(title, description, draftText, slug string, author, createdBy, updatedBy primitive.ObjectID, tags []string) (*Blog, error) {
	now := time.Now()
	b := Blog{
		Title:       title,
		Description: description,
		DraftText:   draftText,
		Tags:        tags,
		Author:      author,
		Slug:        slug,
		Score:       0.01,
		IsSubmitted: false,
		IsDraft:     true,
		IsPublished: false,
		Status:      true,
		CreatedBy:   createdBy,
		UpdatedBy:   updatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := b.Validate(); err != nil {
		return nil, err
	}
	return &b, nil
}

func (blog *Blog) Validate() error {
	validate := validator.New()
	return validate.Struct(blog)
}

func (*Blog) EnsureIndexes(db mongo.Database) {
	indexes := []mongod.IndexModel{
		{
			Keys: bson.D{{Key: "title", Value: "text"}, {Key: "description", Value: "text"}},
			Options: options.Index().SetWeights(bson.M{
				"title":       3,
				"description": 1,
			}),
		},
		{Keys: bson.D{{Key: "_id", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "slug", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "isPublished", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "_id", Value: 1}, {Key: "isPublished", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "tags", Value: 1}, {Key: "isPublished", Value: 1}, {Key: "status", Value: 1}}},
	}

	mongo.NewQueryBuilder[Blog](db, CollectionName).Query(context.Background()).CreateIndexes(indexes)
}
