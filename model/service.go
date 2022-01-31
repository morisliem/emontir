package model

import (
	"context"
	"database/sql"
	"time"

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

	FavServiceBaseModel struct {
		UserID    string    `db:"user_id"`
		ServiceID int       `db:"service_id"`
		CreatedAt time.Time `db:"created_at"`
	}
)

type Service interface {
	GetAllServices(ctx context.Context, limit, offset int, condition SortNFilter) ([]ServiceBaseModel, error)
	// nolint
	SearchService(ctx context.Context, limit, offset int, keyboard string, condition SortNFilter) ([]ServiceBaseModel, error)
	AddFavService(ctx context.Context, userID string, serviceID int) error
	RemoveFavService(ctx context.Context, userID string, serviceID int) error
	ListOfFavServices(ctx context.Context, userID string) ([]int, error)
	GetServiceByID(ctx context.Context, serviceID int) (*ServiceBaseModel, error)
	GetFavServiceByUserIDNServiceID(ctx context.Context, userID string, serviceID int) (int, error)
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

	getServiceByID    = "getServiceByID"
	getServiceByIDSQL = `SELECT "id", "title", "description", "rating", 
						"price", "picture" FROM "services" 
						WHERE "id" = $1`

	getServicesSortHighestRating    = "AllServicesSortHighestRating"
	getServicesSortHighestRatingSQL = `SELECT "id", "title", "description", "rating", 
									"price", "picture" FROM "services" 
									WHERE "rating" > $3 AND "category" LIKE '%'||$4||'%' 
									ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	getServicesSortHighestPrice    = "AllServicesSortHighestPrice"
	getServicesSortHighestPriceSQL = `SELECT "id", "title", "description", "rating", 
									"price", "picture" FROM "services" 
									WHERE "rating" > $3 AND "category" LIKE '%'||$4||'%' 
									ORDER BY "price" DESC LIMIT $1 OFFSET $2`

	getServicesSortLowestPrice    = "AllServicesSortLowestPrice"
	getServicesSortLowestPriceSQL = `SELECT "id", "title", "description", "rating", 
									"price", "picture" FROM "services" 
									WHERE "rating" > $3 AND "category" LIKE '%'||$4||'%' 
									ORDER BY "price" LIMIT $1 OFFSET $2`

	getServicesSortTitleDesc    = "AllServicesSortTitleDesc"
	getServicesSortTitleDescSQL = `SELECT "id", "title", "description", "rating", 
								"price", "picture" FROM "services" 
								WHERE "rating" > $3 AND "category" LIKE '%'||$4||'%' 
								ORDER BY "title" DESC LIMIT $1 OFFSET $2`

	getServicesSortTitleAsc    = "AllServicesSortTitleAsc"
	getServicesSortTitleAscSQL = `SELECT "id", "title", "description", "rating", 
								"price", "picture" FROM "services" 
								WHERE "rating" > $3 AND "category" LIKE '%'||$4||'%' 
								ORDER BY "title" LIMIT $1 OFFSET $2`

	getServicesSortHighestRatingTypePopular    = "AllServicesSortHighestRatingTypePopular"
	getServicesSortHighestRatingTypePopularSQL = `SELECT "id", "title", "description", "rating",
												"price", "picture" FROM "services"
												WHERE "rating" > $3 AND "number_of_order" >= 30
												ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	getServicesSortHighestPriceTypePopular    = "AllServicesSortHighestPriceTypePopular"
	getServicesSortHighestPriceTypePopularSQL = `SELECT "id", "title", "description", "rating",
												"price", "picture" FROM "services"
												WHERE "rating" > $3 AND "number_of_order" >= 30
												ORDER BY "price" DESC LIMIT $1 OFFSET $2`

	getServicesSortLowestPriceTypePopular    = "AllServicesSortLowestPriceTypePopular"
	getServicesSortLowestPriceTypePopularSQL = `SELECT "id", "title", "description", "rating",
												"price", "picture" FROM "services"
												WHERE "rating" > $3 AND "number_of_order" >= 30
												ORDER BY "price" LIMIT $1 OFFSET $2`

	getServicesSortTitleAscTypePopular    = "AllServicesSortTitleAscTypePopular"
	getServicesSortTitleAscTypePopularSQL = `SELECT "id", "title", "description", "rating",
											"price", "picture" FROM "services"
											WHERE "rating" > $3 AND "number_of_order" >= 30
											ORDER BY "title" LIMIT $1 OFFSET $2`

	getServicesSortTitleDescTypePopular    = "AllServicesSortTitleDescTypePopular"
	getServicesSortTitleDescTypePopularSQL = `SELECT "id", "title", "description", "rating",
											"price", "picture" FROM "services"
											WHERE "rating" > $3 AND "number_of_order" >= 30
											ORDER BY "title" DESC LIMIT $1 OFFSET $2`

	// ================================================================== //

	searchService    = "SearchService"
	searchServiceSQL = `SELECT "id", "title", "description", "rating", "price", "picture" 
						FROM "services" WHERE "title" LIKE '%'||$1||'%' 
						ORDER BY "rating" DESC LIMIT $2 OFFSET $3`

	searchServicesSortHighestRating    = "SearchAllServicesSortHighestRating"
	searchServicesSortHighestRatingSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%'
										AND "category" LIKE '%'||$5||'%'
										ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	searchServicesSortHighestPrice    = "SearchAllServicesSortHighestPrice"
	searchServicesSortHighestPriceSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										AND "category" LIKE '%'||$5||'%'
										ORDER BY "price" DESC LIMIT $1 OFFSET $2`

	searchServicesSortLowestPrice    = "SearchAllServicesSortLowestPrice"
	searchServicesSortLowestPriceSQL = `SELECT "id", "title", "description", "rating", 
										"price", "picture" FROM "services" 
										WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
										AND "category" LIKE '%'||$5||'%'
										ORDER BY "price" LIMIT $1 OFFSET $2`

	searchServicesSortTitleDesc    = "SearchAllServicesSortTitleDesc"
	searchServicesSortTitleDescSQL = `SELECT "id", "title", "description", "rating", 
									"price", "picture" FROM "services" 
									WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
									AND "category" LIKE '%'||$5||'%'
									ORDER BY "title" DESC LIMIT $1 OFFSET $2`

	searchServicesSortTitleAsc    = "SearchAllServicesSortTitleAsc"
	searchServicesSortTitleAscSQL = `SELECT "id", "title", "description", "rating", 
									"price", "picture" FROM "services" 
									WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%' 
									AND "category" LIKE '%'||$5||'%'
									ORDER BY "title" LIMIT $1 OFFSET $2`

	searchServicesSortHighestRatingTypePopular    = "SearchAllServicesSortHighestRatingTypePopular"
	searchServicesSortHighestRatingTypePopularSQL = `SELECT "id", "title", "description", "rating"
													"price", "picture" FROM "services"
													WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%'
													AND "number_of_order" >= 30
													ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	searchServicesSortHighestPriceTypePopular    = "SearchAllServicesSortHighestPriceTypePopular"
	searchServicesSortHighestPriceTypePopularSQL = `SELECT "id", "title", "description", "rating"
													"price", "picture" FROM "services"
													WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%'
													AND "number_of_order" >= 30
													ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	searchServicesSortLowestPriceTypePopular    = "SearchAllServicesSortLowestPriceTypePopular"
	searchServicesSortLowestPriceTypePopularSQL = `SELECT "id", "title", "description", "rating"
													"price", "picture" FROM "services"
													WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%'
													AND "number_of_order" >= 30
													ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	searchServicesSortTitleAscTypePopular    = "SearchAllServicesSortTitleAscTypePopular"
	searchServicesSortTitleAscTypePopularSQL = `SELECT "id", "title", "description", "rating"
													"price", "picture" FROM "services"
													WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%'
													AND "number_of_order" >= 30
													ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	searchServicesSortTitleDescTypePopular    = "SearchAllServicesSortTitleDescTypePopular"
	searchServicesSortTitleDescTypePopularSQL = `SELECT "id", "title", "description", "rating"
													"price", "picture" FROM "services"
													WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%'
													AND "number_of_order" >= 30
													ORDER BY "rating" DESC LIMIT $1 OFFSET $2`

	// ================================================================== //

	addFavService    = "addFavoriteService"
	addFavServiceSQL = `INSERT INTO "user_fav_services" ("user_id","service_id","created_at")
						VALUES ($1,$2,$3)`

	removeFavService    = "removeFavoriteService"
	removeFavServiceSQL = `DELETE FROM "user_fav_services" WHERE "user_id" = $1 AND "service_id" = $2`

	getUserFavServices    = "getUserFavoriteServices"
	getUserFavServicesSQL = `SELECT "service_id" FROM "user_fav_services" 
							WHERE "user_id" = $1 ORDER BY "created_at" DESC`

	getFavServiceByUserIDNServiceID    = "getFavServiceByUserIDNServiceID"
	getFavServiceByUserIDNServiceIDSQL = `SELECT "service_id" FROM "user_fav_services" 
										WHERE "user_id" = $1 AND "service_id" = $2`

	serviceQueries = map[string]string{
		getServices:                                getServicesSQL,
		getServicesSortHighestRating:               getServicesSortHighestRatingSQL,
		getServicesSortHighestPrice:                getServicesSortHighestPriceSQL,
		getServicesSortLowestPrice:                 getServicesSortLowestPriceSQL,
		getServicesSortTitleAsc:                    getServicesSortTitleAscSQL,
		getServicesSortTitleDesc:                   getServicesSortTitleDescSQL,
		getServicesSortHighestRatingTypePopular:    getServicesSortHighestRatingTypePopularSQL,
		getServicesSortHighestPriceTypePopular:     getServicesSortHighestPriceTypePopularSQL,
		getServicesSortLowestPriceTypePopular:      getServicesSortLowestPriceTypePopularSQL,
		getServicesSortTitleAscTypePopular:         getServicesSortTitleAscTypePopularSQL,
		getServicesSortTitleDescTypePopular:        getServicesSortTitleDescTypePopularSQL,
		searchService:                              searchServiceSQL,
		searchServicesSortHighestRating:            searchServicesSortHighestRatingSQL,
		searchServicesSortHighestPrice:             searchServicesSortHighestPriceSQL,
		searchServicesSortLowestPrice:              searchServicesSortLowestPriceSQL,
		searchServicesSortTitleAsc:                 getServicesSortTitleAscSQL,
		searchServicesSortTitleDesc:                getServicesSortTitleDescSQL,
		searchServicesSortHighestRatingTypePopular: searchServicesSortHighestRatingTypePopularSQL,
		searchServicesSortHighestPriceTypePopular:  searchServicesSortHighestPriceTypePopularSQL,
		searchServicesSortLowestPriceTypePopular:   searchServicesSortLowestPriceTypePopularSQL,
		searchServicesSortTitleAscTypePopular:      searchServicesSortTitleAscTypePopularSQL,
		searchServicesSortTitleDescTypePopular:     searchServicesSortTitleDescTypePopularSQL,
		addFavService:                              addFavServiceSQL,
		removeFavService:                           removeFavServiceSQL,
		getUserFavServices:                         getUserFavServicesSQL,
		getServiceByID:                             getServiceByIDSQL,
		getFavServiceByUserIDNServiceID:            getFavServiceByUserIDNServiceIDSQL,
	}

	getServiceSort = map[string]string{
		"highest_rating": getServicesSortHighestRatingSQL,
		"highest_price":  getServicesSortHighestPriceSQL,
		"lowest_price":   getServicesSortLowestPriceSQL,
		"name_a-z":       getServicesSortTitleAscSQL,
		"name_z-a":       getServicesSortTitleDescSQL,
	}

	getServiceSortValTypePopular = map[string]string{
		"highest_rating": getServicesSortHighestRatingTypePopularSQL,
		"highest_price":  getServicesSortHighestPriceTypePopularSQL,
		"lowest_price":   getServicesSortLowestPriceTypePopularSQL,
		"name_a-z":       getServicesSortTitleAscTypePopularSQL,
		"name_z-a":       getServicesSortTitleDescTypePopularSQL,
	}

	searchServiceSort = map[string]string{
		"highest_rating": searchServicesSortHighestRatingSQL,
		"highest_price":  searchServicesSortHighestPriceSQL,
		"lowest_price":   searchServicesSortLowestPriceSQL,
		"name_a-z":       searchServicesSortTitleAscSQL,
		"name_z-a":       searchServicesSortTitleDescSQL,
	}

	searchServiceSortValTypePopular = map[string]string{
		"highest_rating": searchServicesSortHighestRatingTypePopularSQL,
		"highest_price":  searchServicesSortHighestPriceTypePopularSQL,
		"lowest_price":   searchServicesSortLowestPriceTypePopularSQL,
		"name_a-z":       searchServicesSortTitleAscTypePopularSQL,
		"name_z-a":       searchServicesSortTitleDescTypePopularSQL,
	}
)

// nolint
func (c *service) GetAllServices(ctx context.Context, limit, offset int, condition SortNFilter) ([]ServiceBaseModel, error) {
	var result []ServiceBaseModel
	if condition.Type == "all" {
		// nolint
		if err := c.db.SelectContext(ctx, &result, getServiceSort[condition.Sort], limit, offset, condition.Rating, ""); err != nil {
			return nil, err
		}
		return result, nil
	}
	if condition.Type == "popular" {
		// nolint
		if err := c.db.SelectContext(ctx, &result, getServiceSortValTypePopular[condition.Sort], limit, offset, condition.Rating); err != nil {
			return nil, err
		}
		return result, nil
	}
	// nolint
	if err := c.db.SelectContext(ctx, &result, getServiceSort[condition.Sort], limit, offset, condition.Rating, condition.Type); err != nil {
		return nil, err
	}
	return result, nil
}

// nolint
func (c *service) SearchService(ctx context.Context, limit, offset int, keyword string, condition SortNFilter) ([]ServiceBaseModel, error) {
	var result []ServiceBaseModel
	if condition.Type == "all" {
		// nolint
		if err := c.db.SelectContext(ctx, &result, searchServiceSort[condition.Sort], limit, offset, condition.Rating, keyword, ""); err != nil {
			return nil, err
		}
		return result, nil
	}
	if condition.Type == "popular" {
		// nolint
		if err := c.db.SelectContext(ctx, &result, searchServiceSortValTypePopular[condition.Sort], limit, offset, condition.Rating, keyword); err != nil {
			return nil, err
		}
		return result, nil
	}
	// nolint
	if err := c.db.SelectContext(ctx, &result, searchServiceSort[condition.Sort], limit, offset, condition.Rating, keyword, condition.Type); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *service) AddFavService(ctx context.Context, userID string, serviceID int) error {
	_, err := c.queries[addFavService].ExecContext(ctx, userID, serviceID, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (c *service) RemoveFavService(ctx context.Context, userID string, serviceID int) error {
	_, err := c.queries[removeFavService].ExecContext(ctx, userID, serviceID)
	if err != nil {
		return err
	}
	return nil
}

func (c *service) ListOfFavServices(ctx context.Context, userID string) ([]int, error) {
	var result []int
	err := c.queries[getUserFavServices].SelectContext(ctx, &result, userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *service) GetServiceByID(ctx context.Context, serviceID int) (*ServiceBaseModel, error) {
	var result ServiceBaseModel
	err := c.queries[getServiceByID].GetContext(ctx, &result, serviceID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *service) GetFavServiceByUserIDNServiceID(ctx context.Context, userID string, serviceID int) (int, error) {
	var service sql.NullInt64
	err := c.queries[getFavServiceByUserIDNServiceID].QueryRowContext(ctx, userID, serviceID).Scan(&service)
	if err != nil {
		return -1, err
	}
	return int(service.Int64), nil
}
