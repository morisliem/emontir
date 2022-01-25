package model

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type (
	ServiceBaseModel struct {
		ID          int            `db:"id"`
		Title       string         `db:"title"`
		Description sql.NullString `db:"description"`
		Rating      float64        `db:"rating"`
		Price       float64        `db:"price"`
		Picture     sql.NullString `db:"picture"`
	}

	SortNFilter struct {
		Type   string
		Rating float64
		Sort   string
	}
)

type Service interface {
	GetAllServices(ctx context.Context, limit, offset int, condition SortNFilter) ([]ServiceBaseModel, error)
	// nolint
	SearchService(ctx context.Context, limit, offset int, keyboard string, condition SortNFilter) ([]ServiceBaseModel, error)
}

type service struct {
	db      *sqlx.DB
	queries map[string]*sqlx.Stmt
}

func NewService(db *sqlx.DB) Service {
	service := new(service)
	service.db = db
	service.queries = make(map[string]*sqlx.Stmt, len(serviceQueries))
	for k, v := range serviceQueries {
		stmt, err := db.Preparex(v)
		if err != nil {
			log.Fatal().Msg("error : " + err.Error() + "\nservice : " + v)
		}
		service.queries[k] = stmt
	}
	return service
}

var (
	getServices    = "GetAllServices"
	getServicesSQL = `SELECT "id", "title", "description", "rating", 
					"price", "picture" FROM "services" 
					ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	getServicesSortHighestRatingNoType    = "AllServicesSortHighestRatingWithoutType"
	getServicesSortHighestRatingNoTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3
										ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	getServicesSortHighestRatingWithType    = "AllServicesSortHighestRatingWithType"
	getServicesSortHighestRatingWithTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	getServicesSortHighestPriceNoType    = "AllServicesSortHighestPriceWithoutType"
	getServicesSortHighestPriceNoTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3
										ORDER BY "price" DESC LIMIT $1 OFFSET $2`

	getServicesSortHighestPriceWithType    = "AllServicesSortHighestPriceWithType"
	getServicesSortHighestPriceWithTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										ORDER BY "price" DESC LIMIT $1 OFFSET $2`

	getServicesSortLowestPriceNoType    = "AllServicesSortLowestPriceWithoutType"
	getServicesSortLowestPriceNoTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3
										ORDER BY "price" LIMIT $1 OFFSET $2`

	getServicesSortLowestPriceWithType    = "AllServicesSortLowestPriceWithType"
	getServicesSortLowestPriceWithTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										ORDER BY "price" LIMIT $1 OFFSET $2`

	getServicesSortTitleDescNoType    = "AllServicesSortTitleDescWithoutType"
	getServicesSortTitleDescNoTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3
										ORDER BY "title" DESC LIMIT $1 OFFSET $2`

	getServicesSortTitleDescWithType    = "AllServicesSortTitleDescWithType"
	getServicesSortTitleDescWithTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										ORDER BY "title" DESC LIMIT $1 OFFSET $2`

	getServicesSortTitleAscNoType    = "AllServicesSortTitleAscWithoutType"
	getServicesSortTitleAscNoTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3
										ORDER BY "title" LIMIT $1 OFFSET $2`

	getServicesSortTitleAscWithType    = "AllServicesSortTitleAscWithType"
	getServicesSortTitleAscWithTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										ORDER BY "title" LIMIT $1 OFFSET $2`

	// ================================================================== //

	searchService    = "SearchService"
	searchServiceSQL = `SELECT "id", "title", "description", "rating", "price", "picture" 
						FROM "services" WHERE "title" LIKE '%'||$1||'%' 
						ORDER BY "rating" DESC LIMIT $2 OFFSET $3`

	searchServicesSortHighestRatingNoType    = "SearchAllServicesSortHighestRatingWithoutType"
	searchServicesSortHighestRatingNoTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	searchServicesSortHighestRatingWithType    = "SearchAllServicesSortHighestRatingWithType"
	searchServicesSortHighestRatingWithTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%'
										AND "title" LIKE '%'||$5||'%'
										ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	searchServicesSortHighestPriceNoType    = "SearchAllServicesSortHighestPriceWithoutType"
	searchServicesSortHighestPriceNoTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										ORDER BY "price" DESC LIMIT $1 OFFSET $2`

	searchServicesSortHighestPriceWithType    = "SearchAllServicesSortHighestPriceWithType"
	searchServicesSortHighestPriceWithTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										AND "title" LIKE '%'||$5||'%'
										ORDER BY "price" DESC LIMIT $1 OFFSET $2`

	searchServicesSortLowestPriceNoType    = "SearchAllServicesSortLowestPriceWithoutType"
	searchServicesSortLowestPriceNoTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										ORDER BY "price" LIMIT $1 OFFSET $2`

	searchServicesSortLowestPriceWithType    = "SearchAllServicesSortLowestPriceWithType"
	searchServicesSortLowestPriceWithTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										AND "title" LIKE '%'||$5||'%'
										ORDER BY "price" LIMIT $1 OFFSET $2`

	searchServicesSortTitleDescNoType    = "SearchAllServicesSortTitleDescWithoutType"
	searchServicesSortTitleDescNoTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										ORDER BY "title" DESC LIMIT $1 OFFSET $2`

	searchServicesSortTitleDescWithType    = "SearchAllServicesSortTitleDescWithType"
	searchServicesSortTitleDescWithTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										AND "title" LIKE '%'||$5||'%'
										ORDER BY "title" DESC LIMIT $1 OFFSET $2`

	searchServicesSortTitleAscNoType    = "SearchAllServicesSortTitleAscWithoutType"
	searchServicesSortTitleAscNoTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										ORDER BY "title" LIMIT $1 OFFSET $2`

	searchServicesSortTitleAscWithType    = "SearchAllServicesSortTitleAscWithType"
	searchServicesSortTitleAscWithTypeSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										AND "title" LIKE '%'||$5||'%'
										ORDER BY "title" LIMIT $1 OFFSET $2`

	serviceQueries = map[string]string{
		getServices:                             getServicesSQL,
		getServicesSortHighestRatingNoType:      getServicesSortHighestRatingNoTypeSQL,
		getServicesSortHighestRatingWithType:    getServicesSortHighestRatingWithTypeSQL,
		getServicesSortHighestPriceWithType:     getServicesSortHighestPriceWithTypeSQL,
		getServicesSortHighestPriceNoType:       getServicesSortHighestPriceNoTypeSQL,
		getServicesSortLowestPriceWithType:      getServicesSortLowestPriceWithTypeSQL,
		getServicesSortLowestPriceNoType:        getServicesSortLowestPriceNoTypeSQL,
		getServicesSortTitleDescWithType:        getServicesSortTitleDescWithTypeSQL,
		getServicesSortTitleDescNoType:          getServicesSortTitleDescNoTypeSQL,
		getServicesSortTitleAscNoType:           getServicesSortTitleAscNoTypeSQL,
		getServicesSortTitleAscWithType:         getServicesSortTitleAscWithTypeSQL,
		searchService:                           searchServiceSQL,
		searchServicesSortHighestRatingNoType:   searchServicesSortHighestRatingNoTypeSQL,
		searchServicesSortHighestRatingWithType: searchServicesSortHighestRatingWithTypeSQL,
		searchServicesSortHighestPriceWithType:  searchServicesSortHighestPriceWithTypeSQL,
		searchServicesSortHighestPriceNoType:    searchServicesSortHighestPriceNoTypeSQL,
		searchServicesSortLowestPriceWithType:   searchServicesSortLowestPriceWithTypeSQL,
		searchServicesSortLowestPriceNoType:     searchServicesSortLowestPriceNoTypeSQL,
		searchServicesSortTitleDescWithType:     searchServicesSortTitleDescWithTypeSQL,
		searchServicesSortTitleDescNoType:       searchServicesSortTitleDescNoTypeSQL,
		searchServicesSortTitleAscNoType:        searchServicesSortTitleAscNoTypeSQL,
		searchServicesSortTitleAscWithType:      searchServicesSortTitleAscWithTypeSQL,
	}

	getServiceSortValNoType = map[string]string{
		"highest_rating": getServicesSortHighestRatingNoTypeSQL,
		"highest_price":  getServicesSortHighestPriceNoTypeSQL,
		"lowest_price":   getServicesSortLowestPriceNoTypeSQL,
		"name_a-z":       getServicesSortTitleAscNoTypeSQL,
		"name_z-a":       getServicesSortTitleDescNoTypeSQL,
	}

	getServiceSortValWithType = map[string]string{
		"highest_rating": getServicesSortHighestRatingWithTypeSQL,
		"highest_price":  getServicesSortHighestPriceWithTypeSQL,
		"lowest_price":   getServicesSortLowestPriceWithTypeSQL,
		"name_a-z":       getServicesSortTitleAscWithTypeSQL,
		"name_z-a":       getServicesSortTitleDescWithTypeSQL,
	}

	searchServiceSortValNoType = map[string]string{
		"highest_rating": searchServicesSortHighestRatingNoTypeSQL,
		"highest_price":  searchServicesSortHighestPriceNoTypeSQL,
		"lowest_price":   searchServicesSortLowestPriceNoTypeSQL,
		"name_a-z":       searchServicesSortTitleAscNoTypeSQL,
		"name_z-a":       searchServicesSortTitleDescNoTypeSQL,
	}

	searchServiceSortValWithType = map[string]string{
		"highest_rating": searchServicesSortHighestRatingWithTypeSQL,
		"highest_price":  searchServicesSortHighestPriceWithTypeSQL,
		"lowest_price":   searchServicesSortLowestPriceWithTypeSQL,
		"name_a-z":       searchServicesSortTitleAscWithTypeSQL,
		"name_z-a":       searchServicesSortTitleDescWithTypeSQL,
	}
)

// nolint
func (c *service) GetAllServices(ctx context.Context, limit, offset int, condition SortNFilter) ([]ServiceBaseModel, error) {
	var result []ServiceBaseModel
	if condition.Type == "" {
		// nolint
		if err := c.db.SelectContext(ctx, &result, getServiceSortValNoType[condition.Sort], limit, offset, condition.Rating); err != nil {
			return nil, err
		}
		return result, nil
	}
	// nolint
	if err := c.db.SelectContext(ctx, &result, getServiceSortValWithType[condition.Sort], limit, offset, condition.Rating, condition.Type); err != nil {
		return nil, err
	}
	return result, nil
}

// nolint
func (c *service) SearchService(ctx context.Context, limit, offset int, keyword string, condition SortNFilter) ([]ServiceBaseModel, error) {
	var result []ServiceBaseModel
	if condition.Type == "" {
		// nolint
		if err := c.db.SelectContext(ctx, &result, searchServiceSortValNoType[condition.Sort], limit, offset, condition.Rating, keyword); err != nil {
			return nil, err
		}
		return result, nil
	}
	// nolint
	if err := c.db.SelectContext(ctx, &result, searchServiceSortValWithType[condition.Sort], limit, offset, condition.Rating, keyword, condition.Type); err != nil {
		return nil, err
	}
	return result, nil
}
