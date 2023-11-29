package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/commedesvlados/grpc_service/internal/domain/models"
	"github.com/commedesvlados/grpc_service/internal/lib/jwt"
	"github.com/commedesvlados/grpc_service/internal/lib/logger/sl"
	"github.com/commedesvlados/grpc_service/internal/storage"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
)

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) RegisterUser(ctx context.Context, email, password string) (int64, error) {
	const fn = "auth.RegisterUser"

	log := a.log.With(
		slog.String("fn", fn),
		slog.String("email", email),
	)

	log.Info("register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("failed to save user", sl.Err(err))
			return 0, fmt.Errorf("%s: %w", fn, ErrUserExists)
		}

		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("user is registered")
	return id, nil
}

func (a *Auth) LoginUser(ctx context.Context, email, password string, appID int) (string, error) {
	const fn = "auth.LoginUser"

	log := a.log.With(
		slog.String("fn", fn),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", fn, ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", fn, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		log.Error("failed to get app", sl.Err(err))
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return token, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const fn = "auth.IsAdmin"

	log := a.log.With(
		slog.String("fn", fn),
		slog.Int64("userID", userID),
	)

	log.Info("checking if user if admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Error("failed to check if user is admin", sl.Err(err))
			return false, fmt.Errorf("%s: %w", fn, ErrInvalidAppID)
		}

		log.Error("failed to check if user is admin", sl.Err(err))
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
