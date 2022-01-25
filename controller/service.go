package controller

import (
	"context"
	"e-montir/api/handler"
	"e-montir/model"
	"e-montir/pkg/validator"
	"errors"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
)

type serviceCtx struct {
	serviceModel model.Service
}

type Service interface {
	GetAllServices(ctx context.Context, condition *ServiceListRequest) (ServiceListsResponse, error)
	SearchService(ctx context.Context, condition *SearchServiceRequest) (ServiceListsResponse, error)
}

func NewService(serviceModel model.Service) Service {
	return &serviceCtx{
		serviceModel: serviceModel,
	}
}

type (
	ServiceListRequest struct {
		Page         int
		Limit        int
		PageString   string
		LimitString  string
		Type         string
		Rating       float64
		RatingString string
		Sort         string
	}

	ServiceItem struct {
		ID          int     `json:"id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Rating      float64 `json:"rating"`
		Price       float64 `json:"price"`
		Picture     string  `json:"picture"`
	}
	ServiceListsResponse struct {
		Services   []ServiceItem `json:"data"`
		Pagination Pagination    `json:"pagination"`
	}

	Pagination struct {
		NextPage int `json:"next_page"`
	}

	SearchServiceRequest struct {
		Page         int
		Limit        int
		PageString   string
		LimitString  string
		Keyword      string
		Type         string
		Rating       float64
		RatingString string
		Sort         string
	}
)

const defaultSort = "highest price"

// nolint
func (req *ServiceListRequest) ValidateServiceListRequest() ([]handler.Fields, error) {
	var count int
	var fields []handler.Fields
	pageValid := true
	limitValid := true
	ratingValid := true
	err := validator.ValidatePage(req.PageString)
	if err != nil {
		pageValid = false
		count++
		fields = append(fields, handler.Fields{
			Name:    "page",
			Message: err.Error(),
		})
	}

	if pageValid {
		page, pageErr := strconv.Atoi(req.PageString)
		if pageErr != nil {
			count++
			pageValid = false
			fields = append(fields, handler.Fields{
				Name:    "page",
				Message: pageErr.Error(),
			})
		}
		if page < 1 && pageValid {
			count++
			fields = append(fields, handler.Fields{
				Name:    "page",
				Message: "page must be more than 0",
			})
		}
		req.Page = page
	}

	err = validator.ValidateLimit(req.LimitString)
	if err != nil {
		limitValid = false
		count++
		fields = append(fields, handler.Fields{
			Name:    "limit",
			Message: err.Error(),
		})
	}

	if limitValid {
		limit, limitErr := strconv.Atoi(req.LimitString)
		if limitErr != nil {
			count++
			limitValid = false
			fields = append(fields, handler.Fields{
				Name:    "limit",
				Message: limitErr.Error(),
			})
		}
		if limit < 1 && limitValid {
			count++
			fields = append(fields, handler.Fields{
				Name:    "limit",
				Message: "limit must be more than 0",
			})
		}
		req.Limit = limit
	}

	err = validator.ValidateFilterRating(req.RatingString)
	if err != nil {
		ratingValid = false
		count++
		fields = append(fields, handler.Fields{
			Name:    "rating",
			Message: err.Error(),
		})
	}

	if ratingValid {
		if req.RatingString == "" {
			req.Rating = 0
		} else {
			rating, ratingErr := strconv.ParseFloat(req.RatingString, 64)
			if ratingErr != nil {
				count++
				fields = append(fields, handler.Fields{
					Name:    "rating",
					Message: ratingErr.Error(),
				})
			}
			req.Rating = rating
		}
	}

	err = validator.ValidateSort(req.Sort)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "sort",
			Message: err.Error(),
		})
	}

	if req.Sort == "" {
		req.Sort = defaultSort
	}

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

// nolint
func (req *SearchServiceRequest) ValidateSearchService() ([]handler.Fields, error) {
	var count int
	var fields []handler.Fields
	pageValid := true
	limitValid := true
	ratingValid := true
	err := validator.ValidatePage(req.PageString)
	if err != nil {
		pageValid = false
		count++
		fields = append(fields, handler.Fields{
			Name:    "page",
			Message: err.Error(),
		})
	}

	if pageValid {
		page, pageErr := strconv.Atoi(req.PageString)
		if pageErr != nil {
			count++
			pageValid = false
			fields = append(fields, handler.Fields{
				Name:    "page",
				Message: pageErr.Error(),
			})
		}
		if page < 1 && pageValid {
			count++
			fields = append(fields, handler.Fields{
				Name:    "page",
				Message: "page must be more than 0",
			})
		}
		req.Page = page
	}

	err = validator.ValidateLimit(req.LimitString)
	if err != nil {
		limitValid = false
		count++
		fields = append(fields, handler.Fields{
			Name:    "limit",
			Message: err.Error(),
		})
	}

	if limitValid {
		limit, limitErr := strconv.Atoi(req.LimitString)
		if limitErr != nil {
			count++
			limitValid = false
			fields = append(fields, handler.Fields{
				Name:    "limit",
				Message: limitErr.Error(),
			})
		}
		if limit < 1 && limitValid {
			count++
			fields = append(fields, handler.Fields{
				Name:    "limit",
				Message: "limit must be more than 0",
			})
		}
		req.Limit = limit
	}

	err = validator.ValidateKeyword(req.Keyword)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "keyword",
			Message: err.Error(),
		})
	}

	err = validator.ValidateFilterRating(req.RatingString)
	if err != nil {
		ratingValid = false
		count++
		fields = append(fields, handler.Fields{
			Name:    "rating",
			Message: err.Error(),
		})
	}

	if ratingValid {
		if req.RatingString == "" {
			req.Rating = 0
		} else {
			rating, ratingErr := strconv.ParseFloat(req.RatingString, 64)
			if ratingErr != nil {
				count++
				fields = append(fields, handler.Fields{
					Name:    "rating",
					Message: ratingErr.Error(),
				})
			}
			req.Rating = rating
		}
	}

	err = validator.ValidateSort(req.Sort)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "sort",
			Message: err.Error(),
		})
	}

	if req.Sort == "" {
		req.Sort = defaultSort
	}

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

func (c *serviceCtx) GetAllServices(ctx context.Context, condition *ServiceListRequest) (ServiceListsResponse, error) {
	result := make([]ServiceItem, 0)
	offset := (condition.Page - 1) * condition.Limit
	nextPageOffset := condition.Page * condition.Limit
	res, err := c.serviceModel.GetAllServices(ctx, condition.Limit, offset, model.SortNFilter{
		Type:   condition.Type,
		Rating: condition.Rating,
		Sort:   condition.Sort,
	})
	if err != nil {
		log.Error().Err(fmt.Errorf("error when GetAllServices: %w", err)).Send()
		return ServiceListsResponse{}, &handler.InternalServerError
	}
	for _, v := range res {
		tmp := ServiceItem{
			ID:          v.ID,
			Title:       v.Title,
			Description: v.Description.String,
			Rating:      v.Rating,
			Price:       v.Price,
			Picture:     v.Picture.String,
		}
		result = append(result, tmp)
	}

	nextPageRes, err := c.serviceModel.GetAllServices(ctx, condition.Limit, nextPageOffset, model.SortNFilter{
		Type:   condition.Type,
		Rating: condition.Rating,
		Sort:   condition.Sort,
	})
	if err != nil {
		log.Error().Err(fmt.Errorf("error when GetAllServices: %w", err)).Send()
		return ServiceListsResponse{}, &handler.InternalServerError
	}

	if len(res) <= condition.Limit && len(nextPageRes) == 0 {
		return ServiceListsResponse{
			Services: result,
			Pagination: Pagination{
				NextPage: 0,
			},
		}, nil
	}
	return ServiceListsResponse{
		Services: result,
		Pagination: Pagination{
			NextPage: condition.Page + 1,
		},
	}, nil
}

func (c *serviceCtx) SearchService(ctx context.Context, condition *SearchServiceRequest) (ServiceListsResponse, error) {
	result := make([]ServiceItem, 0)
	offset := (condition.Page - 1) * condition.Limit
	nextPageOffset := condition.Page * condition.Limit
	res, err := c.serviceModel.SearchService(ctx, condition.Limit, offset, condition.Keyword, model.SortNFilter{
		Type:   condition.Type,
		Rating: condition.Rating,
		Sort:   condition.Sort,
	})
	if err != nil {
		log.Error().Err(fmt.Errorf("error when SearchService: %w", err)).Send()
		return ServiceListsResponse{}, &handler.InternalServerError
	}

	for _, v := range res {
		tmp := ServiceItem{
			ID:          v.ID,
			Title:       v.Title,
			Description: v.Description.String,
			Rating:      v.Rating,
			Price:       v.Price,
			Picture:     v.Picture.String,
		}
		result = append(result, tmp)
	}

	// nolint
	nextPageRes, err := c.serviceModel.SearchService(ctx, condition.Limit, nextPageOffset, condition.Keyword, model.SortNFilter{
		Type:   condition.Type,
		Rating: condition.Rating,
		Sort:   condition.Sort,
	})
	if err != nil {
		log.Error().Err(fmt.Errorf("error when SearchService: %w", err)).Send()
		return ServiceListsResponse{}, &handler.InternalServerError
	}

	if len(res) <= condition.Limit && len(nextPageRes) == 0 {
		return ServiceListsResponse{
			Services: result,
			Pagination: Pagination{
				NextPage: 0,
			},
		}, nil
	}
	return ServiceListsResponse{
		Services: result,
		Pagination: Pagination{
			NextPage: condition.Page + 1,
		},
	}, nil
}
