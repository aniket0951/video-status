package repositories

import (
	"context"
	dbconfig "github.com/aniket0951/Chatrapati-Maharaj/db-config"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var videoVerificationCollection *mongo.Collection = dbconfig.GetCollection(dbconfig.DB, "video_verification")
var publishCollection *mongo.Collection = dbconfig.GetCollection(dbconfig.DB, "video_publish")
var verificationNotificationCollection *mongo.Collection = dbconfig.GetCollection(dbconfig.DB, "verification_notification")

type VideoVerificationRepository interface {
	Init() (context.Context, context.CancelFunc)
	CreateVerification(verification models.VideoVerification) error
	GetAllVideosVerification() ([]models.VideoVerification, error)

	CreatePublish(publish models.VideoPublish) error
	GetAllPublishData() ([]models.VideoPublish, error)

	CreateVerificationNotification(notification models.VideoVerificationNotification) error
	GetUserVerificationNotification(userId primitive.ObjectID) ([]models.VideoVerificationNotification, error)
}

type videoverification struct {
	videoVerificationConnection *mongo.Collection
	videoPublishConnection      *mongo.Collection
	notificationConnection      *mongo.Collection
}

func NewVideoVerificationRepository() VideoVerificationRepository {
	return &videoverification{
		videoVerificationConnection: videoVerificationCollection,
		videoPublishConnection:      publishCollection,
		notificationConnection:      verificationNotificationCollection,
	}
}

func (db *videoverification) Init() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, cancel
}

func (db *videoverification) CreateVerification(verification models.VideoVerification) error {
	ctx, cancel := db.Init()
	defer cancel()

	_, err := db.videoVerificationConnection.InsertOne(ctx, &verification)

	return err
}

func (db *videoverification) GetAllVideosVerification() ([]models.VideoVerification, error) {

	filter := []bson.M{
		{"$limit": 10},
	}

	cursor, curErr := db.videoVerificationConnection.Aggregate(context.Background(), filter)

	if curErr != nil {
		return nil, curErr
	}

	var verificationData []models.VideoVerification

	if err := cursor.All(context.Background(), &verificationData); err != nil {
		return nil, err
	}

	return verificationData, nil
}

func (db *videoverification) CreatePublish(publish models.VideoPublish) error {
	ctx, cancel := db.Init()
	defer cancel()

	_, err := db.videoPublishConnection.InsertOne(ctx, publish)

	return err
}
func (db *videoverification) GetAllPublishData() ([]models.VideoPublish, error) {
	filter := []bson.M{
		{"$limit": 10},
	}

	ctx, cancel := db.Init()
	defer cancel()

	cursor, curErr := db.videoPublishConnection.Aggregate(ctx, filter)

	if curErr != nil {
		return nil, curErr
	}

	var publishData []models.VideoPublish

	if err := cursor.All(ctx, &publishData); err != nil {
		return nil, err
	}

	return publishData, nil
}

func (db *videoverification) CreateVerificationNotification(notification models.VideoVerificationNotification) error {
	ctx, cancel := db.Init()
	defer cancel()

	_, err := db.notificationConnection.InsertOne(ctx, notification)

	if err != nil {
		return err
	}

	return nil
}

func (db *videoverification) GetUserVerificationNotification(userId primitive.ObjectID) ([]models.VideoVerificationNotification, error) {

	ctx, cancel := db.Init()
	defer cancel()

	filter := bson.D{
		bson.E{Key: "user_id", Value: userId},
	}

	cursor, curErr := db.notificationConnection.Find(ctx, filter)

	if curErr != nil {
		return nil, curErr
	}

	var notifications []models.VideoVerificationNotification

	if err := cursor.All(context.TODO(), &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}