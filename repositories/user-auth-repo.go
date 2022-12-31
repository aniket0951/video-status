package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	dbconfig "github.com/aniket0951/Chatrapati-Maharaj/db-config"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthRepository interface {
	CreateEndUser(user models.Users) (models.Users, error)
	CreateAdminUser(adminuser models.AdminUser) (models.AdminUser, error)
	GetAdminUserById(adminId primitive.ObjectID) (models.AdminUser, error)
	UpdateAdminUserInfo(adminuser models.AdminUser) error
	GetAllAdminUsers() ([]models.AdminUser, error)
	ValidateAdminUser(email string) (models.AdminUser, error)
	DuplicateMobile(mobile string) bool
	DuplicateEmail(email string) (models.AdminUser, bool)

	AddAdminUserAddress(address models.AdminUserAddressInfo) error
	CheckUserAddressAlreadyAdded(userId primitive.ObjectID) (models.AdminUserAddressInfo, error)
	UpdateAdminAddress(address models.AdminUserAddressInfo) error

	Init() (context.Context, context.CancelFunc)
}

var userAddressCollection *mongo.Collection = dbconfig.GetCollection(dbconfig.DB, "user_address")

type userauthrepository struct {
	userconnection        *mongo.Collection
	userAddressConnection *mongo.Collection
}

func NewUserAuthRepository(userCollection *mongo.Collection) UserAuthRepository {
	return &userauthrepository{
		userconnection:        userCollection,
		userAddressConnection: userAddressCollection,
	}
}

func (db *userauthrepository) Init() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, cancel
}

func (db *userauthrepository) CreateEndUser(user models.Users) (models.Users, error) {
	ctx, cancel := db.Init()
	defer cancel()
	currentTime := time.Now()
	currentTime.Format("2006.01.02 15:04:05")

	user.ID = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(currentTime)
	user.UpdatedAt = primitive.NewDateTimeFromTime(currentTime)

	if user.UserType == "" {
		user.UserType = "end_user"
	}

	result, err := db.userconnection.InsertOne(ctx, &user)
	if err != nil {
		return models.Users{}, err
	}

	var newUser models.Users

	filter := bson.D{
		bson.E{Key: "_id", Value: result.InsertedID},
	}

	db.userconnection.FindOne(ctx, filter).Decode(&newUser)

	return newUser, nil
}

func (db *userauthrepository) CreateAdminUser(adminuser models.AdminUser) (models.AdminUser, error) {
	ctx, cancel := db.Init()
	defer cancel()

	adminuser.ID = primitive.NewObjectID()
	adminuser.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	adminuser.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	adminuser.Password = hasAndSalt([]byte(adminuser.Password))
	adminuser.UserType = "admin_user"

	_, insError := db.userconnection.InsertOne(ctx, &adminuser)

	if insError != nil {
		return models.AdminUser{}, insError
	}

	return db.GetAdminUserById(adminuser.ID)
}

func (db *userauthrepository) GetAdminUserById(adminId primitive.ObjectID) (models.AdminUser, error) {

	filter := bson.D{
		bson.E{Key: "_id", Value: adminId},
	}

	//pipline := []bson.M{
	//	bson.M{
	//		"$match": bson.M{
	//			"_id": adminId,
	//		},
	//	},
	//	bson.M{
	//		"$lookup": bson.M{
	//			"from": "user_address",
	//			"localfield": "_id",
	//			"foreignField": "user"
	//		},
	//	},
	//}

	ctx, cancel := db.Init()
	defer cancel()
	var adminUser models.AdminUser

	db.userconnection.FindOne(ctx, filter).Decode(&adminUser)
	fmt.Println("admin user ==> ", adminUser)
	if (adminUser == models.AdminUser{}) {
		return models.AdminUser{}, errors.New(helper.DATA_NOT_FOUND)
	}

	return adminUser, nil
}

func (db *userauthrepository) UpdateAdminUserInfo(adminuser models.AdminUser) error {
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "username", Value: adminuser.UserName},
			bson.E{Key: "contact", Value: adminuser.MobileNumber},
			bson.E{Key: "email", Value: adminuser.Email},
			bson.E{Key: "user_type", Value: adminuser.UserType},
			bson.E{Key: "updated_at", Value: primitive.NewDateTimeFromTime(time.Now())},
		}},
	}

	ctx, cancel := db.Init()
	defer cancel()

	res, err := db.userconnection.UpdateByID(ctx, adminuser.ID, update)

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("user not found for update")
	}

	return nil
}

func (db *userauthrepository) GetAllAdminUsers() ([]models.AdminUser, error) {
	ctx, cancel := db.Init()
	defer cancel()

	allAdminUsers := []models.AdminUser{}

	cursor, curErr := db.userconnection.Find(ctx, bson.M{})

	if curErr != nil {
		return []models.AdminUser{}, curErr
	}

	if err := cursor.All(context.TODO(), &allAdminUsers); err != nil {
		return []models.AdminUser{}, err
	}

	return allAdminUsers, nil

}

func (db *userauthrepository) ValidateAdminUser(email string) (models.AdminUser, error) {
	user, res := db.DuplicateEmail(email)
	fmt.Println(user, res)
	if !res {
		return user, nil
	}

	return models.AdminUser{}, errors.New(helper.DATA_NOT_FOUND)
}

func (db *userauthrepository) DuplicateMobile(mobile string) bool {
	ctx, cancel := db.Init()
	defer cancel()
	filter := bson.D{
		bson.E{Key: "contact", Value: mobile},
	}

	var user models.Users
	res := db.userconnection.FindOne(ctx, filter).Decode(&user)

	return res == mongo.ErrNoDocuments

}

func (db *userauthrepository) DuplicateEmail(email string) (models.AdminUser, bool) {
	filter := bson.D{
		bson.E{Key: "email", Value: email},
	}

	ctx, cancel := db.Init()
	defer cancel()

	var adminUser models.AdminUser

	res := db.userconnection.FindOne(ctx, filter).Decode(&adminUser)
	return adminUser, res == mongo.ErrNoDocuments
}

func (db *userauthrepository) AddAdminUserAddress(address models.AdminUserAddressInfo) error {
	ctx, cancel := db.Init()
	defer cancel()

	_, err := db.userAddressConnection.InsertOne(ctx, &address)

	if err != nil {
		return err
	}
	return nil
}

func (db *userauthrepository) CheckUserAddressAlreadyAdded(userId primitive.ObjectID) (models.AdminUserAddressInfo, error) {
	filter := bson.D{
		bson.E{Key: "userId", Value: userId},
	}

	ctx, cancel := db.Init()
	defer cancel()

	userAddress := models.AdminUserAddressInfo{}

	db.userAddressConnection.FindOne(ctx, filter).Decode(&userAddress)

	if (userAddress == models.AdminUserAddressInfo{}) {
		return models.AdminUserAddressInfo{}, nil
	}

	return userAddress, errors.New("user address already exits, please try to update address")
}

func (db *userauthrepository) UpdateAdminAddress(address models.AdminUserAddressInfo) error {
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "state", Value: address.State},
			bson.E{Key: "city", Value: address.City},
			bson.E{Key: "address", Value: address.Addressline},
			bson.E{Key: "updated_at", Value: primitive.NewDateTimeFromTime(time.Now())},
		}},
	}

	ctx, cancel := db.Init()
	defer cancel()

	res, err := db.userAddressConnection.UpdateByID(ctx, address.ID, update)

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("address not found for update")
	}

	return nil

}

func hasAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to has a password")
	}
	return string(hash)
}
