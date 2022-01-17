package controller

import (
	"context"
	"e-montir/api/handler"
	"e-montir/model"
	"e-montir/pkg/validator"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type reviewCtx struct {
	reviewModel model.Review
	cartModel   model.Cart
}

type Review interface {
	AddServiceReview(ctx context.Context, userID string, form *ReviewBaseModel) error
}

func NewReview(reviewModel model.Review, cartModel model.Cart) Review {
	return &reviewCtx{
		reviewModel: reviewModel,
		cartModel:   cartModel,
	}
}

type (
	ReviewBaseModel struct {
		ID              int       `json:"id"`
		Feedback        string    `json:"feedback"`
		RatingString    string    `json:"rating"`
		CreatedAt       time.Time `json:"created_at"`
		Rating          float64
		ServiceID       int
		ServiceIDString string
		OrderID         string
	}
)

func (req *ReviewBaseModel) ValidateReviewRequest() ([]handler.Fields, error) {
	var count int
	var fields []handler.Fields
	serviceIDValid := true
	ratingValid := true
	err := validator.ValidateRating(req.RatingString)
	if err != nil {
		ratingValid = false
		count++
		fields = append(fields, handler.Fields{
			Name:    "rating",
			Message: err.Error(),
		})
	}

	if ratingValid {
		rating, ratingErr := strconv.ParseFloat(req.RatingString, 64)
		if ratingErr != nil {
			count++
			ratingValid = false
			fields = append(fields, handler.Fields{
				Name:    "rating",
				Message: ratingErr.Error(),
			})
		}
		if rating < 0.1 && ratingValid || rating > 5 && ratingValid {
			count++
			fields = append(fields, handler.Fields{
				Name:    "rating",
				Message: "rating must be between 0.1 to 5",
			})
		}
		req.Rating = rating
	}

	err = validator.ValidateServiceID(req.ServiceIDString)
	if err != nil {
		serviceIDValid = false
		count++
		fields = append(fields, handler.Fields{
			Name:    "service_id",
			Message: err.Error(),
		})
	}

	if serviceIDValid {
		serviceID, serviceIDErr := strconv.Atoi(req.ServiceIDString)
		if serviceIDErr != nil {
			count++
			serviceIDValid = false
			fields = append(fields, handler.Fields{
				Name:    "service_id",
				Message: serviceIDErr.Error(),
			})
		}
		if serviceID < 1 && serviceIDValid {
			count++
			fields = append(fields, handler.Fields{
				Name:    "service_id",
				Message: "service_id must be more than 0",
			})
		}
		req.ServiceID = serviceID
	}

	err = validator.ValidateFeedback(req.Feedback)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "feedback",
			Message: err.Error(),
		})
	}

	err = validator.ValidateOrderID(req.OrderID)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "order_id",
			Message: err.Error(),
		})
	}

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

func (c *reviewCtx) AddServiceReview(ctx context.Context, userID string, form *ReviewBaseModel) error {
	isServiceAvailable, err := c.cartModel.IsServiceAvailable(ctx, form.ServiceID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when checking isServiceAvailable : %w", err)).Send()
		return err
	}

	if !isServiceAvailable {
		return &handler.ServiceNotExists
	}

	isReviewed, err := c.reviewModel.IsServiceReviewed(ctx, form.OrderID, form.ServiceID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when checking isServiceReviewed : %w", err)).Send()
		return err
	}

	if isReviewed {
		return &handler.ServiceIsReviewed
	}

	err = c.reviewModel.AddServiceReview(ctx, &model.ReviewBaseModel{
		UserID:    userID,
		Feedback:  form.Feedback,
		Rating:    form.Rating,
		ServiceID: form.ServiceID,
		OrderID:   form.OrderID,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Error().Err(fmt.Errorf("error when AddServiceReview : %w", err)).Send()
		return err
	}

	return nil
}
