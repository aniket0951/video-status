package repositories

import (
	"context"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
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

var videoCategoryCollection *mongo.Collection = dbconfig.GetCollection(dbconfig.DB, "video_category")
var videosCollection *mongo.Collection = dbconfig.GetCollection(dbconfig.DB, "videos")

type VideoRepository interface {
	CreateCategory(categories models.VideoCategories) (models.VideoCategories, error)
	UpdateCategory(categories models.VideoCategories) (models.VideoCategories, error)
	GetAllCategory() ([]models.VideoCategories, error)
	GetCategoryById(catId primitive.ObjectID) (models.VideoCategories, error)
	DeleteCategory(categoryId primitive.ObjectID) error
	DuplicateCategory(categoryName string) (bool, error)

	AddVideo(video models.Videos, file multipart.File) error
	GetAllVideos() ([]models.Videos, error)
	GetVideoByID(videoId primitive.ObjectID) (models.Videos, error)
	UpdateVideo(video models.Videos) error
	DeleteVideo(videoId primitive.ObjectID) error

	Init() (context.Context, context.CancelFunc)
}

type videocategoriesrepo struct {
	collection       *mongo.Collection
	videoscollection *mongo.Collection
}

func NewVideoCategoriesRepository() VideoRepository {
	return &videocategoriesrepo{
		collection:       videoCategoryCollection,
		videoscollection: videosCollection,
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

	result := []models.VideoCategories{}
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
	db.collection.FindOne(ctx, filter).Decode(&category)

	if (category == models.VideoCategories{}) {
		return false, nil

	}

	return true, errors.New("this category already exits")

}

func (db *videocategoriesrepo) AddVideo(video models.Videos, file multipart.File) error {

	video.ID = primitive.NewObjectID()
	video.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	video.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	tempFile, err := ioutil.TempFile("static", "upload-*.mp4")

	if err != nil {
		return err
	}

	defer tempFile.Close()

	fileBytes, fileReader := ioutil.ReadAll(file)

	if fileReader != nil {
		return fileReader
	}

	tempFile.Write(fileBytes)
	defer file.Close()
	defer tempFile.Close()

	video.VideoPath = path.Base(tempFile.Name())

	ctx, cancel := db.Init()
	defer cancel()

	_, insErr := db.videoscollection.InsertOne(ctx, &video)

	if insErr != nil {
		return insErr
	}

	return nil
}

func (db *videocategoriesrepo) GetAllVideos() ([]models.Videos, error) {
	ctx, cancel := db.Init()
	defer cancel()

	queryOptions := options.Find().SetSort(bson.D{{"_id", -1}})

	cursor, curErr := db.videoscollection.Find(ctx, bson.M{}, queryOptions)

	if curErr != nil {
		return []models.Videos{}, curErr
	}

	videos := []models.Videos{}

	if err := cursor.All(ctx, &videos); err != nil {
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

	var fileRemoveErr error

	if strings.Contains(video.VideoPath, "static") {
		fileRemoveErr = os.Remove(video.VideoPath)
	} else {
		fileRemoveErr = os.Remove("static/" + video.VideoPath)
	}

	return fileRemoveErr
}
