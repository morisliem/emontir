package v1

import (
	"e-montir/api/handler"
	"e-montir/controller"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
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
	request.Type = strings.ToLower(request.Type)
	request.RatingString = r.URL.Query().Get("rating")
	request.Sort = r.URL.Query().Get("sort")
	request.Sort = strings.ToLower(request.Sort)
	userID := handler.GetTokenClaim(r.Context()).ID

	fieldsErr, err := request.ValidateServiceListRequest()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	res, err := c.serviceController.GetAllServices(r.Context(), userID, request)
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
	request.Type = strings.ToLower(request.Type)
	request.RatingString = r.URL.Query().Get("rating")
	request.Sort = r.URL.Query().Get("sort")
	request.Sort = strings.ToLower(request.Sort)

	fieldsErr, err := request.ValidateSearchServiceRequest()
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

func (c *ServiceHandler) AddFavService(w http.ResponseWriter, r *http.Request) {
	request := new(controller.AddOrRemoveService)
	request.UserID = handler.GetTokenClaim(r.Context()).ID
	request.ServiceIDString = chi.URLParam(r, "service_id")

	fieldsErr, err := request.ValidateAddOrRemoveService()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	err = c.serviceController.AddFavService(r.Context(), request.UserID, request.ServiceID)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, handler.DefaultSuccess{Success: true})
}

func (c *ServiceHandler) RemoveFavService(w http.ResponseWriter, r *http.Request) {
	request := new(controller.AddOrRemoveService)
	request.UserID = handler.GetTokenClaim(r.Context()).ID
	request.ServiceIDString = chi.URLParam(r, "service_id")

	fieldsErr, err := request.ValidateAddOrRemoveService()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	err = c.serviceController.RemoveFavService(r.Context(), request.UserID, request.ServiceID)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, handler.DefaultSuccess{Success: true})
}

func (c *ServiceHandler) ListOfFavServices(w http.ResponseWriter, r *http.Request) {
	userID := handler.GetTokenClaim(r.Context()).ID
	res, err := c.serviceController.ListOfFavServices(r.Context(), userID)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, res)
}
