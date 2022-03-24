package repo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/samuelmahr/listings/internal/models"
	"net/url"
	"strconv"
	"strings"
)

type ListingsRepository interface {
	GetListings(ctx context.Context, queryParams url.Values) ([]models.Listing, models.Pagination, error)
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

const (
	defaultPage     = 1
	defaultPageSize = 10
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

func (ar *ListingsRepoType) GetListings(ctx context.Context, queryParams url.Values) ([]models.Listing, models.Pagination, error) {
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
	page, limit, offset, err := BuildPagination(queryParams)
	if err != nil {
		return nil, models.Pagination{}, err
	}

	listQueryTemplate := sq.Select("id", "street", "status", "price", "bedrooms", "bathrooms", "sq_ft", "lat", "lng").From("features.listings")
	listQueryFiltered := BuildFilterQuery(queryParams, filterParams, listQueryTemplate)
	listQueryPaginated := listQueryFiltered.Limit(limit).Offset(offset)

	sql, args, err := listQueryPaginated.ToSql()
	if err != nil {
		return []models.Listing{}, models.Pagination{}, errors.Wrap(err, "error getting listings")
	}

	rows, err := ar.db.Queryx(sql, args...)
	if err != nil {
		return []models.Listing{}, models.Pagination{}, errors.Wrap(err, "error getting listings")
	}

	listings := make([]models.Listing, 0)
	for rows.Next() {
		var a models.Listing
		if err := rows.StructScan(&a); err != nil {
			return []models.Listing{}, models.Pagination{}, errors.Wrap(err, "error getting appointments")
		}

		listings = append(listings, a)
	}

	return listings, models.Pagination{
		Page:     page,
		PageSize: limit,
	}, nil
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

func BuildPagination(queryParams url.Values) (uint64, uint64, uint64, error) {
	var err error

	page := defaultPage
	pageSize := defaultPageSize

	pages, ok := queryParams["page"]
	if ok {
		page, err = strconv.Atoi(pages[0])
		if err != nil {
			return 0, 0, 0, errors.New("could not convert page to number")
		}
	}

	pageSizes, ok := queryParams["page_size"]
	if ok {
		pageSize, err = strconv.Atoi(pageSizes[0])
		if err != nil {
			return 0, 0, 0, errors.New("could not convert page_size to number")
		}
	}

	offset := pageSize * (page - 1)
	if pageSize < 0 || offset < 0 {
		return 0, 0, 0, errors.New("can not have negative page or page_size")
	}

	return uint64(page), uint64(pageSize), uint64(offset), nil
}
