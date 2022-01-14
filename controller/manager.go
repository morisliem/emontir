package controller

import (
	"e-montir/model"
	"sync"
)

type Manager interface {
	Auth() Auth
	User() User
	Service() Service
	Timeslot() Timeslot
	Cart() Cart
	Order() Order
	Payment() Payment
}

type manager struct {
	modelManager model.Manager
}

func NewManager(modelManager model.Manager) Manager {
	sm := &manager{
		modelManager: modelManager,
	}
	return sm
}

var (
	authControllerOnce sync.Once
	authController     Auth
)

func (c *manager) Auth() Auth {
	authControllerOnce.Do(func() {
		authController = NewAuth(c.modelManager.User())
	})
	return authController
}

var (
	userControllerOnce sync.Once
	userController     User
)

func (c *manager) User() User {
	userControllerOnce.Do(func() {
		userController = NewUser(c.modelManager.User())
	})
	return userController
}

var (
	serviceControllerOnce sync.Once
	serviceController     Service
)

func (c *manager) Service() Service {
	serviceControllerOnce.Do(func() {
		serviceController = NewService(c.modelManager.Service())
	})
	return serviceController
}

var (
	timeslotControllerOnce sync.Once
	timeslotController     Timeslot
)

func (c *manager) Timeslot() Timeslot {
	timeslotControllerOnce.Do(func() {
		timeslotController = NewTimeslot(c.modelManager.Timeslot())
	})
	return timeslotController
}

var (
	cartControllerOnce sync.Once
	cartController     Cart
)

func (c *manager) Cart() Cart {
	cartControllerOnce.Do(func() {
		cartController = NewCart(c.modelManager.Cart(), c.modelManager.User())
	})
	return cartController
}

var (
	orderControllerOnce sync.Once
	orderController     Order
)

func (c *manager) Order() Order {
	orderControllerOnce.Do(func() {
		orderController = NewOrder(c.modelManager.Order(), c.modelManager.Cart(), c.modelManager.User())
	})
	return orderController
}

var (
	paymentControllerOnce sync.Once
	paymentController     Payment
)

func (c *manager) Payment() Payment {
	paymentControllerOnce.Do(func() {
		paymentController = NewPayment(c.modelManager.Order(), c.modelManager.Cart())
	})
	return paymentController
}
