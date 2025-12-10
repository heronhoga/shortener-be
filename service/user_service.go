package service

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/heronhoga/shortener-be/model"
	"github.com/heronhoga/shortener-be/repository"
	"github.com/heronhoga/shortener-be/util/auth"
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

	// Validate phone number format
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	if !phoneRegex.MatchString(requests.Phone) {
		return errors.New("phone number format is not valid")
	}

	// Validate password length
	if len(requests.Password) < 8 {
		return errors.New("password length is less than 8")
	}

	// check if there is existing email/username
	available, err := s.repo.CheckExistingEmailUsername(ctx, requests.Email, requests.Username)
	if err != nil {
		return errors.New("There is existing user with the email/username")
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
		CreatedAt: time.Now().UTC(),
	}

	err = s.repo.InsertUser(ctx, newUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) LoginUser(ctx context.Context, requests *model.LoginUser) (string, error) {
	    switch requests.Provider {
    case "local":
        return s.loginLocal(ctx, requests.Email, requests.Password)
    case "google":
        return s.loginGoogle(ctx, requests.Token)
    default:
        return "", errors.New("unsupported login provider")
    }
}

func (s *UserService) loginLocal(ctx context.Context, email string, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)

	if err != nil {
		return "", errors.New("db error - getting existing user's data")
	}

	if user == nil {
		return "", errors.New("invalid email/password")
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email/password")
	}

	// generate token
	token, err := auth.GenerateToken(user.ID.String())
		if err != nil {
		return "", errors.New("error generating token")
	}

	return token, nil

}

func (s *UserService) loginGoogle(ctx context.Context, token string) (string, error) {
	return "", nil
}