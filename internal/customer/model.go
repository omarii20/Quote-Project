package customer

import "time"

type Customer struct {
	ID         int64     `json:"id"`
	BusinessID int64     `json:"business_id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Email      *string   `json:"email,omitempty"`
	Address    *string   `json:"address,omitempty"`
	Notes      *string   `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpdateCustomerRequest struct {
	Name    *string `json:"name"`
	Phone   *string `json:"phone"`
	Email   *string `json:"email"`
	Address *string `json:"address"`
	Notes   *string `json:"notes"`
}
