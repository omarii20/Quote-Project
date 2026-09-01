package quote

import "errors"

func validateQuoteStatus(status string) error {
	switch status {
	case "draft", "sent", "viewed", "approved", "rejected", "expired":
		return nil
	default:
		return errors.New("invalid quote status")
	}
}
