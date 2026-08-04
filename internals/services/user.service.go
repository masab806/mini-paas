package services

import (
	"context"
	"errors"
	"mini-paas/ent"
	"mini-paas/internals/repositories"
	"mini-paas/internals/config"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService( repo repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, email string, username string, password string) (*ent.User, error) {

	exists, err := s.repo.ExistsByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("User Already Exists!")
	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if hashErr != nil {
		return nil, hashErr
	}

	newUser, userErr := s.repo.CreateUser(ctx, email, username, string(hashedPassword))

	if userErr != nil {
		return  nil , userErr
	}

	return newUser, nil

}

func (s *UserService) LoginUser(ctx context.Context, email string, password string) (string, error) {
	exists, err := s.repo.ExistsByEmail(ctx, email)

	if err != nil {
		return "", err
	}

	if !exists {
		return "", errors.New("Invalid Credentials")
	}

	user, userErr := s.repo.GetByEmail(ctx, email)

	if userErr != nil {
		return "", userErr
	}

	matchErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if matchErr != nil {
		return "", matchErr
	}

	token, tokenErr := config.GenerateToken(user.Email, user.ID, user.Username)

	if tokenErr != nil {
		return "", tokenErr
	}

	return token, nil

}

