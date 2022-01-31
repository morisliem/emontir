package fcm

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/rs/zerolog/log"
)

type NotifFcm struct {
	To       string
	Title    string
	Body     string
	Redirect string
}

func SendNotification(ctx context.Context, notif NotifFcm) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("failed to establish connection: %s", err.Error()))
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error getting Messaging client: %s", err.Error()))
	}

	msg := &messaging.Message{
		Data: map[string]string{
			"redirect": notif.Redirect,
		},
		Token: notif.To,
		Notification: &messaging.Notification{
			Title: notif.Title,
			Body:  notif.Body,
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: notif.Title,
						Body:  notif.Body,
					},
				},
			},
		},
	}

	res, err := client.Send(ctx, msg)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error when sending message to client: %s", err.Error()))
	}

	fmt.Println(res)
}
