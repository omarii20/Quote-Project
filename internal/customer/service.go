package customer

import (
	"context"
	"errors"
	"strings"
)

var ErrCustomerNotFound = errors.New("customer not found")
var ErrBusinessNotFound = errors.New("business not found")
var ErrCustomerHasQuotes = errors.New("customer has existing quotes")

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

<<<<<<< Updated upstream
// GetCustomer retrieves a customer by its ID from the database.
func (s *Service) GetCustomer(ctx context.Context, id int64) (*Customer, error) {
	return s.repo.GetByID(ctx, id)
}

// Get CustomersByBusinessID retrieves all customers associated with a specific business ID.
=======
// GetCustomer retrieves a customer by its ID,
// only if it belongs to the given business.
func (s *Service) GetCustomer(ctx context.Context, id int64, businessID int64) (*Customer, error) {
	return s.repo.GetByID(ctx, id, businessID)
}

// GetCustomersByBusinessID retrieves all customers associated with a specific business ID.
>>>>>>> Stashed changes
func (s *Service) GetCustomersByBusinessID(ctx context.Context, businessID int64) ([]Customer, error) {

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

<<<<<<< Updated upstream
// UpdateCustomer updates an existing customer's information in the database.
func (s *Service) UpdateCustomer(ctx context.Context, id int64, req *UpdateCustomerRequest) (*Customer, error) {

=======
// UpdateCustomer updates an existing customer's information,
// only if it belongs to the given business.
func (s *Service) UpdateCustomer(ctx context.Context, id int64, businessID int64, req *UpdateCustomerRequest) (*Customer, error) {
>>>>>>> Stashed changes
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

	return s.repo.Update(ctx, id, req)
}

<<<<<<< Updated upstream
// DeleteCustomer deletes a customer from the database by its ID.
func (s *Service) DeleteCustomer(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
=======
// DeleteCustomer deletes a customer,
// only if it belongs to the given business
// and has no existing quotes.
func (s *Service) DeleteCustomer(ctx context.Context, id int64, businessID int64) error {
	hasQuotes, err := s.repo.HasQuotes(ctx, id, businessID)
	if err != nil {
		return err
	}

	if hasQuotes {
		return ErrCustomerHasQuotes
	}

	return s.repo.Delete(ctx, id, businessID)
>>>>>>> Stashed changes
}
