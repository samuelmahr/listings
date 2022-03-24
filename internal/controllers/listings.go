package controllers

import (
	"github.com/gorilla/mux"
	"github.com/samuelmahr/listings/internal/configuration"
	"github.com/samuelmahr/listings/internal/models"
	"github.com/samuelmahr/listings/internal/repo"
	"net/http"
)

type V1ListingsController struct {
	config *configuration.AppConfig
	repo   repo.ListingsRepository
}

func NewV1AppointmentsController(c *configuration.AppConfig, aRepo repo.ListingsRepository) V1ListingsController {
	return V1ListingsController{
		config: c,
		repo:   aRepo,
	}
}

func (a *V1ListingsController) RegisterRoutes(v1 *mux.Router) {
	v1.Path("/listings").Name("GetListings").Handler(http.HandlerFunc(a.GetListings)).Methods(http.MethodGet)
}

func (a *V1ListingsController) GetListings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParams := r.URL.Query()

	listings, pagination, err := a.repo.GetListings(ctx, queryParams)
	if err != nil {
		respondError(ctx, w, 500, "error retrieving listings", err)
		return
	}

	response := models.ListingResponse{
		Type:       "FeatureCollection",
		Pagination: pagination,
	}

	features := make([]models.Feature, len(listings))
	for i, l := range listings {
		features[i] = l.ToFeature()
	}

	response.Features = features

	respondModel(ctx, w, http.StatusOK, response)
	return
}
