package service

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/heronhoga/shortener-be/model"
	"github.com/heronhoga/shortener-be/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterNewUser(ctx context.Context, requests *model.RegisterUser) error {
	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(requests.Email) {
		return errors.New("email format is not valid")
	}

	// Validate password length
	if len(requests.Password) < 8 {
		return errors.New("password length is less than 8")
	}

	// check if there is existing email/username
	available, err := s.repo.CheckExistingEmailUsername(ctx, requests.Email, requests.Username)
	if err != nil {
		return errors.New("Internal server error - query check existing email and username")
	}

	if !available {
		return errors.New("There is existing user with the email/username")
	}

	// generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requests.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Error generating password")
	}

	// create id for new user
	userId, err := uuid.NewV7()
	if err != nil {
		return errors.New("Error creating id for new user")
	}
	//create new User
	newUser := &model.User{
		ID: userId,
		Email: requests.Email,
		Username: requests.Username,
		Password: string(hashedPassword),
		Phone: requests.Phone,
		CreatedAt: time.Now(),
	}

	err = s.repo.InsertUser(ctx, newUser)
	if err != nil {
		return err
	}

	return nil
}