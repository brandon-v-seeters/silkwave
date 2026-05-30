package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
)

const usersCollection = "Users"

var ErrUserCredentialsNotFound = errors.New("user credentials not found")

type UserCredentials struct {
	Password string `json:"password"`
	Key      string `json:"key"`
}

type UserRepository struct {
	db *database.ArangoDB
}

func NewUserRepository(db *database.ArangoDB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	col, err := r.db.GetCollection(ctx, usersCollection)
	if err != nil {
		return nil, err
	}

	user.CreatedAt = time.Now()

	meta, err := col.CreateDocument(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.Key = meta.Key
	user.ID = string(meta.ID)
	user.Rev = meta.Rev

	return user, nil
}

func (r *UserRepository) CreateWithCredentials(ctx context.Context, user *models.UserWithCredentials) (*models.UserWithCredentials, error) {
	col, err := r.db.GetCollection(ctx, usersCollection)
	if err != nil {
		return nil, err
	}

	user.CreatedAt = time.Now()

	meta, err := col.CreateDocument(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user with credentials: %w", err)
	}

	user.Key = meta.Key
	user.ID = string(meta.ID)
	user.Rev = meta.Rev

	return user, nil
}

func (r *UserRepository) GetCredentialsByEmail(ctx context.Context, email string) (*UserCredentials, error) {
	query := /*aql*/ `
		FOR user IN Users
		FILTER user.email == @email
		RETURN {
			password: user.password,
			key: user._key
		}
	`

	credentials, err := database.QueryOne[UserCredentials](ctx, r.db, query, map[string]interface{}{"email": email})
	if err != nil {
		if err.Error() == "no results found" {
			return nil, ErrUserCredentialsNotFound
		}

		return nil, fmt.Errorf("failed to get user credentials by email: %w", err)
	}

	return credentials, nil
}

func (r *UserRepository) GetByKey(ctx context.Context, key string) (*models.ClientUser, error) {
	q := /*aql*/ `
			FOR user IN Users
			FILTER user._key == @key
			LET artist = FIRST(
				FOR userArtist IN UsersArtists
				FILTER userArtist._from == CONCAT("Users/", user._key)
				RETURN DOCUMENT(userArtist._to)
			)
			RETURN MERGE(user, { artist: artist })
		`

	return database.QueryOne[models.ClientUser](ctx, r.db, q, map[string]interface{}{"key": key})
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := "FOR u IN Users RETURN u"

	users, err := database.Query[models.User](ctx, r.db, query, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, key string, user *models.User) (*models.User, error) {
	col, err := r.db.GetCollection(ctx, usersCollection)
	if err != nil {
		return nil, err
	}

	meta, err := col.UpdateDocument(ctx, key, user)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	user.Key = meta.Key
	user.ID = string(meta.ID)
	user.Rev = meta.Rev

	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, key string) error {
	col, err := r.db.GetCollection(ctx, usersCollection)
	if err != nil {
		return err
	}

	_, err = col.RemoveDocument(ctx, key)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
