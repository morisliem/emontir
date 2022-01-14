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
)

type Service interface {
	GetAllServices(ctx context.Context, limit, offset int) ([]ServiceBaseModel, error)
	SearchService(ctx context.Context, limit, offset int, keyboard string) ([]ServiceBaseModel, error)
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
	serviceGetAllServices       = "GetAllService"
	dataToRetrieveForAllService = `"id", "title", "description", "rating", "price", "picture"`
	getAllServiceOrderBy        = `ORDER BY "rating" DESC LIMIT $1 OFFSET $2`
	serviceGetAllServicesSQL    = `SELECT` + dataToRetrieveForAllService + `FROM "services"` + getAllServiceOrderBy

	searchService                  = "SearchService"
	dataToRetrieveForSearchService = `SELECT "id", "title", "description", "rating", "price", "picture"`
	searchServiceWhereStmt         = `FROM "services" WHERE "title" LIKE '%'||$1||'%' `
	searchServiceOrderBy           = `ORDER BY "rating" DESC LIMIT $2 OFFSET $3`
	searchServiceSQL               = dataToRetrieveForSearchService + searchServiceWhereStmt + searchServiceOrderBy
	serviceQueries                 = map[string]string{
		serviceGetAllServices: serviceGetAllServicesSQL,
		searchService:         searchServiceSQL,
	}
)

func (c *service) GetAllServices(ctx context.Context, limit, offset int) ([]ServiceBaseModel, error) {
	var result []ServiceBaseModel
	if err := c.queries[serviceGetAllServices].SelectContext(ctx, &result, limit, offset); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *service) SearchService(ctx context.Context, limit, offset int, keyword string) ([]ServiceBaseModel, error) {
	var result []ServiceBaseModel
	if err := c.queries[searchService].SelectContext(ctx, &result, keyword, limit, offset); err != nil {
		return nil, err
	}
	return result, nil
}
