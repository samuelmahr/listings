package models

import (
	"math"
	"time"
)

// Listing models database table
type Listing struct {
	ID            int64     `json:"id" db:"id"`
	Street        string    `json:"street" db:"street"`
	Status        string    `json:"status" db:"status"`
	Price         int       `json:"price" db:"price"`
	Bedrooms      int       `json:"bedrooms" db:"bedrooms"`
	Bathrooms     int       `json:"bathrooms" db:"bathrooms"`
	SquareFootage int       `json:"sq_ft" db:"sq_ft"`
	Latitude      float64   `json:"latitude" db:"lat"`
	Longitude     float64   `json:"longitude" db:"lng"`
	CreatedAt     time.Time `json:"-" db:"created_at"`
	UpdatedAt     time.Time `json:"-" db:"updated_at"`
}

func (l *Listing) ToFeature() Feature {
	return Feature{
		Type: "features",
		Geometry: Geometry{
			Type:        "Point",
			Coordinates: []float64{toFixed(l.Longitude, 1), toFixed(l.Latitude, 1)},
		},
		Properties: Properties{
			ID:            l.ID,
			Street:        l.Street,
			Status:        l.Status,
			Price:         l.Price,
			Bedrooms:      l.Bedrooms,
			Bathrooms:     l.Bathrooms,
			SquareFootage: l.SquareFootage,
		},
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

// Geometry nested object in API response body
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// Properties nested object in API response body
type Properties struct {
	ID            int64  `json:"id"`
	Street        string `json:"street"`
	Status        string `json:"status"`
	Price         int    `json:"price"`
	Bedrooms      int    `json:"bedrooms"`
	Bathrooms     int    `json:"bathrooms"`
	SquareFootage int    `json:"sq_ft"`
}

// Feature each listing item
type Feature struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

// ListingResponse response to match API contract
type ListingResponse struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
	Pagination
}

type Pagination struct {
	Page     uint64 `json:"page"`
	PageSize uint64 `json:"page_size"`
}
