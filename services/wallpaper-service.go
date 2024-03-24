package services

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"path"
	"reflect"
	"strings"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	notificationmanager "github.com/aniket0951/Chatrapati-Maharaj/notification_manager"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WallPaperService interface {
	AddWallPaper(file multipart.File, wallPaper models.WallPaper) error
	GetWallPapers(isActive bool) (dto.GetWallPapersDTO, error)
	ActiveInActiveWallPaper(videoId string, isActive bool) error
	WallPaperLiked(wallPaperId string) error
	FetchRecentWallPapers(isActive bool) ([]models.WallPaper, error)
}

type service struct {
	wallPaperRepo      repositories.WallPaperRepository
	notifcationService notificationmanager.NotificationManager
}

func NewWallPaperService(repo repositories.WallPaperRepository, notificationManager notificationmanager.NotificationManager) WallPaperService {
	return &service{
		wallPaperRepo:      repo,
		notifcationService: notificationManager,
	}
}

func (serv *service) AddWallPaper(file multipart.File, wallPaper models.WallPaper) error {
	tempFile, err := ioutil.TempFile("static/wallpaper", "wallpaper-*.png")

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

	wallPaper.FilePath = path.Base(tempFile.Name())

	filePath := path.Base(tempFile.Name())

	wallPaper.FilePath = filePath
	wallPaper.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	wallPaper.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	wallPaper.Category = helper.RECENT

	return serv.wallPaperRepo.AddWallPaper(wallPaper)
}

func (serv *service) GetWallPapers(isActive bool) (dto.GetWallPapersDTO, error) {
	result, err := serv.wallPaperRepo.GetWallPapers(isActive)

	if err != nil {
		return dto.GetWallPapersDTO{}, err
	}

	if len(result) == 0 {
		return dto.GetWallPapersDTO{}, errors.New("video not found")
	}

	wallPaper_data := dto.GetWallPapersDTO{
		Recent: []models.WallPaper{},
		Olds:   []models.WallPaper{},
	}

	for i := range result {
		path := result[i].FilePath

		if strings.Contains(path, "static") {
			result[i].FilePath = "http://localhost:5000/" + path
		} else {
			result[i].FilePath = "http://localhost:5000/static/wallpaper/" + path
		}

		if result[i].Category == helper.RECENT {
			wallPaper_data.Recent = append(wallPaper_data.Recent, result[i])
		} else {
			wallPaper_data.Olds = append(wallPaper_data.Olds, result[i])
		}
	}

	if len(wallPaper_data.Recent) == 0 {
		wallPaper_data.Recent = append(wallPaper_data.Recent, wallPaper_data.Olds[:len(wallPaper_data.Olds)/2]...)
	}

	return wallPaper_data, err
}

func (serv *service) FetchRecentWallPapers(isActive bool) ([]models.WallPaper, error) {
	result, err := serv.wallPaperRepo.FetchRecentWallPapers(isActive)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("recent WallPaper not found")
	}

	for i := range result {
		path := result[i].FilePath

		if strings.Contains(path, "static") {
			result[i].FilePath = "http://localhost:5000/" + path
		} else {
			result[i].FilePath = "http://localhost:5000/static/wallpaper/" + path
		}
	}

	return result, err
}

func (serv *service) ActiveInActiveWallPaper(videoId string, isActive bool) error {
	objId, err := primitive.ObjectIDFromHex(videoId)

	if err != nil {
		return err
	}
	if err := serv.wallPaperRepo.ActiveInActiveWallPaper(objId, isActive); err != nil {
		return err
	}

	// check if isActive is true then unset the recent tag for first one
	if isActive {
		// trigger the notification from here
		// for this fetch wallpaper title
		var description string
		wallPaperData, err := serv.wallPaperRepo.GetWallPaper(objId)
		if err != nil || reflect.DeepEqual(wallPaperData, &models.WallPaper{}) {
			description = "New WallPaper has been uploaded!"
		} else {
			description = wallPaperData.Title
		}

		notificationMessage := messaging.Message{
			Notification: &messaging.Notification{
				Title:    "Jay Bhavai !",
				Body:     description,
				ImageURL: "",
			},
		}

		// notificationService := notificationmanager.NotificationManager{}
		// notificationService.NotifyAllUser(&notificationMessage)

		serv.notifcationService.NotifyAllUser(&notificationMessage)
		if err := serv.wallPaperRepo.UnsetRecentCategory(); err != nil {
			fmt.Println("UpdateRecent Category Error : ", err)
		}
		return nil
	}

	return nil
}

func (serv *service) WallPaperLiked(wallPaperId string) error {

	objId, err := primitive.ObjectIDFromHex(wallPaperId)

	if err != nil {
		return err
	}

	return serv.wallPaperRepo.WallPaperLiked(objId)
}
