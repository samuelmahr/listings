package controllers

import (
	"context"
	"encoding/json"
	"github.com/samuelmahr/listings/internal/configuration"
	"github.com/samuelmahr/listings/internal/models"
	"github.com/samuelmahr/listings/internal/repo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

var aRepo repo.ListingsRepository
var config *configuration.AppConfig
var listingController V1ListingsController

func setup() {
	var err error
	config, err = configuration.Configure()
	if err != nil {
		panic("configuration error")
	}
}

func teardown() {

}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestV1Appointments_ListScheduledAppointments(t *testing.T) {
	type args struct {
		ctx   context.Context
		query url.Values
		aRepo repo.MockListings
	}

	tests := []struct {
		name     string
		args     args
		response int
		errMsg   string
	}{
		{
			name: "happy path min bed",
			args: args{
				ctx: context.TODO(),
				query: url.Values{
					"min_bed": []string{"1"},
				},
				aRepo: repo.MockListings{
					GetListingsResponse: []models.Listing{
						{
							ID:            1,
							Street:        "1434 Peace Dr",
							Status:        "active",
							Price:         85000,
							Bedrooms:      3,
							Bathrooms:     1,
							SquareFootage: 1000,
							Latitude:      38.4983765,
							Longitude:     -89.9696759,
							CreatedAt:     time.Now(),
							UpdatedAt:     time.Now(),
						},
					}},
			},
			response: http.StatusOK,
		},
	}

	endpoint := "/listings"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aRepo = &tt.args.aRepo

			listingController = NewV1AppointmentsController(config, aRepo)

			getHandler := http.HandlerFunc(listingController.GetListings)

			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				t.Fatal(err)
			}

			req.URL.RawQuery = tt.args.query.Encode()
			response := httptest.NewRecorder()
			getHandler.ServeHTTP(response, req)
			assert.Equal(t, tt.response, response.Code)

			if tt.response != http.StatusOK {
				resp := make(map[string]string)
				err = json.Unmarshal(response.Body.Bytes(), &resp)
				assert.Equal(t, tt.errMsg, resp["error"])
			}
		})
	}
}
