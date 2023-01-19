package repositories

import (
	"context"
	"errors"
	"fmt"
	dbconfig "github.com/aniket0951/Chatrapati-Maharaj/db-config"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var videoVerificationCollection = dbconfig.GetCollection(dbconfig.DB, "video_verification")
var publishCollection = dbconfig.GetCollection(dbconfig.DB, "video_publish")
var verificationNotificationCollection = dbconfig.GetCollection(dbconfig.DB, "verification_notification")
var videoProcessHistoryCollection = dbconfig.GetCollection(dbconfig.DB, "video_process_history")

type VideoVerificationRepository interface {
	Init() (context.Context, context.CancelFunc)
	CreateVerification(verification models.VideoVerification) error
	GetAllVideosVerification() ([]models.VideoVerification, error)
	ApproveOrDeniedVideo(videoId primitive.ObjectID, verificationStatus string) error
	VideosForVerification(tag string) ([]models.VideoVerification, error)
	GetVideoVerificationByVideoId(videoId primitive.ObjectID) (models.VideoVerification, error)
	DeleteVerificationByVideoId(videoId primitive.ObjectID) error

	PublishedVideo(publish models.VideoPublish) error
	GetAllPublishData() ([]models.VideoPublish, error)
	GetPublishVideoByVideoId(videoId primitive.ObjectID) (models.VideoPublish, error)
	DeletePublishByVideoId(videoId primitive.ObjectID) error
	CreateVideoProcessHistory(history models.VideoProcessHistory) error

	CreateVerificationNotification(notification models.VideoVerificationNotification) error
	GetUserVerificationNotification(userId primitive.ObjectID) ([]models.VideoVerificationNotification, error)
}

type videoverification struct {
	videoVerificationConnection   *mongo.Collection
	videoPublishConnection        *mongo.Collection
	notificationConnection        *mongo.Collection
	videoProcessHistoryConnection *mongo.Collection
}

func NewVideoVerificationRepository() VideoVerificationRepository {
	return &videoverification{
		videoVerificationConnection:   videoVerificationCollection,
		videoPublishConnection:        publishCollection,
		notificationConnection:        verificationNotificationCollection,
		videoProcessHistoryConnection: videoProcessHistoryCollection,
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
func (db *videoverification) ApproveOrDeniedVideo(videoId primitive.ObjectID, verificationStatus string) error {
	filter := bson.D{
		bson.E{Key: "video_id", Value: videoId},
	}

	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "verification_status", Value: verificationStatus},
			bson.E{Key: "updated_at", Value: primitive.NewDateTimeFromTime(time.Now())},
		}},
	}

	ctx, cancel := db.Init()
	defer cancel()

	res := db.videoVerificationConnection.FindOneAndUpdate(ctx, filter, update)

	if res.Err() == mongo.ErrNoDocuments {
		return errors.New("no verification found for this video")
	}

	return nil
}
func (db *videoverification) VideosForVerification(tag string) ([]models.VideoVerification, error) {

	filter := []bson.M{
		bson.M{
			"$match": bson.M{
				"verification_status": tag,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "videos",
				"localField":   "video_id",
				"foreignField": "_id",
				"as":           "video_info",
			},
		},
	}

	ctx, cancel := db.Init()
	defer cancel()

	cursor, curErr := db.videoVerificationConnection.Aggregate(ctx, filter)

	if curErr != nil {
		return nil, curErr
	}

	var approveVideo []models.VideoVerification

	if err := cursor.All(context.TODO(), &approveVideo); err != nil {
		return nil, err
	}

	return approveVideo, nil
}
func (db *videoverification) GetVideoVerificationByVideoId(videoId primitive.ObjectID) (models.VideoVerification, error) {
	filter := bson.D{
		bson.E{Key: "video_id", Value: videoId},
	}

	ctx, cancel := db.Init()
	defer cancel()

	videoVerification := models.VideoVerification{}
	err := db.videoVerificationConnection.FindOne(ctx, filter).Decode(&videoVerification)

	if err != nil {
		return models.VideoVerification{}, err
	}

	return videoVerification, nil
}
func (db *videoverification) DeleteVerificationByVideoId(videoId primitive.ObjectID) error {
	filter := bson.D{
		bson.E{Key: "video_id", Value: videoId},
	}

	ctx, cancel := db.Init()
	defer cancel()

	res, err := db.videoVerificationConnection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("verification not found for delete")
	}

	return nil
}

func (db *videoverification) PublishedVideo(publish models.VideoPublish) error {

	filter := bson.D{
		bson.E{Key: "video_id", Value: publish.VideoId},
	}

	opts := options.FindOneAndReplace().SetUpsert(true)

	ctx, cancel := db.Init()
	defer cancel()

	db.videoPublishConnection.FindOneAndReplace(ctx, filter, publish, opts)
	return nil
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
func (db *videoverification) GetPublishVideoByVideoId(videoId primitive.ObjectID) (models.VideoPublish, error) {
	filter := bson.D{
		bson.E{Key: "video_id", Value: videoId},
	}

	ctx, cancel := db.Init()
	defer cancel()

	publishData := models.VideoPublish{}
	err := db.videoPublishConnection.FindOne(ctx, filter).Decode(&publishData)
	if err != nil {
		return models.VideoPublish{}, err
	}

	return publishData, nil
}
func (db *videoverification) DeletePublishByVideoId(videoId primitive.ObjectID) error {
	filter := bson.D{
		bson.E{Key: "video_id", Value: videoId},
	}

	ctx, cancel := db.Init()
	defer cancel()

	res, err := db.videoPublishConnection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("video publish not found")
	}

	return nil

}
func (db *videoverification) CreateVideoProcessHistory(history models.VideoProcessHistory) error {
	ctx, cancel := db.Init()
	defer cancel()

	res, err := db.videoProcessHistoryConnection.InsertOne(ctx, history)

	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func (db *videoverification) CreateVerificationNotification(notification models.VideoVerificationNotification) error {

	filter := bson.D{
		bson.E{Key: "video_id", Value: notification.VideoId},
	}

	opts := options.FindOneAndReplace().SetUpsert(true)

	ctx, cancel := db.Init()
	defer cancel()

	//_, err := db.notificationConnection.InsertOne(ctx, notification)

	db.notificationConnection.FindOneAndReplace(ctx, filter, notification, opts)

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
