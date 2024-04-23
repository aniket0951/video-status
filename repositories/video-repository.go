package repositories

import (
	"context"
	"errors"

	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"

	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	dbconfig "github.com/aniket0951/Chatrapati-Maharaj/db-config"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var videoCategoryCollection = dbconfig.GetCollection(dbconfig.DB, "video_category")
var videosCollection = dbconfig.GetCollection(dbconfig.DB, "videos")

type VideoRepository interface {
	CreateCategory(categories models.VideoCategories) (models.VideoCategories, error)
	UpdateCategory(categories models.VideoCategories) (models.VideoCategories, error)
	GetAllCategory() ([]models.VideoCategories, error)
	GetCategoryById(catId primitive.ObjectID) (models.VideoCategories, error)
	DeleteCategory(categoryId primitive.ObjectID) error
	DuplicateCategory(categoryName string) (bool, error)


	AddVideo(video models.Videos, file multipart.File) error
	AddVideo2(video models.Videos) error
	AddVideo(video models.Videos, file multipart.File) (primitive.ObjectID, error)

	AddVideo(video models.Videos, file multipart.File) (primitive.ObjectID, error)

	GetAllVideos() ([]models.Videos, error)
	GetVideoByID(videoId primitive.ObjectID) (models.Videos, error)
	UpdateVideo(video models.Videos) error
	UpdateVideoVerification(video models.Videos) error
	DeleteVideo(videoId primitive.ObjectID) error


	FetchInActiveVideos() ([]models.Videos, error)
	ActiveVideo(video_id primitive.ObjectID, isActive bool) error
	IncreaseDownloadCount(video_id primitive.ObjectID) error

	IsFileKeyExists(fileKey string) (bool, error)

	VideoFullDetails(videoId primitive.ObjectID) (interface{}, error)

	VideoFullDetails(videoId primitive.ObjectID) (interface{}, error)

	Init() (context.Context, context.CancelFunc)
}

type videocategoriesrepo struct {
	collection          *mongo.Collection
	videoscollection    *mongo.Collection
	userVideoConnection *mongo.Collection
}

func NewVideoCategoriesRepository() VideoRepository {
	return &videocategoriesrepo{
		collection:          videoCategoryCollection,
		videoscollection:    videosCollection,
		userVideoConnection: userVideoCollection,
	}
}

func (db *videocategoriesrepo) Init() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, cancel
}

func (db *videocategoriesrepo) CreateCategory(categories models.VideoCategories) (models.VideoCategories, error) {

	categories.ID = primitive.NewObjectID()
	categories.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	categories.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	ctx, cancel := db.Init()
	defer cancel()

	_, err := db.collection.InsertOne(ctx, &categories)

	if err != nil {
		return models.VideoCategories{}, err
	}

	return categories, nil
}

func (db *videocategoriesrepo) UpdateCategory(categories models.VideoCategories) (models.VideoCategories, error) {
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "category_name", Value: categories.CategoryName},
			bson.E{Key: "category_desc", Value: categories.CategoryDescription},
			bson.E{Key: "is_active", Value: categories.IsCategoryActive},
			bson.E{Key: "updated_at", Value: primitive.NewDateTimeFromTime(time.Now())},
		}},
	}

	ctx, cancel := db.Init()
	defer cancel()

	result, err := db.collection.UpdateByID(ctx, categories.ID, update)

	if err != nil {
		return models.VideoCategories{}, err
	}

	if result.MatchedCount == 0 {
		return models.VideoCategories{}, errors.New("category not found for update")
	}

	if result.ModifiedCount == 0 {
		return models.VideoCategories{}, errors.New(helper.UPDATE_FAILED)
	}

	return categories, nil
}

func (db *videocategoriesrepo) GetAllCategory() ([]models.VideoCategories, error) {

	ctx, cancel := db.Init()
	defer cancel()

	cursor, curErr := db.collection.Find(ctx, bson.D{})

	if curErr != nil {
		return []models.VideoCategories{}, curErr
	}

	var result []models.VideoCategories
	err := cursor.All(ctx, &result)

	if err != nil {
		return []models.VideoCategories{}, err
	}

	return result, nil
}

func (db *videocategoriesrepo) GetCategoryById(catId primitive.ObjectID) (models.VideoCategories, error) {
	filter := bson.D{
		bson.E{Key: "_id", Value: catId},
	}

	ctx, cancel := db.Init()
	defer cancel()

	videoCategory := models.VideoCategories{}

	res := db.collection.FindOne(ctx, filter).Decode(&videoCategory)

	if res == mongo.ErrNoDocuments {
		return models.VideoCategories{}, errors.New("video category not found")
	}

	return videoCategory, nil
}

func (db *videocategoriesrepo) DeleteCategory(categoryId primitive.ObjectID) error {

	filter := bson.D{
		bson.E{Key: "_id", Value: categoryId},
	}

	ctx, cancel := db.Init()
	defer cancel()

	res, err := db.collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("category not found for update")
	}

	return nil
}

func (db *videocategoriesrepo) DuplicateCategory(categoryName string) (bool, error) {
	filter := bson.D{
		bson.E{Key: "category_name", Value: categoryName},
	}

	ctx, cancel := db.Init()
	defer cancel()
	category := models.VideoCategories{}
	_ = db.collection.FindOne(ctx, filter).Decode(&category)

	if (category == models.VideoCategories{}) {
		return false, nil

	}

	return true, errors.New("this category already exits")
}

func (db *videocategoriesrepo) AddVideo(video models.Videos, file multipart.File) (primitive.ObjectID, error) {

	video.ID = primitive.NewObjectID()
	video.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	video.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	video.IsVideoActive = false
	video.IsVerified = false
	video.IsPublished = false

	tempFile, err := ioutil.TempFile("static", "upload-*.mp4")

	if err != nil {
		return primitive.NewObjectID(), err
	}

	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(tempFile)

	fileBytes, fileReader := ioutil.ReadAll(file)

	if fileReader != nil {
		return video.ID, fileReader
	}

	_, err = tempFile.Write(fileBytes)
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(tempFile)

	video.VideoPath = path.Base(tempFile.Name())

	ctx, cancel := db.Init()
	defer cancel()

	_, insErr := db.videoscollection.InsertOne(ctx, &video)

	if insErr != nil {
		return primitive.NewObjectID(), insErr
	}

	return video.ID, nil
}


func (db *videocategoriesrepo) AddVideo2(video models.Videos) error {
	video.ID = primitive.NewObjectID()
	video.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	video.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	ctx, cancel := db.Init()
	defer cancel()

	_, err := db.videoscollection.InsertOne(ctx, &video)

	return err
}

func (db *videocategoriesrepo) GetAllVideos() ([]models.Videos, error) {


	queryOptions := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	filter := bson.M{"is_active": true}
	cursor, curErr := db.videoscollection.Find(ctx, filter, queryOptions)

	filter := []bson.M{
		bson.M{
			"$match": bson.M{
				"is_active": true,
			},
		},

		{"$sort": bson.M{"_id": -1}},
		{"$limit": 5},
	}

	cursor, curErr := db.videoscollection.Aggregate(context.TODO(), filter)

func (db *videocategoriesrepo) GetAllVideos() ([]models.Videos, error) {

	filter := []bson.M{
		bson.M{
			"$match": bson.M{
				"is_active": true,
			},
		},

		{"$sort": bson.M{"_id": -1}},
		{"$limit": 5},
	}

	cursor, curErr := db.videoscollection.Aggregate(context.TODO(), filter)


	if curErr != nil {
		return []models.Videos{}, curErr
	}

	var videos []models.Videos

	if err := cursor.All(context.TODO(), &videos); err != nil {
		return []models.Videos{}, err
	}

	return videos, nil
}

func (db *videocategoriesrepo) UpdateVideo(video models.Videos) error {
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "video_title", Value: video.VideoTitle},
			bson.E{Key: "video_desc", Value: video.VideoDescription},
			bson.E{Key: "is_active", Value: video.IsVideoActive},
			bson.E{Key: "v_cat_id", Value: video.VideoCategoriesID},
			bson.E{Key: "updated_at", Value: primitive.NewDateTimeFromTime(time.Now())},
		}},
	}

	ctx, cancel := db.Init()
	defer cancel()

	res, upErr := db.videoscollection.UpdateByID(ctx, video.ID, update)

	if upErr != nil {
		return upErr
	}

	if res.MatchedCount == 0 {
		return errors.New("video not found to update")
	}

	return nil
}

func (db *videocategoriesrepo) UpdateVideoVerification(video models.Videos) error {

	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "is_verified", Value: video.IsVerified},
			bson.E{Key: "is_published", Value: video.IsPublished},
			bson.E{Key: "is_active", Value: video.IsVideoActive},
			bson.E{Key: "updated_at", Value: primitive.NewDateTimeFromTime(time.Now())},
		}},
	}

	ctx, cancel := db.Init()
	defer cancel()

	res, upErr := db.videoscollection.UpdateByID(ctx, video.ID, update)

	if upErr != nil {
		return upErr
	}

	if res.MatchedCount == 0 {
		return errors.New("video not found to update")
	}

	return nil
}

func (db *videocategoriesrepo) GetVideoByID(videoId primitive.ObjectID) (models.Videos, error) {
	filter := bson.D{
		bson.E{Key: "_id", Value: videoId},
	}

	ctx, cancel := db.Init()
	defer cancel()

	video := models.Videos{}

	res := db.videoscollection.FindOne(ctx, filter).Decode(&video)
	if res == mongo.ErrNoDocuments {
		return models.Videos{}, errors.New("video not found for delete")
	}

	return video, nil
}

func (db *videocategoriesrepo) DeleteVideo(videoId primitive.ObjectID) error {
	video, err := db.GetVideoByID(videoId)
	if err != nil {
		return err
	}

	ctx, cancel := db.Init()
	defer cancel()

	filter := bson.D{
		bson.E{Key: "_id", Value: videoId},
	}

	res, err := db.videoscollection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("failed to delete the video")
	}


	path := video.VideoPath
	var fileRemoveErr error
	if strings.Contains(path, "static") {
		fileRemoveErr = os.Remove(path)
	} else {
		fileRemoveErr = os.Remove("./static/" + path)
	}

	// fileRemoveErr = os.Remove(video.VideoPath)

	var fileRemoveErr error

	if strings.Contains(video.VideoPath, "static") {
		fileRemoveErr = os.Remove(video.VideoPath)
	} else {
		fileRemoveErr = os.Remove("static/" + video.VideoPath)
	}


	return fileRemoveErr
}


func (db *videocategoriesrepo) FetchInActiveVideos() ([]models.Videos, error) {
	ctx, cancel := db.Init()
	defer cancel()

	filter := bson.M{"is_active": false}

	cursor, err := db.videoscollection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var result []models.Videos

	if err := cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	return result, err
}

func (db *videocategoriesrepo) ActiveVideo(video_id primitive.ObjectID, isActive bool) error {
	ctx, cancel := db.Init()
	defer cancel()
	update := bson.M{
		"$set": bson.M{
			"is_active":  isActive,
			"updated_at": time.Now(),
		},
	}

	result, err := db.videoscollection.UpdateByID(ctx, video_id, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("video not found to update")
	}

	return nil
}

func (db *videocategoriesrepo) IncreaseDownloadCount(video_id primitive.ObjectID) error {
	ctx, cancel := db.Init()
	defer cancel()

	update := bson.M{
		"$inc": bson.M{
			"download_count": 1,
		},
	}

	filter := bson.M{"_id": video_id}

	_, err := db.videoscollection.UpdateOne(ctx, filter, update)

	return err
}

func (db *videocategoriesrepo) IsFileKeyExists(fileKey string) (bool, error) {
	filter := bson.M{"video_path": fileKey}

	ctx, cancel := db.Init()
	defer cancel()
	var video_data models.Videos
	err := db.videoscollection.FindOne(ctx, filter).Decode(&video_data)

	if err != nil {
		return false, err
	}

	if reflect.DeepEqual(video_data, &models.Videos{}) {
		return false, mongo.ErrNilDocument
	}

	return true, nil

func (db *videocategoriesrepo) VideoFullDetails(videoId primitive.ObjectID) (interface{}, error) {
	filter := []bson.M{
		bson.M{
			"$match": bson.M{
				"video_id": videoId,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user_data",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "videos",
				"localField":   "video_id",
				"foreignField": "_id",
				"as":           "videos_data",
			},
		},
	}

	cursor, curErr := db.userVideoConnection.Aggregate(context.TODO(), filter)

	if curErr != nil {
		return nil, curErr
	}

	var videoDetail []bson.M

	if err := cursor.All(context.TODO(), &videoDetail); err != nil {
		return nil, err
	}

	return videoDetail, nil
}
