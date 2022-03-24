package routers

import (
	"github.com/gorilla/mux"
	"github.com/samuelmahr/listings/internal/configuration"
	"github.com/samuelmahr/listings/internal/controllers"
	"github.com/samuelmahr/listings/internal/repo"
)

type V1Router struct {
	config *configuration.AppConfig
	uRepo  repo.ListingsRepoType
}

func NewV1Router(c *configuration.AppConfig, uRepo repo.ListingsRepoType) V1Router {
	return V1Router{config: c, uRepo: uRepo}
}

// Register initialize all routes
func (v *V1Router) Register(root *mux.Router) {
	r := root.PathPrefix("/v1").Subrouter()

	appointmentsController := controllers.NewV1AppointmentsController(v.config, &v.uRepo)
	appointmentsController.RegisterRoutes(r)
}
