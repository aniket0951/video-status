package notificationmanager

import (
	"context"
	"time"

	dbconfig "github.com/aniket0951/Chatrapati-Maharaj/db-config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var tokenConnection = dbconfig.GetCollection(dbconfig.DB, "fcm_token")

type NotificationManagerRepo interface {
	Init() (context.Context, context.CancelFunc)
	AddNewToken(tokenData TokenData) error
	GetTokens() ([]TokenData, error)
}

type nmRepo struct {
	tokenCollection *mongo.Collection
}

func NewNotificationManagerRepo() NotificationManagerRepo {
	return &nmRepo{tokenCollection: tokenConnection}
}

func (repo *nmRepo) Init() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx, cancel
}

func (repo *nmRepo) AddNewToken(tokenData TokenData) error {
	ctx, cancel := repo.Init()
	defer cancel()

	_, err := repo.tokenCollection.InsertOne(ctx, &tokenData)
	return err
}

func (repo *nmRepo) GetTokens() ([]TokenData, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	cursor, err := repo.tokenCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var tokenData []TokenData
	if err := cursor.All(context.Background(), &tokenData); err != nil {
		return nil, err
	}

	return tokenData, nil
}
