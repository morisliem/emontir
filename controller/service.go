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
	GetAllServices(ctx context.Context, page, limit int) (ServiceListsResponse, error)
	SearchService(ctx context.Context, page, limit int, keyword string) (ServiceListsResponse, error)
}

func NewService(serviceModel model.Service) Service {
	return &serviceCtx{
		serviceModel: serviceModel,
	}
}

type (
	ServiceListRequest struct {
		Page        int
		Limit       int
		PageString  string
		LimitString string
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
		Page        int
		Limit       int
		PageString  string
		LimitString string
		Keyword     string
	}
)

func (req *ServiceListRequest) ValidateServiceListRequest() ([]handler.Fields, error) {
	var count int
	var fields []handler.Fields
	pageValid := true
	limitValid := true
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

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

func (req *SearchServiceRequest) ValidateSearchService() ([]handler.Fields, error) {
	var count int
	var fields []handler.Fields
	pageValid := true
	limitValid := true
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

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

func (c *serviceCtx) GetAllServices(ctx context.Context, page, limit int) (ServiceListsResponse, error) {
	result := make([]ServiceItem, 0)
	offset := (page - 1) * limit
	nextPageOffset := page * limit
	res, err := c.serviceModel.GetAllServices(ctx, limit, offset)
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

	nextPageRes, err := c.serviceModel.GetAllServices(ctx, limit, nextPageOffset)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when GetAllServices: %w", err)).Send()
		return ServiceListsResponse{}, &handler.InternalServerError
	}

	if len(res) <= limit && len(nextPageRes) == 0 {
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
			NextPage: page + 1,
		},
	}, nil
}

func (c *serviceCtx) SearchService(ctx context.Context, page, limit int, keyword string) (ServiceListsResponse, error) {
	result := make([]ServiceItem, 0)
	offset := (page - 1) * limit
	nextPageOffset := page * limit
	res, err := c.serviceModel.SearchService(ctx, limit, offset, keyword)
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

	nextPageRes, err := c.serviceModel.SearchService(ctx, limit, nextPageOffset, keyword)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when SearchService: %w", err)).Send()
		return ServiceListsResponse{}, &handler.InternalServerError
	}

	if len(res) <= limit && len(nextPageRes) == 0 {
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
			NextPage: page + 1,
		},
	}, nil
}
