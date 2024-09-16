package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	EmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserById(ID int) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.AddUser(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

//mapping struck input ke user
//simpan struck user melalui repository

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("No user found on that email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) EmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}
	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {
	// find user by id
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, nil
	}
	// change attribute avatar file name
	user.AvatarFileName = fileLocation

	// save changes avatar file name
	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, nil
	}
	return updatedUser, nil
}

func (s *service) GetUserById(ID int) (User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, nil
	}

	if user.ID == 0 {
		return user, errors.New("No user found with that ID")
	}
	return user, nil
}
