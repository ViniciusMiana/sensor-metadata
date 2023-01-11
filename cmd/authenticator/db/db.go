package db

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "users"

// User represents a user with a username and password
type User struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password"  json:"password"`
	Role     string `bson:"role" json:"role"`
}

// UserStore represents a store for users
type UserStore struct {
	client   *mongo.Client
	database *mongo.Database
	users    *mongo.Collection
}

// NewUserStore creates a new user store
func NewUserStore(uri, databaseName string) (*UserStore, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	// TODO add credentials for connection
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	database := client.Database(databaseName)
	users := database.Collection(collectionName)
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err = users.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return nil, err
	}
	// TODO move this to service
	pass := os.Getenv("ROOT_PASSWORD")
	if pass == "" {
		pass = "1234"
	}
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	root := User{
		Username: "root",
		Password: string(hashedPasswordBytes),
		Role:     "ADMIN",
	}
	_, err = users.InsertOne(ctx, root)
	if err != nil {
		// We are ignoring this error to avoid root kidnapping.
		fmt.Println("Error while inserting root " + err.Error())
	}
	return &UserStore{client: client, database: database, users: users}, nil
}

// AddUser adds a new user to the store
func (store *UserStore) AddUser(ctx context.Context, user User) error {
	_, err := store.users.InsertOne(ctx, user)
	return err
}

func (store *UserStore) FindByUserName(ctx context.Context, username string) (*User, error) {
	filter := bson.M{"username": username}
	var result User
	err := store.users.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
