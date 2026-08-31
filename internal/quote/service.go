package quote

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateQuote validates and calculates the pricing for a quote, then saves it to the database.
func (s *Service) CreateQuote(ctx context.Context, q *Quote) error {
	if q.BusinessID <= 0 {
		return errors.New("business id is required")
	}

	if q.CustomerID <= 0 {
		return errors.New("customer id is required")
	}

	// Check if the businessExists
	businessExists, err := s.repo.BusinessExists(ctx, q.BusinessID)
	if err != nil {
		return err
	}

	if !businessExists {
		return errors.New("business not found")
	}

	// Check if the customer belongs to the specified business
	customerBelongs, err := s.repo.CustomerBelongsToBusiness(ctx, q.CustomerID, q.BusinessID)
	if err != nil {
		return err
	}

	if !customerBelongs {
		return errors.New("customer not found or does not belong to business")
	}

	// Trim and validate the pricing method
	q.PricingMethod = strings.TrimSpace(q.PricingMethod)

	if q.PricingMethod == "" {
		q.PricingMethod = "items"
	}

	if q.PricingMethod != "items" && q.PricingMethod != "manual" {
		return errors.New("invalid pricing method")
	}

	if q.PricingMethod == "manual" && len(q.Items) > 0 {
		return errors.New("manual pricing cannot contain quote items")
	}

	if q.Title != nil {
		trimmed := strings.TrimSpace(*q.Title)
		q.Title = &trimmed
	}

	if q.Description != nil {
		trimmed := strings.TrimSpace(*q.Description)
		q.Description = &trimmed
	}

	if q.Notes != nil {
		trimmed := strings.TrimSpace(*q.Notes)
		q.Notes = &trimmed
	}

	if q.Status == "" {
		q.Status = "draft"
	}

	if q.PricingMethod == "items" {
		if err := calculateItemsPricing(q); err != nil {
			return err
		}
	}

	if q.PricingMethod == "manual" {
		if err := calculateManualPricing(q); err != nil {
			return err
		}
	}

	if err := calculateFinalPrice(q); err != nil {
		return err
	}

	return s.repo.Create(ctx, q)
}

// GetNextQuoteNumber retrieves the next available quote number for a given business.
func (s *Service) GetNextQuoteNumber(ctx context.Context, businessID int64) (string, error) {
	if businessID <= 0 {
		return "", errors.New("business id is required")
	}

	businessExists, err := s.repo.BusinessExists(ctx, businessID)
	if err != nil {
		return "", err
	}

	if !businessExists {
		return "", errors.New("business not found")
	}

	return s.repo.GetNextQuoteNumber(ctx, businessID)
}

// GetQuoteByID retrieves a quote by its ID.
func (s *Service) GetQuoteByID(ctx context.Context, quoteID int64) (*Quote, error) {
	if quoteID <= 0 {
		return nil, errors.New("quote id is required")
	}

	return s.repo.GetByID(ctx, quoteID)
}

// GetQuotesByBusinessID retrieves quotes for a specific business within a given period.
func (s *Service) GetQuotesByBusinessID(ctx context.Context, businessID int64, period string) ([]Quote, error) {

	if businessID <= 0 {
		return nil, errors.New("business id is required")
	}

	businessExists, err := s.repo.BusinessExists(ctx, businessID)
	if err != nil {
		return nil, err
	}

	if !businessExists {
		return nil, errors.New("business not found")
	}

	period = strings.TrimSpace(period)

	if period == "" {
		period = "month"
	}

	switch period {
	case "today", "week", "month", "all":
		// valid
	default:
		return nil, errors.New("invalid period")
	}

	return s.repo.GetByBusinessID(ctx, businessID, period)
}
