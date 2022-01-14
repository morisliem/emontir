package v1

import (
	"e-montir/api/handler"
	"e-montir/controller"
	"net/http"
)

type TimeslotHandler struct {
	timeslotController controller.Timeslot
}

func NewTimeslotHandler(timeslotController controller.Timeslot) TimeslotHandler {
	return TimeslotHandler{
		timeslotController: timeslotController,
	}
}

func (c *TimeslotHandler) ListOfTimeslot(w http.ResponseWriter, r *http.Request) {
	request := new(controller.TimeslotRequest)
	request.Date = r.URL.Query().Get("date")

	fieldsErr, err := request.ValidateTimeslotRequest()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	res, err := c.timeslotController.GetTimeslot(r.Context(), request.Date)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, res)
}
