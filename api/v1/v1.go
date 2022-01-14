package v1

import (
	"e-montir/controller"
	"e-montir/pkg/mailer"
)

type Handler struct {
	Auth     AuthHandler
	Service  ServiceHandler
	Timeslot TimeslotHandler
	Cart     CartHandler
	User     UserHandler
	Order    OrderHandler
	Payment  PaymentHandler
}

func GetHandler(c controller.Manager, mailerCfg *mailer.Config) Handler {
	return Handler{
		Auth:     NewAuthHandler(c.Auth(), mailerCfg),
		Service:  NewServiceHandler(c.Service()),
		Timeslot: NewTimeslotHandler(c.Timeslot()),
		Cart:     NewCartHandler(c.Cart()),
		User:     NewUserHandler(c.User()),
		Order:    NewOrderHandler(c.Order()),
		Payment:  NewPaymentHandler(c.Payment(), c.Order()),
	}
}
