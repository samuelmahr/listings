package repo

import (
	"context"
	"github.com/samuelmahr/listings/internal/models"
	"net/url"
)

// MockListings is an implementation of ListingsRepository to set values to use as a mock when testing
type MockListings struct {
	GetListingsResponse []models.Listing
	GetListingsErr      error
}

func (m *MockListings) GetListings(ctx context.Context, queryParams url.Values) ([]models.Listing, models.Pagination, error) {
	return m.GetListingsResponse, models.Pagination{}, m.GetListingsErr
}
