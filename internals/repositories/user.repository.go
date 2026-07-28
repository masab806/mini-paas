package repositories

import (
	"context"
	"mini-paas/ent"
	"mini-paas/ent/user"
)

type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, email string, username string, password string) (*ent.User, error)
}

type userRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) UserRepository {
	return &userRepository{
		client: client,
	}
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.client.User.Query().Where(user.Email(email)).Exist(ctx)
}

func (r *userRepository) CreateUser(ctx context.Context, email string, username string, password string) (*ent.User, error) {
	return r.client.User.Create().SetEmail(email).SetUserName(username).SetPassword(password).Save(ctx)
}