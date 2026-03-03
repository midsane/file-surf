package user

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/midsane/file-surf/internal/tenant"
)

type Service struct {
	repo       *Repository
	tenantRepo *tenant.Repository
}

func NewService(repo *Repository, tenantRepo *tenant.Repository) *Service {
	return &Service{
		repo:       repo,
		tenantRepo: tenantRepo,
	}
}

func (s *Service) CreateUser(ctx context.Context, tenantID, email string) (*User, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, errors.New("email is required")
	}

	// Check tenant exists
	t, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.New("tenant not found")
	}

	user := &User{
		ID:       uuid.New().String(),
		TenantID: tenantID,
		Email:    email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUsers(ctx context.Context, tenantID string) ([]User, error) {
	return s.repo.GetByTenant(ctx, tenantID)
}