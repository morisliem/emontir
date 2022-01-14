package model

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type (
	TimeslotBaseModel struct {
		ID          string    `db:"id"`
		Time        string    `db:"time"`
		EmployeeNum int       `db:"employee_num"`
		Date        time.Time `db:"date"` // yyyy-mm-dd
	}
)

type Timeslot interface {
	GetTimeslot(ctx context.Context, date string) ([]TimeslotBaseModel, error)
}

type timeslot struct {
	db      *sqlx.DB
	queries map[string]*sqlx.Stmt
}

func NewTimeslot(db *sqlx.DB) Timeslot {
	timeslot := new(timeslot)
	timeslot.db = db
	timeslot.queries = make(map[string]*sqlx.Stmt, len(timeslotQueries))
	for k, v := range timeslotQueries {
		stmt, err := db.Preparex(v)
		if err != nil {
			log.Fatal().Msg("error : " + err.Error() + "\ntimeslot : " + v)
		}
		timeslot.queries[k] = stmt
	}
	return timeslot
}

var (
	getTimeslot     = "GetTimeslot"
	getTimeslotSQL  = `SELECT "time","employee_num" FROM "time_slots" WHERE "date" = $1 ORDER BY "time"`
	timeslotQueries = map[string]string{
		getTimeslot: getTimeslotSQL,
	}
)

func (c *timeslot) GetTimeslot(ctx context.Context, date string) ([]TimeslotBaseModel, error) {
	var result []TimeslotBaseModel
	if err := c.queries[getTimeslot].SelectContext(ctx, &result, date); err != nil {
		return nil, err
	}
	return result, nil
}
