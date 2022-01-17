package v1

import (
	"e-montir/api/handler"
	"e-montir/controller"
	"net/http"

	"github.com/go-chi/chi"
)

type ReviewHandler struct {
	reviewController controller.Review
}

func NewReviewHandler(reviewController controller.Review) ReviewHandler {
	return ReviewHandler{
		reviewController: reviewController,
	}
}

func (c *ReviewHandler) AddServiceReview(w http.ResponseWriter, r *http.Request) {
	request := new(controller.ReviewBaseModel)
	request.ServiceIDString = chi.URLParam(r, "service_id")
	request.OrderID = chi.URLParam(r, "order_id")
	userID := handler.GetTokenClaim(r.Context()).ID

	if err := handler.DecodeJSON(r, request); err != nil {
		handler.ResponseError(w, &handler.ParsePayloadError)
		return
	}

	fieldsErr, err := request.ValidateReviewRequest()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	err = c.reviewController.AddServiceReview(r.Context(), userID, request)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}

	handler.GenerateResponse(w, http.StatusOK, handler.DefaultSuccess{Success: true})
}
