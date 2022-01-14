package controller

import (
	"context"
	"e-montir/api/handler"
	"e-montir/model"
	"e-montir/pkg/validator"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

type timeslotCtx struct {
	timeslotModel model.Timeslot
}

type Timeslot interface {
	GetTimeslot(ctx context.Context, date string) (ListOfTimeslotResponse, error)
}

func NewTimeslot(timeslotModel model.Timeslot) Timeslot {
	return &timeslotCtx{
		timeslotModel: timeslotModel,
	}
}

type (
	TimeslotItem struct {
		Time        string `json:"time"`
		EmployeeNum int    `json:"employee_num"`
	}

	TimeslotList struct {
		Date string         `json:"date"`
		Data []TimeslotItem `json:"data"`
	}

	ListOfTimeslotResponse struct {
		Timeslot TimeslotList `json:"time_slot"`
	}

	TimeslotRequest struct {
		Date string // yyyy-mm-dd
	}
)

func (req *TimeslotRequest) ValidateTimeslotRequest() ([]handler.Fields, error) {
	var count int
	var fields []handler.Fields
	date, err := validator.ValidateDate(req.Date)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "date",
			Message: err.Error(),
		})
	}

	req.Date = date
	if count != 0 {
		return fields, errors.New("validation-failed")
	}
	return nil, nil
}

func (c *timeslotCtx) GetTimeslot(ctx context.Context, date string) (ListOfTimeslotResponse, error) {
	var timeslotList TimeslotList
	timeslotItem := make([]TimeslotItem, 0)
	res, err := c.timeslotModel.GetTimeslot(ctx, date)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when GetTimeslot : %w", err)).Send()
		return ListOfTimeslotResponse{}, &handler.InternalServerError
	}

	for _, v := range res {
		tmp := TimeslotItem{
			Time:        v.Time,
			EmployeeNum: v.EmployeeNum,
		}
		timeslotItem = append(timeslotItem, tmp)
	}

	timeslotList.Date = date
	timeslotList.Data = timeslotItem
	return ListOfTimeslotResponse{
		Timeslot: timeslotList,
	}, nil
}
