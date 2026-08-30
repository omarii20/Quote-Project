package business

import (
	"context"
	"errors"
	"strings"
)

var ErrBusinessNotFound = errors.New("business not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateBusiness creates a new business in the database.
func (s *Service) CreateBusiness(ctx context.Context, b *Business) error {
	b.BusinessName = strings.TrimSpace(b.BusinessName)
	b.OwnerName = strings.TrimSpace(b.OwnerName)
	b.Phone = strings.TrimSpace(b.Phone)
	b.FirebaseUID = strings.TrimSpace(b.FirebaseUID)

	if b.BusinessName == "" {
		return errors.New("business name is required")
	}

	if b.OwnerName == "" {
		return errors.New("owner name is required")
	}

	if b.Phone == "" {
		return errors.New("phone is required")
	}

	if b.FirebaseUID == "" {
		return errors.New("firebase uid is required")
	}

	return s.repo.Create(ctx, b)
}

// GetBusiness retrieves a business by its ID from the database.
func (s *Service) GetBusiness(ctx context.Context, id int64) (*Business, error) {
	return s.repo.GetByID(ctx, id)
}
