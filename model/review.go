package model

import (
	"context"
	"database/sql"
	"e-montir/api/handler"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type (
	ReviewBaseModel struct {
		ID        int       `db:"id"`
		UserID    string    `db:"user_id"`
		Feedback  string    `db:"feedback"`
		Rating    float64   `db:"rating"`
		ServiceID int       `db:"service_id"`
		OrderID   string    `db:"order_id"`
		CreatedAt time.Time `db:"created_at"`
	}
)

type Review interface {
	AddServiceReview(ctx context.Context, param *ReviewBaseModel) error
	IsServiceReviewed(ctx context.Context, orderID string, serviceID int) (bool, error)
}

type review struct {
	db      *sqlx.DB
	queries map[string]*sqlx.Stmt
}

func NewReview(db *sqlx.DB) Review {
	review := new(review)
	review.db = db
	review.queries = make(map[string]*sqlx.Stmt, len(reviewQueries))
	for k, v := range reviewQueries {
		stmt, err := db.Preparex(v)
		if err != nil {
			log.Fatal().Msg("error : " + err.Error() + "\nreview : " + v)
		}
		review.queries[k] = stmt
	}
	return review
}

var (
	setServiceReview      = "serviceReview"
	setServiceReviewField = `("user_id", "feedback", "rating", "service_id", "order_id", "created_at")`
	setServiceReviewSQL   = `INSERT INTO "feedbacks" ` + setServiceReviewField + ` VALUES ($1,$2,$3,$4,$5,$6)`

	updateServiceRating    = "updateServiceRating"
	updateServiceRatingSQL = `UPDATE "services" SET "rating" = round((float8(rating + $2)/2)::numeric, 2) WHERE "id" = $1`

	getReviewByOrderIDNServiceID    = "getReviewByOrderIDNServiceID"
	getReviewByOrderIDNServiceIDSQL = `SELECT "user_id", "rating", "feedback" FROM "feedbacks" WHERE "service_id" = $1 AND "order_id" = $2`

	reviewQueries = map[string]string{
		setServiceReview:             setServiceReviewSQL,
		updateServiceRating:          updateServiceRatingSQL,
		getReviewByOrderIDNServiceID: getReviewByOrderIDNServiceIDSQL,
	}
)

func (c *review) AddServiceReview(ctx context.Context, param *ReviewBaseModel) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if rollback := tx.Rollback(); rollback == nil {
			log.Info().Msg("rolling back changes")
		}
	}()

	if param.Feedback == "" {
		_, err := tx.ExecContext(ctx, setServiceReviewSQL, param.UserID, nil, param.Rating, param.ServiceID, param.OrderID, param.CreatedAt)
		if err != nil {
			return err
		}
	} else {
		_, err := tx.ExecContext(ctx, setServiceReviewSQL, param.UserID, param.Feedback, param.Rating, param.ServiceID, param.OrderID, param.CreatedAt)
		if err != nil {
			return err
		}
	}

	row, err := tx.ExecContext(ctx, updateServiceRatingSQL, param.ServiceID, param.Rating)
	if err != nil {
		return err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil || rowAffected != 1 {
		return &handler.InternalServerError
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *review) IsServiceReviewed(ctx context.Context, orderID string, serviceID int) (bool, error) {
	var result ReviewBaseModel
	var feedback sql.NullString
	err := c.queries[getReviewByOrderIDNServiceID].QueryRowContext(ctx, serviceID, orderID).Scan(&result.UserID, &result.Rating, &feedback)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
