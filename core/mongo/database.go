package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewObjectID(id string) (primitive.ObjectID, error) {
	i, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		err = errors.New(id + " is not a valid mongo id")
	}
	return i, err
}

type Document[T any] interface {
	EnsureIndexes(Database)
	GetValue() *T
	Validate() error
}

type Database interface {
	GetInstance() *database
	Connect()
	Disconnect()
}

type database struct {
	*mongo.Database
	context     context.Context
	user        string
	pwd         string
	host        string
	port        uint16
	name        string
	minPoolSize uint16
	maxPoolSize uint16
}

func NewDatabase(ctx context.Context, env *config.Env) Database {
	db := database{
		context:     ctx,
		user:        env.DBUser,
		pwd:         env.DBUserPwd,
		host:        env.DBHost,
		port:        env.DBPort,
		name:        env.DBName,
		minPoolSize: env.DBMinPoolSize,
		maxPoolSize: env.DBMaxPoolSize,
	}
	return &db
}

func (db *database) GetInstance() *database {
	return db
}

func (db *database) Connect() {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/%s",
		db.user, db.pwd, db.host, db.port, db.name,
	)

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetMaxPoolSize(uint64(db.maxPoolSize))
	clientOptions.SetMaxPoolSize(uint64(db.minPoolSize))

	fmt.Println("Connecting Mongo...")
	client, err := mongo.Connect(db.context, clientOptions)
	if err != nil {
		log.Fatal("Connection to Mongo Failed!: ", err)
	}

	err = client.Ping(db.context, nil)
	if err != nil {
		log.Panic("Pinging to Mongo Failed!: ", err)
	}
	fmt.Println("Connected to Mongo!")

	db.Database = client.Database(db.name)
}

func (db *database) Disconnect() {
	fmt.Println("Disconnecting Mongo...")
	err := db.Client().Disconnect(db.context)
	if err != nil {
		log.Panic(err)
	}
}
