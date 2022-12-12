package services

import (
	"errors"
	"log"

	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthService interface {
	CreateEndUser(user dto.RegisterEndUserDTO) (dto.RegisterEndUserDTO, error)
	CreateAdminUser(user dto.CreateAdminUserDTO) (dto.GetAdminUserDTO, error)
	ValidateAdminUser(email string, pass string) (dto.GetAdminUserDTO, error)
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

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)

	if err != nil {
		log.Println(err)
		return false
	}
	return true

}
