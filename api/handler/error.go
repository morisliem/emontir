package handler

type EmontirError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var (
	DuplicatedEmailError       = EmontirError{Code: "AUTH-400-01", Message: "email used"}
	LoginFailedError           = EmontirError{Code: "AUTH-400-02", Message: "incorrect password" + "/" + "email"}
	ActivationEmailFailedError = EmontirError{Code: "AUTH-400-03", Message: "activation email failed"}
	ActivationLinkExpired      = EmontirError{Code: "AUTH-400-04", Message: "email activation link expired"}
	UnauthorizedError          = EmontirError{Code: "AUTH-401-01", Message: "token invalid"}
	EmailNotActivatedError     = EmontirError{Code: "AUTH-422-01", Message: "email not verified"}
	ParsePayloadError          = EmontirError{Code: "SERVER-400-01", Message: "failed to parse payload"}
	// nolint(gosec) // false positive
	CartAppointmentAvailable    = EmontirError{Code: "SERVER-400-02", Message: "appointment is exists, remove appointment before change appointment date or time"}
	NoEmployeeError             = EmontirError{Code: "SERVER-400-03", Message: "no employee available"}
	OrderHasBeenPaid            = EmontirError{Code: "SERVER-400-04", Message: "order has been paid"}
	ServiceIsReviewed           = EmontirError{Code: "SERVER-400-05", Message: "service has been reviewed"}
	ServiceIsAlreadyFav         = EmontirError{Code: "SERVER-400-06", Message: "service is already in the favorite list"}
	ServiceNotExists            = EmontirError{Code: "SERVER-404-01", Message: "service not exists"}
	CartAppointmentNotAvailable = EmontirError{Code: "SERVER-404-02", Message: "appointment not exists"}
	OrderNotExists              = EmontirError{Code: "SERVER-404-03", Message: "cannot make payment to not exist order"}
	FavServiceNotExists         = EmontirError{Code: "SERVER-404-04", Message: "favorite service not exists"}
	InternalServerError         = EmontirError{Code: "SERVER-500-01", Message: "server error"}
)

const (
	ValidationFailed = "validation-failed"
)

func (e *EmontirError) Error() string { return e.Message }
