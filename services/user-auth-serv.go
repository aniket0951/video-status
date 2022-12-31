package services

import (
	"errors"
	"log"
	"time"

	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/mashingan/smapping"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthService interface {
	CreateEndUser(user dto.RegisterEndUserDTO) (dto.RegisterEndUserDTO, error)
	CreateAdminUser(user dto.CreateAdminUserDTO) (dto.GetAdminUserDTO, error)
	ValidateAdminUser(email string, pass string) (dto.GetAdminUserDTO, error)
	GetUserById(adminId primitive.ObjectID) (dto.GetAdminUserDTO, error)
	GetAllAdminUsers() ([]dto.GetAdminUserDTO, error)

	AddAdminUserAddress(address dto.CreateAdminUserAddress) error
	GetAdminUserAddress(userId primitive.ObjectID) (dto.GetAdminUserAddress, error)
	UpdateAdminAddress(address dto.UpdateAdminAddressDTO) error
}

type userauthservice struct {
	repo repositories.UserAuthRepository
}

func NewUserAuthService(repo repositories.UserAuthRepository) UserAuthService {
	return &userauthservice{
		repo: repo,
	}
}

func (ser *userauthservice) CreateEndUser(user dto.RegisterEndUserDTO) (dto.RegisterEndUserDTO, error) {

	newUserToCreate := models.Users{}

	smapping.FillStruct(&newUserToCreate, smapping.MapFields(&user))

	checkDuplicate := ser.repo.DuplicateMobile(newUserToCreate.MobileNumber)

	if !checkDuplicate {
		return dto.RegisterEndUserDTO{}, errors.New(helper.MOBILE_EXITS)
	}

	createdUser, err := ser.repo.CreateEndUser(newUserToCreate)

	if err != nil {
		return dto.RegisterEndUserDTO{}, err
	}

	smapping.FillStruct(&user, smapping.MapFields(&createdUser))

	return user, nil

}

func (ser *userauthservice) CreateAdminUser(user dto.CreateAdminUserDTO) (dto.GetAdminUserDTO, error) {
	userToCreate := models.AdminUser{}

	smpErr := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))

	if smpErr != nil {
		return dto.GetAdminUserDTO{}, smpErr
	}

	_, isEmailDuplicate := ser.repo.DuplicateEmail(user.Email)

	if !isEmailDuplicate {
		return dto.GetAdminUserDTO{}, errors.New(helper.EMAIL_EXITS)
	}

	newUser, err := ser.repo.CreateAdminUser(userToCreate)
	if err != nil {
		return dto.GetAdminUserDTO{}, err
	}

	newAdminUserDTO := dto.GetAdminUserDTO{}
	smapping.FillStruct(&newAdminUserDTO, smapping.MapFields(&newUser))

	return newAdminUserDTO, nil

}

func (ser *userauthservice) ValidateAdminUser(email string, pass string) (dto.GetAdminUserDTO, error) {
	adminUser, err := ser.repo.ValidateAdminUser(email)

	if err != nil {
		return dto.GetAdminUserDTO{}, err
	}
	if comparePassword(adminUser.Password, []byte(pass)) {
		adminUserDTO := dto.GetAdminUserDTO{}

		smapping.FillStruct(&adminUserDTO, smapping.MapFields(&adminUser))

		return adminUserDTO, nil
	}

	return dto.GetAdminUserDTO{}, errors.New("password not matched")
}

func (ser *userauthservice) GetUserById(adminId primitive.ObjectID) (dto.GetAdminUserDTO, error) {
	res, err := ser.repo.GetAdminUserById(adminId)

	if err != nil {
		return dto.GetAdminUserDTO{}, err
	}

	adminUser := dto.GetAdminUserDTO{}

	smapping.FillStruct(&adminUser, smapping.MapFields(&res))

	return adminUser, nil

}

func (ser *userauthservice) GetAllAdminUsers() ([]dto.GetAdminUserDTO, error) {
	users, err := ser.repo.GetAllAdminUsers()

	if err != nil {
		return []dto.GetAdminUserDTO{}, err
	}

	adminUsers := []dto.GetAdminUserDTO{}

	for i := range users {
		temp := dto.GetAdminUserDTO{}
		smapping.FillStruct(&temp, smapping.MapFields(users[i]))
		adminUsers = append(adminUsers, temp)
	}

	return adminUsers, nil
}

func (ser *userauthservice) AddAdminUserAddress(address dto.CreateAdminUserAddress) error {
	addressToCreate := models.AdminUserAddressInfo{}

	if smpErr := smapping.FillStruct(&addressToCreate, smapping.MapFields(address)); smpErr != nil {
		return smpErr
	}

	userObjId, objErr := primitive.ObjectIDFromHex(string(address.UserID.Hex()))

	if objErr != nil {
		return objErr
	}

	_, userAddErr := ser.repo.CheckUserAddressAlreadyAdded(userObjId)

	if userAddErr != nil {
		return userAddErr
	}

	addressToCreate.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	addressToCreate.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	err := ser.repo.AddAdminUserAddress(addressToCreate)

	return err
}

func (ser *userauthservice) GetAdminUserAddress(userId primitive.ObjectID) (dto.GetAdminUserAddress, error) {
	userAddress, err := ser.repo.CheckUserAddressAlreadyAdded(userId)

	if err == nil {
		return dto.GetAdminUserAddress{}, errors.New(helper.DATA_NOT_FOUND)
	}

	address := dto.GetAdminUserAddress{}

	smapping.FillStruct(&address, smapping.MapFields(&userAddress))

	return address, nil
}

func (ser *userauthservice) UpdateAdminAddress(address dto.UpdateAdminAddressDTO) error {
	addressToUpdate := models.AdminUserAddressInfo{}

	if smpErr := smapping.FillStruct(&addressToUpdate, smapping.MapFields(address)); smpErr != nil {
		return smpErr
	}

	upErr := ser.repo.UpdateAdminAddress(addressToUpdate)

	return upErr
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)

	if err != nil {
		log.Println(err)
		return false
	}
	return true

}
