package customer

import (
	"context"
	"errors"
	"strings"
)

var ErrCustomerNotFound = errors.New("customer not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateCustomer(ctx context.Context, c *Customer) error {
	c.Name = strings.TrimSpace(c.Name)
	c.Phone = strings.TrimSpace(c.Phone)

	if c.BusinessID <= 0 {
		return errors.New("business id is required")
	}

	if c.Name == "" {
		return errors.New("customer name is required")
	}

	if c.Phone == "" {
		return errors.New("phone is required")
	}

	return s.repo.Create(ctx, c)
}

func (s *Service) GetCustomer(ctx context.Context, id int64) (*Customer, error) {
	return s.repo.GetByID(ctx, id)
}
