package v1

import (
	"e-montir/api/handler"
	"e-montir/controller"
	"net/http"
	"strings"
)

type ServiceHandler struct {
	serviceController controller.Service
}

func NewServiceHandler(serviceController controller.Service) ServiceHandler {
	return ServiceHandler{
		serviceController: serviceController,
	}
}

func (c *ServiceHandler) ListOfServices(w http.ResponseWriter, r *http.Request) {
	request := new(controller.ServiceListRequest)
	request.PageString = r.URL.Query().Get("page")
	request.LimitString = r.URL.Query().Get("limit")
	request.Type = r.URL.Query().Get("type")
	request.RatingString = r.URL.Query().Get("rating")
	request.Sort = r.URL.Query().Get("sort")
	request.Sort = strings.ToLower(request.Sort)

	fieldsErr, err := request.ValidateServiceListRequest()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	res, err := c.serviceController.GetAllServices(r.Context(), request)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, res)
}

func (c *ServiceHandler) SearchService(w http.ResponseWriter, r *http.Request) {
	request := new(controller.SearchServiceRequest)
	request.PageString = r.URL.Query().Get("page")
	request.LimitString = r.URL.Query().Get("limit")
	request.Keyword = r.URL.Query().Get("keyword")
	request.Keyword = strings.ToLower(request.Keyword)
	request.Type = r.URL.Query().Get("type")
	request.RatingString = r.URL.Query().Get("rating")
	request.Sort = r.URL.Query().Get("sort")
	request.Sort = strings.ToLower(request.Sort)

	fieldsErr, err := request.ValidateSearchService()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	res, err := c.serviceController.SearchService(r.Context(), request)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, res)
}
