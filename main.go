package main

import (
	"context"
	"e-montir/api/middleware"
	v1 "e-montir/api/v1"
	"e-montir/controller"
	"e-montir/model"
	"e-montir/pkg/mailer"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	_ = godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(osSigChan)
		os.Exit(0)
	}()

	mailerPort, err := strconv.Atoi(os.Getenv("MAILER_PORT"))
	if err != nil {
		mailerPort = 587
	}

	mailerCfg := &mailer.Config{
		Sender:   os.Getenv("MAILER_SENDER"),
		Username: os.Getenv("MAILER_USERNAME"),
		Password: os.Getenv("MAILER_PASSWORD"),
		Host:     os.Getenv("MAILER_HOST"),
		Port:     mailerPort,
	}

	r := createHandler(mailerCfg)

	readTimeout, err := time.ParseDuration(os.Getenv("READ_TIMEOUT"))
	if err != nil {
		readTimeout = 60000 * time.Millisecond
	}
	writeTimeout, err := time.ParseDuration(os.Getenv("WRITE_TIMEOUT"))
	if err != nil {
		writeTimeout = 60000 * time.Millisecond
	}
	shutdownTimeout, err := time.ParseDuration(os.Getenv("SHUTDOWN_TIMEOUT"))
	if err != nil {
		shutdownTimeout = 60000 * time.Millisecond
	}

	httpServer := &http.Server{
		// Addr: "0.0.0.0:" + port,
		Addr:         os.Getenv("SERVER_ADDR"),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      r,
	}

	shutdownCtx := context.Background()
	if shutdownTimeout > 0 {
		var cancelShotdownTimeout context.CancelFunc
		shutdownCtx, cancelShotdownTimeout = context.WithTimeout(shutdownCtx, shutdownTimeout)
		defer cancelShotdownTimeout()
	}

	err = httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Error().Err(err).Msg(err.Error())
	}
	log.Info().Msg(fmt.Sprintf("serving %s\n", os.Getenv("SERVER_ADDR")))

	go func(httpServer *http.Server) {
		<-osSigChan
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			panic("failed to shutdown gracefully")
		}
	}(httpServer)
}

func createHandler(mailerCfg *mailer.Config) http.Handler {
	m := model.NewManager()
	c := controller.NewManager(m)
	h := v1.GetHandler(c, mailerCfg)
	r := chi.NewRouter()

	r.Route("/api/v1", func(apiRoute chi.Router) {
		apiRoute.Post("/auth/register", h.Auth.Register)
		apiRoute.Post("/auth/login", h.Auth.Login)
		apiRoute.Get("/auth/verify", h.Auth.ActivateEmail)

		apiRoute.With(middleware.ValidateToken()).Get("/services", h.Service.ListOfServices)
		apiRoute.With(middleware.ValidateToken()).Get("/services/search", h.Service.SearchService)
		apiRoute.With(middleware.ValidateToken()).Post("/services/{order_id}/{service_id}/review", h.Review.AddServiceReview)

		apiRoute.With(middleware.ValidateToken()).Get("/timeslot", h.Timeslot.ListOfTimeslot)

		apiRoute.With(middleware.ValidateToken()).Get("/me/address", h.User.ListOfUserLocation)
		apiRoute.With(middleware.ValidateToken()).Post("/me/address", h.User.AddUserLocation)

		apiRoute.With(middleware.ValidateToken()).Get("/cart", h.Cart.GetCheckoutDetail)
		apiRoute.With(middleware.ValidateToken()).Post("/cart/item", h.Cart.AddServiceToCart)
		apiRoute.With(middleware.ValidateToken()).Delete("/cart/item/{service_id}", h.Cart.RemoveServiceFromCart)
		apiRoute.With(middleware.ValidateToken()).Post("/cart/appointment", h.Cart.SetCartAppointment)
		apiRoute.With(middleware.ValidateToken()).Delete("/cart/appointment", h.Cart.RemoveCartAppointment)

		apiRoute.With(middleware.ValidateToken()).Post("/pay/{order_id}", h.Payment.MakePayment)
		apiRoute.Post("/payment/notification", h.Payment.PaymentNotification)

		apiRoute.With(middleware.ValidateToken()).Post("/order", h.Order.PlaceOrder)
		apiRoute.With(middleware.ValidateToken()).Get("/orders", h.Order.OrderLists)
	})

	return r
}
