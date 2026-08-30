package business

import "time"

type Business struct {
	ID           int64     `json:"id"`
	FirebaseUID  string    `json:"firebase_uid"`
	BusinessName string    `json:"business_name"`
	OwnerName    string    `json:"owner_name"`
	Phone        string    `json:"phone"`
	Email        *string   `json:"email,omitempty"`
	Address      *string   `json:"address,omitempty"`
	LogoURL      *string   `json:"logo_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UpdateBusinessRequest struct {
	BusinessName *string `json:"business_name"`
	OwnerName    *string `json:"owner_name"`
	Phone        *string `json:"phone"`
	Email        *string `json:"email"`
	Address      *string `json:"address"`
	LogoURL      *string `json:"logo_url"`
}
