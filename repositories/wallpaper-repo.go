package repositories

import (
	"context"
	"errors"
	"time"

	dbconfig "github.com/aniket0951/Chatrapati-Maharaj/db-config"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var WallPaperCollection *mongo.Collection = dbconfig.GetCollection(dbconfig.DB, "wallpaper")

type WallPaperRepository interface {
	AddWallPaper(wallPaper models.WallPaper) error
	GetWallPapers(isActive bool) ([]models.WallPaper, error)
	ActiveInActiveWallPaper(videoId primitive.ObjectID, isActive bool) error
	WallPaperLiked(wallpapper_id primitive.ObjectID) error
	FetchRecentWallPapers(isActive bool) ([]models.WallPaper, error)
	GetWallPaper(wallPaperId primitive.ObjectID) (models.WallPaper, error)

	UnsetRecentCategory() error
}

type repo struct {
	WallPaperCollection *mongo.Collection
}

func NewWallPaperRepository() WallPaperRepository {
	return &repo{
		WallPaperCollection: WallPaperCollection,
	}
}

func (repo *repo) Init() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx, cancel
}

func (repo *repo) AddWallPaper(wallPaper models.WallPaper) error {
	ctx, cancel := repo.Init()
	defer cancel()

	_, err := repo.WallPaperCollection.InsertOne(ctx, wallPaper)
	return err
}

func (repo *repo) GetWallPapers(isActive bool) ([]models.WallPaper, error) {
	ctx, cancel := repo.Init()
	defer cancel()

	filter := bson.M{
		"is_active": isActive,
	}
	sort := options.Find().SetSort(bson.M{"updated_at": -1})
	cursor, err := repo.WallPaperCollection.Find(ctx, filter, sort)

	if err != nil {
		return nil, err
	}

	var result []models.WallPaper

	if err := cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *repo) FetchRecentWallPapers(isActive bool) ([]models.WallPaper, error) {
	ctx, cancel := repo.Init()
	defer cancel()
	filter := bson.M{
		"updated_at": bson.M{
			"$gte": time.Now().AddDate(0, 0, -10),
		},
		"is_active": isActive,
	}

	cursor, err := repo.WallPaperCollection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var result = []models.WallPaper{}

	if err := cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// update the wallpaper status
func (repo *repo) ActiveInActiveWallPaper(videoId primitive.ObjectID, isActive bool) error {
	ctx, cancel := repo.Init()
	defer cancel()

	filter := bson.M{"_id": videoId}

	update := bson.M{
		"$set": bson.M{
			"is_active":  isActive,
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := repo.WallPaperCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("video not found to update")
	}

	return nil
}

// unset a category for first in the list of previous recent wallpaper
func (repo *repo) UnsetRecentCategory() error {
	ctx, cancel := repo.Init()
	defer cancel()

	filter := bson.M{
		"category": "recent",
	}

	opt := options.FindOneAndUpdate()
	opt.SetSort(bson.M{"_id": 1})

	update := bson.M{
		"$set": bson.M{
			"category":   "",
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	res := repo.WallPaperCollection.FindOneAndUpdate(ctx, filter, update, opt)

	return res.Err()
}

func (db *repo) WallPaperLiked(wallpapper_id primitive.ObjectID) error {
	ctx, cancel := db.Init()
	defer cancel()

	update := bson.M{
		"$inc": bson.M{
			"download_count": 1,
		},
	}

	filter := bson.M{"_id": wallpapper_id}

	_, err := db.WallPaperCollection.UpdateOne(ctx, filter, update)

	return err
}

func (db *repo) GetWallPaper(wallPaperId primitive.ObjectID) (models.WallPaper, error) {
	ctx, cancel := db.Init()
	defer cancel()

	filter := bson.D{
		bson.E{Key: "_id", Value: wallPaperId},
	}

	wallPaper := models.WallPaper{}
	err := db.WallPaperCollection.FindOne(ctx, filter).Decode(&wallPaper)
	return wallPaper, err
}
