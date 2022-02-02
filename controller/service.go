package controller

import (
	"context"
	"database/sql"
	"e-montir/api/handler"
	"e-montir/model"
	"e-montir/pkg/filter"
	"e-montir/pkg/sort"
	"e-montir/pkg/validator"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type serviceCtx struct {
	serviceModel model.Service
}

type Service interface {
	GetAllServices(ctx context.Context, userID string, condition *ServiceListRequest) (ServiceListsResponse, error)
	SearchService(ctx context.Context, condition *SearchServiceRequest) (ServiceListsResponse, error)
	AddFavService(ctx context.Context, userID string, serviceID int) error
	RemoveFavService(ctx context.Context, userID string, serviceID int) error
	ListOfFavServices(ctx context.Context, userID string) (FavServiceListResponse, error)
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

	FavServiceListResponse struct {
		Services []ServiceItem `json:"data"`
	}

	AddOrRemoveService struct {
		ServiceID       int
		ServiceIDString string
		UserID          string
	}
)

const defaultSort = sort.HighestRating
const defaultType = filter.All

// nolint:funlen:gocyclo // concise in 1 function
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
		req.RatingString = strings.TrimSpace(req.RatingString)
		if req.RatingString == "" {
			req.Rating = filter.RatingDefaultVal
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

	req.Sort = strings.TrimSpace(req.Sort)
	if req.Sort == "" {
		req.Sort = defaultSort
	}

	err = validator.ValidateCategory(req.Type)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "category",
			Message: err.Error(),
		})
	}

	req.Type = strings.TrimSpace(req.Type)
	if req.Type == "" {
		req.Type = defaultType
	}

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

// nolint:funlen:gocyclo // concise in 1 function
func (req *SearchServiceRequest) ValidateSearchServiceRequest() ([]handler.Fields, error) {
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
		req.RatingString = strings.TrimSpace(req.RatingString)
		if req.RatingString == "" {
			req.Rating = filter.RatingDefaultVal
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

	req.Sort = strings.TrimSpace(req.Sort)
	if req.Sort == "" {
		req.Sort = defaultSort
	}

	err = validator.ValidateCategory(req.Type)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "type",
			Message: err.Error(),
		})
	}

	req.Type = strings.TrimSpace(req.Type)
	if req.Type == "" {
		req.Type = defaultType
	}

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

func (req *AddOrRemoveService) ValidateAddOrRemoveService() ([]handler.Fields, error) {
	var count int
	var fields []handler.Fields
	serviceIDValid := true

	err := validator.ValidateServiceID(req.ServiceIDString)
	if err != nil {
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
			fields = append(fields, handler.Fields{
				Name:    "service_id",
				Message: serviceIDErr.Error(),
			})
		}
		req.ServiceID = serviceID
	}

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

// popular type will return service with number of order >= 30
func (c *serviceCtx) GetAllServices(ctx context.Context, userID string, condition *ServiceListRequest) (ServiceListsResponse, error) {
	result := make([]ServiceItem, 0)
	offset := (condition.Page - 1) * condition.Limit
	nextPageOffset := condition.Page * condition.Limit

	// userFavServices, err := c.serviceModel.ListOfFavServices(ctx, userID)
	// if err != nil {
	// 	if err != sql.ErrNoRows {
	// 		log.Error().Err(fmt.Errorf("error when get listOfFavServices: %w", err)).Send()
	// 		return ServiceListsResponse{}, &handler.InternalServerError
	// 	}
	// }

	// if err == nil {
	// 	for _, v := range userFavServices {
	// 		service, err := c.serviceModel.GetServiceByID(ctx, v)
	// 		if err != nil {
	// 			log.Error().Err(fmt.Errorf("error when getting getServiceByID: %w", err)).Send()
	// 			return ServiceListsResponse{}, &handler.InternalServerError
	// 		}
	// 		result = append(result, ServiceItem{
	// 			ID:          service.ID,
	// 			Title:       service.Title,
	// 			Description: service.Description.String,
	// 			Rating:      service.Rating,
	// 			Price:       service.Price,
	// 			Picture:     service.Picture.String,
	// 		})
	// 	}
	// }

	res, err := c.serviceModel.GetAllServices(ctx, userID, model.ListCriteria{
		Type:   condition.Type,
		Rating: condition.Rating,
		Sort:   condition.Sort,
		Offset: offset,
		Limit:  condition.Limit,
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

	nextPageRes, err := c.serviceModel.GetAllServices(ctx, userID, model.ListCriteria{
		Type:   condition.Type,
		Rating: condition.Rating,
		Sort:   condition.Sort,
		Offset: nextPageOffset,
		Limit:  condition.Limit,
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

// popular type will return service with number of order >= 30
func (c *serviceCtx) SearchService(ctx context.Context, condition *SearchServiceRequest) (ServiceListsResponse, error) {
	result := make([]ServiceItem, 0)
	offset := (condition.Page - 1) * condition.Limit
	nextPageOffset := condition.Page * condition.Limit
	res, err := c.serviceModel.SearchService(ctx, model.ListCriteria{
		Type:    condition.Type,
		Rating:  condition.Rating,
		Sort:    condition.Sort,
		Offset:  offset,
		Limit:   condition.Limit,
		Keyword: condition.Keyword,
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

	nextPageRes, err := c.serviceModel.SearchService(ctx, model.ListCriteria{
		Type:    condition.Type,
		Rating:  condition.Rating,
		Sort:    condition.Sort,
		Offset:  nextPageOffset,
		Limit:   condition.Limit,
		Keyword: condition.Keyword,
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

func (c *serviceCtx) AddFavService(ctx context.Context, userID string, serviceID int) error {
	_, err := c.serviceModel.GetServiceByID(ctx, serviceID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when getServiceByID: %w", err)).Send()
		if err == sql.ErrNoRows {
			return &handler.ServiceNotExists
		}
		return err
	}

	res, err := c.serviceModel.GetFavServiceByUserIDNServiceID(ctx, userID, serviceID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(fmt.Errorf("error when getFavServiceByUserIDNServiceID: %w", err)).Send()
			return err
		}
	}

	if res != -1 {
		return &handler.ServiceIsAlreadyFav
	}

	err = c.serviceModel.AddFavService(ctx, userID, serviceID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when addFavService: %w", err)).Send()
		return err
	}

	return nil
}

func (c *serviceCtx) RemoveFavService(ctx context.Context, userID string, serviceID int) error {
	_, err := c.serviceModel.GetServiceByID(ctx, serviceID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when getServiceByID: %w", err)).Send()
		if err == sql.ErrNoRows {
			return &handler.ServiceNotExists
		}
		return err
	}

	_, err = c.serviceModel.GetFavServiceByUserIDNServiceID(ctx, userID, serviceID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when getFavServiceByUserIDNServiceID: %w", err)).Send()
		if err == sql.ErrNoRows {
			return &handler.FavServiceNotExists
		}
		return err
	}

	err = c.serviceModel.RemoveFavService(ctx, userID, serviceID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when removeFavService: %w", err)).Send()
		return err
	}

	return nil
}

func (c *serviceCtx) ListOfFavServices(ctx context.Context, userID string) (FavServiceListResponse, error) {
	services := make([]ServiceItem, 0)
	res, err := c.serviceModel.ListOfFavServices(ctx, userID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when getting listOfFavServices: %w", err)).Send()
		return FavServiceListResponse{
			Services: services,
		}, err
	}

	for _, v := range res {
		service, err := c.serviceModel.GetServiceByID(ctx, v)
		if err != nil {
			log.Error().Err(fmt.Errorf("error when getting getServiceByID: %w", err)).Send()
			return FavServiceListResponse{
				Services: services,
			}, err
		}
		services = append(services, ServiceItem{
			ID:          service.ID,
			Title:       service.Title,
			Description: service.Description.String,
			Rating:      service.Rating,
			Price:       service.Price,
			Picture:     service.Picture.String,
		})
	}

	return FavServiceListResponse{
		Services: services,
	}, nil
}
