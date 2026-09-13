package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/synclabs-io/GateForge/internal/domain/entities"
	"github.com/synclabs-io/GateForge/internal/domain/value_objects"
	"github.com/synclabs-io/GateForge/internal/repositories/redis/dto"
	dto2 "github.com/synclabs-io/GateForge/internal/service/dto"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, entity entities.User) (entities.User, error)
}

type SessionRepository interface {
	Save(ctx context.Context, input dto.SaveRequestDTO) (dto.SaveResponseDTO, error)
}

type JWTManager interface {
	GenerateToken(userID *value_objects.ID) (string, error)
}

type UsersService struct {
	repo       UsersRepository
	cache      SessionRepository
	jwt        JWTManager
	hash       value_objects.PasswordHash
	logger     *slog.Logger
	refreshTTL time.Duration
}

func NewUsersService(repo UsersRepository, cache SessionRepository, logger *slog.Logger, jwt JWTManager, hash value_objects.PasswordHash, refreshTTL time.Duration) *UsersService {
	return &UsersService{
		repo:       repo,
		cache:      cache,
		jwt:        jwt,
		hash:       hash,
		logger:     logger,
		refreshTTL: refreshTTL,
	}
}

func (s *UsersService) Register(ctx context.Context, input dto2.RegisterRequestDTO) (dto2.RegisterResponseDTO, error) {
	user, err := entities.NewUser(input.Username, input.Password, s.hash)
	if err != nil {
		s.logger.ErrorContext(ctx, "create user entity", slog.String("error", err.Error()))
		return dto2.RegisterResponseDTO{}, err
	}

	_, err = s.repo.CreateUser(ctx, *user)
	if err != nil {
		s.logger.ErrorContext(ctx, "create user", slog.String("error", err.Error()))
		return dto2.RegisterResponseDTO{}, err
	}

	resp, err := s.cache.Save(ctx, dto.SaveRequestDTO{
		UserID: user.ID,
		TTL:    s.refreshTTL,
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "cache save", slog.String("error", err.Error()))
		return dto2.RegisterResponseDTO{}, err
	}

	access, err := s.jwt.GenerateToken(user.ID)
	if err != nil {
		s.logger.ErrorContext(ctx, "generate jwt", slog.String("error", err.Error()))
		return dto2.RegisterResponseDTO{}, err
	}

	return dto2.RegisterResponseDTO{
		Success: true,
		Access:  access,
		Refresh: resp.Session.Refresh.Raw(),
	}, nil
}
