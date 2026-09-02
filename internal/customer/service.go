package customer

import (
	"context"
	"errors"
	"strings"
)

var ErrCustomerNotFound = errors.New("customer not found")
var ErrBusinessNotFound = errors.New("business not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateCustomer creates a new customer in the database.
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

// GetCustomer retrieves a customer by its ID,
// only if it belongs to the given business.
func (s *Service) GetCustomer(
	ctx context.Context,
	id int64,
	businessID int64,
) (*Customer, error) {
	return s.repo.GetByID(ctx, id, businessID)
}

// GetCustomersByBusinessID retrieves all customers associated with a specific business ID.
func (s *Service) GetCustomersByBusinessID(
	ctx context.Context,
	businessID int64,
) ([]Customer, error) {

	if businessID <= 0 {
		return nil, errors.New("invalid business id")
	}

	exists, err := s.repo.BusinessExists(ctx, businessID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrBusinessNotFound
	}

	customers, err := s.repo.GetByBusinessID(ctx, businessID)
	if err != nil {
		return nil, err
	}

	if customers == nil {
		customers = []Customer{}
	}

	return customers, nil
}

// UpdateCustomer updates an existing customer's information,
// only if it belongs to the given business.
func (s *Service) UpdateCustomer(
	ctx context.Context,
	id int64,
	businessID int64,
	req *UpdateCustomerRequest,
) (*Customer, error) {

	if req.Name != nil {
		value := strings.TrimSpace(*req.Name)

		if value == "" {
			return nil, errors.New("customer name cannot be empty")
		}

		req.Name = &value
	}

	if req.Phone != nil {
		value := strings.TrimSpace(*req.Phone)

		if value == "" {
			return nil, errors.New("phone cannot be empty")
		}

		req.Phone = &value
	}

	return s.repo.Update(ctx, id, businessID, req)
}

// DeleteCustomer deletes a customer,
// only if it belongs to the given business.
func (s *Service) DeleteCustomer(
	ctx context.Context,
	id int64,
	businessID int64,
) error {
	return s.repo.Delete(ctx, id, businessID)
}
