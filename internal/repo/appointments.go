package repo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/samuelmahr/listings/internal/models"
	"net/url"
	"strings"
)

type ListingsRepository interface {
	GetListings(ctx context.Context, queryParams url.Values) ([]models.Listing, error)
}

type ListingsRepoType struct {
	db *sqlx.DB
}

func NewListingsRepository(db *sqlx.DB) ListingsRepoType {
	return ListingsRepoType{
		db: db,
	}
}

const (
	FilterGTE = "GTE"
	FilterLTE = "LTE"
)

type FilterBuildParams struct {
	// Filters that narrow down list
	AllowedFilters map[string]FilterFields
}

type FilterFields struct {
	// Column name
	Field string
	// Filter type LIKE/=/etc
	FilterType string
}

func (ar *ListingsRepoType) GetListings(ctx context.Context, queryParams url.Values) ([]models.Listing, error) {
	filterParams := FilterBuildParams{
		AllowedFilters: map[string]FilterFields{
			"min_price": {
				Field:      "price",
				FilterType: FilterGTE,
			},
			"max_price": {
				Field:      "price",
				FilterType: FilterLTE,
			},
			"min_bed": {
				Field:      "bedrooms",
				FilterType: FilterGTE,
			},
			"max_bed": {
				Field:      "bedrooms",
				FilterType: FilterLTE,
			},
			"min_bath": {
				Field:      "bathrooms",
				FilterType: FilterGTE,
			},
			"max_bath": {
				Field:      "bathrooms",
				FilterType: FilterLTE,
			},
		},
	}

	listQueryTemplate := sq.Select("id", "street", "status", "price", "bedrooms", "bathrooms", "sq_ft", "lat", "lng").From("features.listings")
	listQueryFiltered := BuildFilterQuery(queryParams, filterParams, listQueryTemplate)

	sql, args, err := listQueryFiltered.ToSql()
	if err != nil {
		return []models.Listing{}, errors.Wrap(err, "error getting listings")
	}

	rows, err := ar.db.Queryx(sql, args...)
	if err != nil {
		return []models.Listing{}, errors.Wrap(err, "error getting listings")
	}

	appts := make([]models.Listing, 0)
	for rows.Next() {
		var a models.Listing
		if err := rows.StructScan(&a); err != nil {
			return []models.Listing{}, errors.Wrap(err, "error getting appointments")
		}

		appts = append(appts, a)
	}

	return appts, nil
}

func BuildFilterQuery(queryParams url.Values, filterParams FilterBuildParams, builder sq.SelectBuilder) sq.SelectBuilder {
	for query, values := range queryParams {
		// If endpoint query is an allowed query parameter
		var unpackedValues []string
		for _, v := range values {
			if v == "" {
				continue
			}
			parsed := strings.Split(v, ",")
			unpackedValues = append(unpackedValues, parsed...)
		}

		if len(unpackedValues) == 0 {
			continue
		}

		if filterField, ok := filterParams.AllowedFilters[query]; ok {
			field := filterParams.AllowedFilters[query].Field

			switch filterField.FilterType {
			case FilterLTE:
				builder = builder.Where(sq.LtOrEq{field: unpackedValues[0]})
			case FilterGTE:
				builder = builder.Where(sq.GtOrEq{field: unpackedValues[0]})
			}
		}
	}

	return builder.PlaceholderFormat(sq.Dollar)
}
