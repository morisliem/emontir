package model

import (
	"context"
	"database/sql"
	"e-montir/pkg/filter"
	"e-montir/pkg/sort"
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

	ListCriteria struct {
		Type    string
		Rating  float64
		Sort    string
		Limit   int
		Offset  int
		Keyword string
	}

	FavServiceBaseModel struct {
		UserID    string    `db:"user_id"`
		ServiceID int       `db:"service_id"`
		CreatedAt time.Time `db:"created_at"`
	}
)

type Service interface {
	GetAllServices(ctx context.Context, userID string, condition ListCriteria) ([]ServiceBaseModel, error)
	SearchService(ctx context.Context, condition ListCriteria) ([]ServiceBaseModel, error)
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

const (
	highestRating = ` ORDER BY "rating" DESC LIMIT $1 OFFSET $2`
	highestPrice  = ` ORDER BY "price" DESC LIMIT $1 OFFSET $2`
	lowestPrice   = ` ORDER BY "price" ASC LIMIT $1 OFFSET $2`
	titleAsc      = ` ORDER BY "title" ASC LIMIT $1 OFFSET $2`
	titleDesc     = ` ORDER BY "title" DESC LIMIT $1 OFFSET $2`
)

var (
	baseGetServiceSQLWithType = `SELECT services.id, services.title, services.description, services.rating, 
						services.price, services.picture FROM "service_categories"
						LEFT OUTER JOIN "services" ON services.id = service_categories.service_id
						WHERE services.rating > $4 AND service_categories.category = $5
						AND "id" NOT IN `

	baseGetPopularServiceSQL = `SELECT "id", "title", "description", "rating",
								"price", "picture" FROM "services"
								WHERE "rating" > $4 AND "number_of_order" >= 30
								AND "id" NOT IN `

	baseGetServicesWithoutTypeSQL = `SELECT "id", "title", "description", "rating", 
									"price", "picture" FROM "services"
									WHERE "rating" > $4 AND "id" NOT IN `

	getUserFavServiceSQL = `(SELECT "service_id" FROM "user_fav_services" 
							WHERE "user_id" = $3)`

	getServicesSortHighestRating    = "AllServicesSortHighestRating"
	getServicesSortHighestRatingSQL = baseGetServiceSQLWithType + getUserFavServiceSQL + highestRating

	getServicesSortHighestPrice    = "AllServicesSortHighestPrice"
	getServicesSortHighestPriceSQL = baseGetServiceSQLWithType + getUserFavServiceSQL + highestPrice

	getServicesSortLowestPrice    = "AllServicesSortLowestPrice"
	getServicesSortLowestPriceSQL = baseGetServiceSQLWithType + getUserFavServiceSQL + lowestPrice

	getServicesSortTitleDesc    = "AllServicesSortTitleDesc"
	getServicesSortTitleDescSQL = baseGetServiceSQLWithType + getUserFavServiceSQL + titleDesc

	getServicesSortTitleAsc    = "AllServicesSortTitleAsc"
	getServicesSortTitleAscSQL = baseGetServiceSQLWithType + getUserFavServiceSQL + titleAsc

	getServicesSortHighestRatingTypePopular    = "AllServicesSortHighestRatingTypePopular"
	getServicesSortHighestRatingTypePopularSQL = baseGetPopularServiceSQL + getUserFavServiceSQL + highestRating

	getServicesSortHighestPriceTypePopular    = "AllServicesSortHighestPriceTypePopular"
	getServicesSortHighestPriceTypePopularSQL = baseGetPopularServiceSQL + getUserFavServiceSQL + highestPrice

	getServicesSortLowestPriceTypePopular    = "AllServicesSortLowestPriceTypePopular"
	getServicesSortLowestPriceTypePopularSQL = baseGetPopularServiceSQL + getUserFavServiceSQL + lowestPrice

	getServicesSortTitleAscTypePopular    = "AllServicesSortTitleAscTypePopular"
	getServicesSortTitleAscTypePopularSQL = baseGetPopularServiceSQL + getUserFavServiceSQL + titleAsc

	getServicesSortTitleDescTypePopular    = "AllServicesSortTitleDescTypePopular"
	getServicesSortTitleDescTypePopularSQL = baseGetPopularServiceSQL + getUserFavServiceSQL + titleDesc

	getServicesSortHighestRatingWithoutType    = "getServiceSortHighestRatingWithoutType"
	getServicesSortHighestRatingWithoutTypeSQL = baseGetServicesWithoutTypeSQL + getUserFavServiceSQL + highestRating

	getServicesSortHighestPriceWithoutType    = "AllServicesSortHighestPriceWithoutType"
	getServicesSortHighestPriceWithoutTypeSQL = baseGetServicesWithoutTypeSQL + getUserFavServiceSQL + highestPrice

	getServicesSortLowestPriceWithoutType    = "AllServicesSortLowestPriceWithoutType"
	getServicesSortLowestPriceWithoutTypeSQL = baseGetServicesWithoutTypeSQL + getUserFavServiceSQL + lowestPrice

	getServicesSortTitleDescWithoutType    = "AllServicesSortTitleDescWithoutType"
	getServicesSortTitleDescWithoutTypeSQL = baseGetServicesWithoutTypeSQL + getUserFavServiceSQL + titleDesc

	getServicesSortTitleAscWithoutType    = "AllServicesSortTitleAscWithoutType"
	getServicesSortTitleAscWithoutTypeSQL = baseGetServicesWithoutTypeSQL + getUserFavServiceSQL + titleAsc

	baseSearchServiceWithTypeSQL = `SELECT services.id, services.title, services.description, 
							services.rating, services.price, services.picture 
							FROM "service_categories" LEFT OUTER 
							JOIN "services" ON services.id = "service_categories".service_id
							WHERE services.rating > $3  AND services.title LIKE '%'||$4||'%' 
							AND service_categories.category = $5`

	baseSearchPopularServiceSQL = `SELECT "id", "title", "description", "rating", 
									"price", "picture" FROM "services"
									WHERE "rating" > $3 AND "title" LIKE '%'||$4||'%'
									AND "number_of_order" >= 30`

	searchServiceWithoutTypeSQL = `SELECT "id", "title", "description", "rating", "price", "picture" 
									FROM "services" WHERE "title" LIKE '%'||$4||'%' 
									AND "rating" > $3`

	searchServicesSortHighestRating    = "SearchAllServicesSortHighestRating"
	searchServicesSortHighestRatingSQL = baseSearchServiceWithTypeSQL + highestRating

	searchServicesSortHighestPrice    = "SearchAllServicesSortHighestPrice"
	searchServicesSortHighestPriceSQL = baseSearchServiceWithTypeSQL + highestPrice

	searchServicesSortLowestPrice    = "SearchAllServicesSortLowestPrice"
	searchServicesSortLowestPriceSQL = baseSearchServiceWithTypeSQL + lowestPrice

	searchServicesSortTitleDesc    = "SearchAllServicesSortTitleDesc"
	searchServicesSortTitleDescSQL = baseSearchServiceWithTypeSQL + titleDesc

	searchServicesSortTitleAsc    = "SearchAllServicesSortTitleAsc"
	searchServicesSortTitleAscSQL = baseSearchServiceWithTypeSQL + titleAsc

	searchServicesSortHighestRatingTypePopular    = "SearchAllServicesSortHighestRatingTypePopular"
	searchServicesSortHighestRatingTypePopularSQL = baseSearchPopularServiceSQL + highestRating

	searchServicesSortHighestPriceTypePopular    = "SearchAllServicesSortHighestPriceTypePopular"
	searchServicesSortHighestPriceTypePopularSQL = baseSearchPopularServiceSQL + highestPrice

	searchServicesSortLowestPriceTypePopular    = "SearchAllServicesSortLowestPriceTypePopular"
	searchServicesSortLowestPriceTypePopularSQL = baseSearchPopularServiceSQL + lowestPrice

	searchServicesSortTitleAscTypePopular    = "SearchAllServicesSortTitleAscTypePopular"
	searchServicesSortTitleAscTypePopularSQL = baseSearchPopularServiceSQL + titleAsc

	searchServicesSortTitleDescTypePopular    = "SearchAllServicesSortTitleDescTypePopular"
	searchServicesSortTitleDescTypePopularSQL = baseSearchPopularServiceSQL + titleDesc

	searchServicesSortHighestRatingWithoutType    = "SearchAllServicesSortHighestRatingWithoutType"
	searchServicesSortHighestRatingWithoutTypeSQL = searchServiceWithoutTypeSQL + highestRating

	searchServicesSortHighestPriceWithoutType    = "SearchAllServicesSortHighestPriceWithoutType"
	searchServicesSortHighestPriceWithoutTypeSQL = searchServiceWithoutTypeSQL + highestPrice

	searchServicesSortLowestPriceWithoutType    = "SearchAllServicesSortLowestPriceWithoutType"
	searchServicesSortLowestPriceWithoutTypeSQL = searchServiceWithoutTypeSQL + lowestPrice

	searchServicesSortTitleDescWithoutType    = "SearchAllServicesSortTitleDescWithoutType"
	searchServicesSortTitleDescWithoutTypeSQL = searchServiceWithoutTypeSQL + titleDesc

	searchServicesSortTitleAscWithoutType    = "SearchAllServicesSortTitleAscWithoutType"
	searchServicesSortTitleAscWithoutTypeSQL = searchServiceWithoutTypeSQL + titleAsc

	// ================================================================== //

	getServiceByID    = "getServiceByID"
	getServiceByIDSQL = `SELECT "id", "title", "description", "rating", 
						"price", "picture" FROM "services" 
						WHERE "id" = $1`

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
		getServicesSortHighestRatingWithoutType:    getServicesSortHighestRatingWithoutTypeSQL,
		getServicesSortHighestPriceWithoutType:     getServicesSortHighestPriceWithoutTypeSQL,
		getServicesSortLowestPriceWithoutType:      getServicesSortLowestPriceWithoutTypeSQL,
		getServicesSortTitleAscWithoutType:         getServicesSortTitleAscWithoutTypeSQL,
		getServicesSortTitleDescWithoutType:        getServicesSortTitleDescWithoutTypeSQL,
		searchServicesSortHighestRating:            searchServicesSortHighestRatingSQL,
		searchServicesSortHighestPrice:             searchServicesSortHighestPriceSQL,
		searchServicesSortLowestPrice:              searchServicesSortLowestPriceSQL,
		searchServicesSortTitleAsc:                 searchServicesSortTitleAscSQL,
		searchServicesSortTitleDesc:                searchServicesSortTitleDescSQL,
		searchServicesSortHighestRatingTypePopular: searchServicesSortHighestRatingTypePopularSQL,
		searchServicesSortHighestPriceTypePopular:  searchServicesSortHighestPriceTypePopularSQL,
		searchServicesSortLowestPriceTypePopular:   searchServicesSortLowestPriceTypePopularSQL,
		searchServicesSortTitleAscTypePopular:      searchServicesSortTitleAscTypePopularSQL,
		searchServicesSortTitleDescTypePopular:     searchServicesSortTitleDescTypePopularSQL,
		searchServicesSortHighestRatingWithoutType: searchServicesSortHighestRatingWithoutTypeSQL,
		searchServicesSortHighestPriceWithoutType:  searchServicesSortHighestPriceWithoutTypeSQL,
		searchServicesSortLowestPriceWithoutType:   searchServicesSortLowestPriceWithoutTypeSQL,
		searchServicesSortTitleAscWithoutType:      searchServicesSortTitleAscWithoutTypeSQL,
		searchServicesSortTitleDescWithoutType:     searchServicesSortTitleDescWithoutTypeSQL,
		addFavService:                              addFavServiceSQL,
		removeFavService:                           removeFavServiceSQL,
		getUserFavServices:                         getUserFavServicesSQL,
		getServiceByID:                             getServiceByIDSQL,
		getFavServiceByUserIDNServiceID:            getFavServiceByUserIDNServiceIDSQL,
	}

	getServiceSort = map[string]string{
		sort.HighestRating: getServicesSortHighestRatingSQL,
		sort.HighestPrice:  getServicesSortHighestPriceSQL,
		sort.LowestPrice:   getServicesSortLowestPriceSQL,
		sort.NameAsc:       getServicesSortTitleAscSQL,
		sort.NameDesc:      getServicesSortTitleDescSQL,
	}

	getServiceSortWithoutType = map[string]string{
		sort.HighestRating: getServicesSortHighestRatingWithoutTypeSQL,
		sort.HighestPrice:  getServicesSortHighestPriceWithoutTypeSQL,
		sort.LowestPrice:   getServicesSortLowestPriceWithoutTypeSQL,
		sort.NameAsc:       getServicesSortTitleAscWithoutTypeSQL,
		sort.NameDesc:      getServicesSortTitleDescWithoutTypeSQL,
	}

	getServiceSortValTypePopular = map[string]string{
		sort.HighestRating: getServicesSortHighestRatingTypePopularSQL,
		sort.HighestPrice:  getServicesSortHighestPriceTypePopularSQL,
		sort.LowestPrice:   getServicesSortLowestPriceTypePopularSQL,
		sort.NameAsc:       getServicesSortTitleAscTypePopularSQL,
		sort.NameDesc:      getServicesSortTitleDescTypePopularSQL,
	}

	searchServiceSort = map[string]string{
		sort.HighestRating: searchServicesSortHighestRatingSQL,
		sort.HighestPrice:  searchServicesSortHighestPriceSQL,
		sort.LowestPrice:   searchServicesSortLowestPriceSQL,
		sort.NameAsc:       searchServicesSortTitleAscSQL,
		sort.NameDesc:      searchServicesSortTitleDescSQL,
	}

	searchServiceSortWithoutType = map[string]string{
		sort.HighestRating: searchServicesSortHighestRatingWithoutTypeSQL,
		sort.HighestPrice:  searchServicesSortHighestPriceWithoutTypeSQL,
		sort.LowestPrice:   searchServicesSortLowestPriceWithoutTypeSQL,
		sort.NameAsc:       searchServicesSortTitleAscWithoutTypeSQL,
		sort.NameDesc:      searchServicesSortTitleDescWithoutTypeSQL,
	}

	searchServiceSortValTypePopular = map[string]string{
		sort.HighestRating: searchServicesSortHighestRatingTypePopularSQL,
		sort.HighestPrice:  searchServicesSortHighestPriceTypePopularSQL,
		sort.LowestPrice:   searchServicesSortLowestPriceTypePopularSQL,
		sort.NameAsc:       searchServicesSortTitleAscTypePopularSQL,
		sort.NameDesc:      searchServicesSortTitleDescTypePopularSQL,
	}
)

func (c *service) GetAllServices(ctx context.Context, userID string, condition ListCriteria) ([]ServiceBaseModel, error) {
	var result []ServiceBaseModel
	if condition.Type == filter.All {
		if err := c.db.SelectContext(ctx, &result,
			getServiceSortWithoutType[condition.Sort],
			condition.Limit, condition.Offset, userID,
			condition.Rating); err != nil {
			return nil, err
		}
		return result, nil
	}
	if condition.Type == filter.Popular {
		if err := c.db.SelectContext(ctx, &result,
			getServiceSortValTypePopular[condition.Sort],
			condition.Limit, condition.Offset, userID,
			condition.Rating); err != nil {
			return nil, err
		}
		return result, nil
	}
	if err := c.db.SelectContext(ctx, &result,
		getServiceSort[condition.Sort], condition.Limit,
		condition.Offset, condition.Rating, userID,
		condition.Type); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *service) SearchService(ctx context.Context, condition ListCriteria) ([]ServiceBaseModel, error) {
	var result []ServiceBaseModel
	if condition.Type == filter.All {
		if err := c.db.SelectContext(ctx, &result,
			searchServiceSortWithoutType[condition.Sort],
			condition.Limit, condition.Offset, condition.Rating,
			condition.Keyword); err != nil {
			return nil, err
		}
		return result, nil
	}
	if condition.Type == filter.Popular {
		if err := c.db.SelectContext(ctx, &result,
			searchServiceSortValTypePopular[condition.Sort],
			condition.Limit, condition.Offset,
			condition.Rating, condition.Keyword); err != nil {
			return nil, err
		}
		return result, nil
	}
	if err := c.db.SelectContext(ctx, &result,
		searchServiceSort[condition.Sort], condition.Limit,
		condition.Offset, condition.Rating,
		condition.Keyword, condition.Type); err != nil {
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
