package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	dbconfig "github.com/aniket0951/Chatrapati-Maharaj/db-config"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var videoCategoryCollection *mongo.Collection = dbconfig.GetCollection(dbconfig.DB, "video_category")

type VideoRepository interface {
	CreateCategory(categories models.VideoCategories) (models.VideoCategories, error)
	UpdateCategory(categories models.VideoCategories) (models.VideoCategories, error)
	GetAllCategory() ([]models.VideoCategories, error)
	DeleteCategory(categoryId primitive.ObjectID) error
	DuplicateCategory(categoryName string) (bool, error)
	Init() (context.Context, context.CancelFunc)
}

type videocategoriesrepo struct {
	collection *mongo.Collection
}

func NewVideoCategoriesRepository() VideoRepository {
	return &videocategoriesrepo{
		collection: videoCategoryCollection,
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

	fmt.Println("id ==> ", categories.ID)

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
